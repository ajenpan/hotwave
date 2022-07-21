package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "hotwave/game/niuniu"
	"hotwave/logger"
	"hotwave/service/battle"
	battleHandler "hotwave/service/battle/handler"
	gwclient "hotwave/service/gateway/client"
	"hotwave/transport"
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

	app.Action = func(c *cli.Context) error {
		h := battleHandler.New()
		grpcConn, err := grpc.Dial("localhost:20000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		gwc := &gwclient.GRPCClient{
			GrpcConn: grpcConn,
			NodeID:   "battle",
			NodeType: "battle",

			OnConnStatusFunc: func(c *gwclient.GRPCClient, ss transport.SessionStat) {
				if ss == transport.Connected {
				} else {
					go c.Reconnect()
				}
			},
			OnUserMessageFunc: h.OnUserMessage,
		}

		battle.LogicCreator.Store.Range(func(key, value any) bool {
			logger.Info("reg game:", key.(string))
			return true
		})

		gwc.Reconnect()
		defer gwc.Close()

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}
