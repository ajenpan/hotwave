package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"

	"hotwave/logger"
	"hotwave/service/auth/common"
	"hotwave/transport/tcp"
	"hotwave/utils/rsagen"
	utilSignal "hotwave/utils/signal"
)

const PrivateKeyFile = "private.pem"
const PublicKeyFile = "public.pem"

func ReadRSAKey() ([]byte, []byte, error) {
	privateRaw, err := os.ReadFile(PrivateKeyFile)
	if err != nil {
		privateKey, publicKey, err := rsagen.GenerateRsaPem(512)
		if err != nil {
			return nil, nil, err
		}
		privateRaw = []byte(privateKey)
		os.WriteFile(PrivateKeyFile, []byte(privateKey), 0644)
		os.WriteFile(PublicKeyFile, []byte(publicKey), 0644)
	}
	publicRaw, err := os.ReadFile(PublicKeyFile)
	if err != nil {
		return nil, nil, err
	}
	return privateRaw, publicRaw, nil
}

func main() {
	app := cli.NewApp()
	app.Action = RealMain
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func RealMain(c *cli.Context) error {
	args := c.Args()
	if args.Len() < 2 {
		return fmt.Errorf("args error")
	}
	remote := args.Get(0)
	uname := args.Get(1)
	fmt.Println("uname:", uname)
	uid, _ := strconv.ParseUint(uname, 10, 64)

	privateRaw, _, err := ReadRSAKey()
	if err != nil {
		return err
	}
	pk, err := rsagen.ParseRsaPrivateKeyFromPem(privateRaw)
	if err != nil {
		return err
	}
	token, err := common.GenerateToken(pk, &common.UserClaims{
		UID:   uid,
		UName: uname,
		Role:  "user",
	})

	if err != nil {
		return err
	}

	opts := &tcp.ClientOptions{
		RemoteAddress: remote,
		OnMessage: func(s *tcp.Socket, t *tcp.THVPacket) {

		},
		OnConnStat: func(s *tcp.Socket, ok bool) {
			if ok {
				fmt.Println("connected")
			} else {
				fmt.Println("disconnected")
			}
		},
		Token: token,
	}

	client := tcp.NewClient(opts)

	err = client.Connect()
	if err != nil {
		return err
	}

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())

	return nil
}
