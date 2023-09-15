package main

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/cli/v2"

	_ "hotwave/games/niuniu"
	"hotwave/logger"
	battleHandler "hotwave/service/battle/handler"
	"hotwave/transport/tcp"
	"hotwave/utils/rsagen"
	utilSignal "hotwave/utils/signal"
)

var (
	Name       string = "unknown"
	Version    string = "unknown"
	GitCommit  string = "unknown"
	BuildAt    string = "unknown"
	BuildBy    string = runtime.Version()
	RunnningOS string = runtime.GOOS + "/" + runtime.GOARCH
)

func longVersion() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintln(buf, "project:", Name)
	fmt.Fprintln(buf, "version:", Version)
	fmt.Fprintln(buf, "git commit:", GitCommit)
	fmt.Fprintln(buf, "build at:", BuildAt)
	fmt.Fprintln(buf, "build by:", BuildBy)
	fmt.Fprintln(buf, "running OS/Arch:", RunnningOS)
	return buf.String()
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Run() error {
	Name = "battle"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(longVersion())
	}

	app := cli.NewApp()
	app.Version = Version
	app.Name = Name
	app.Action = RealMain

	err := app.Run(os.Args)
	return err
}

var listenAt string = ":12345"

func LoadAuthPublicKey() (*rsa.PublicKey, error) {
	publicRaw, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}
	pk, err := rsagen.ParseRsaPublicKeyFromPem(publicRaw)
	return pk, err
}

func RealMain(c *cli.Context) error {
	pk, err := LoadAuthPublicKey()
	if err != nil {
		panic(err)
	}

	h := battleHandler.New()

	listener, err := tcp.NewServer(tcp.ServerOptions{
		Address:   listenAt,
		OnMessage: h.OnMessage,
		OnConn:    h.OnConn,
		AuthTokenChecker: func(tokenRaw string) (*tcp.UserInfo, error) {
			claims := make(jwt.MapClaims)
			token, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (interface{}, error) {
				return pk, nil
			})
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, fmt.Errorf("invalid token")
			}
			ret := &tcp.UserInfo{}
			if uname, has := claims["aud"]; has {
				ret.Uname = uname.(string)
			}
			if uidstr, has := claims["uid"]; has {
				uid, _ := strconv.ParseUint(uidstr.(string), 10, 64)
				ret.Uid = uid
			}
			if role, has := claims["rid"]; has {
				ret.Role = role.(string)
			}
			return ret, nil
		},
	})

	if err != nil {
		panic(err)
	}

	go listener.Start()
	defer listener.Stop()

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	return nil
}
