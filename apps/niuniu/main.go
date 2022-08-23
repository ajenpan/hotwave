package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"hotwave/event"
	log "hotwave/logger"
	gwAuth "hotwave/service/gateway/auth"
	gwHandler "hotwave/service/gateway/handler"
	gwProto "hotwave/service/gateway/proto"
	tcpGate "hotwave/transport/tcp"
	utilSignal "hotwave/utils/signal"
)

var Version string = "unknown"
var GitCommit string = "unknown"
var BuildAt string = "unknown"
var BuildBy string = "unknown"
var Name string = "niuniu"

var ConfigPath string = ""
var ListenAddr string = ""
var PrintConf bool = false

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("project:", Name)
		fmt.Println("version:", Version)
		fmt.Println("git commit:", GitCommit)
		fmt.Println("build at:", BuildAt)
		fmt.Println("build by:", BuildBy)
	}

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       "config.yaml",
			Destination: &ConfigPath,
		}, &cli.StringFlag{
			Name:        "listen",
			Aliases:     []string{"l"},
			Value:       ":10110",
			Destination: &ListenAddr,
		}, &cli.BoolFlag{
			Name:        "print-config",
			Destination: &PrintConf,
			Hidden:      true,
		},
	}

	app.Action = RealMain

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func RealMain(c *cli.Context) error {
	publisher := &event.GrpcEventPublisher{}

	gw, err := gwHandler.NewGateway(gwHandler.GatewayOption{
		AuthClient: &gwAuth.FakeAuth{},
		Publisher:  publisher,
	})

	if err != nil {
		panic(err)
	}

	gw.AddAllowlistByMsg(&gwProto.LoginGateRequest{})

	gate := tcpGate.NewServer(tcpGate.ServerOptions{
		Address:   ":10000",
		OnMessage: gw.OnGateMessage,
		OnConn:    gw.OnGateConnStat,
	})

	if err := gate.Start(); err != nil {
		panic(err)
	}
	defer gate.Stop()

	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
	return nil
}
