package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	_ "hotwave/games/niuniu"
	"hotwave/logger"
	battleHandler "hotwave/service/battle/handler"
	"hotwave/transport/tcp"
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

func RealMain(c *cli.Context) error {

	h := battleHandler.New()

	listener, err := tcp.NewServer(tcp.ServerOptions{
		Address:   listenAt,
		OnMessage: h.OnMessage,
		OnConn:    h.OnConn,
		AuthTokenChecker: func(s string) (*tcp.UserInfo, error) {
			return &tcp.UserInfo{
				Uid:   1,
				Uname: "1",
				Role:  "test",
			}, nil
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
