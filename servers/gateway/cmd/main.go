package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	frame "hotwave"
	"hotwave/logger"
	httpgate "hotwave/servers/gateway/gate/http"
	tcpgate "hotwave/servers/gateway/gate/tcp"
	"hotwave/servers/gateway/handler"
	"hotwave/servers/gateway/proto"
	utilSignal "hotwave/util/signal"
)

var Name = "gateway"
var Version string = "unknow"
var GitCommit string = "unknow"
var BuildAt string = "unknow"
var BuildBy string = "unknow"

func main() {
	err := RealMain()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(-1)
	}
}

func RealMain() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("project:", Name)
		fmt.Println("version:", Version)
		fmt.Println("git commit:", GitCommit)
		fmt.Println("build at:", BuildAt)
		fmt.Println("build by:", BuildBy)
	}
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		core := frame.New(
			frame.Name(Name),
			frame.Version(Version),
		)

		gate := handler.NewGater(core)
		proto.RegisterGatewayServer(core, gate)

		if err := core.Start(); err != nil {
			return err
		}
		defer core.Stop()

		httpListener := httpgate.NewServer(httpgate.Options{
			Address: ":10087",
			Adapter: gate,
		})
		if err := httpListener.Start(); err != nil {
			panic(err)
		}
		defer httpListener.Stop()

		tcpListener := tcpgate.NewServer(&tcpgate.ServerOptions{
			Adapter:          gate,
			HeatbeatInterval: time.Second * 20,
			Address:          ":10086",
		})
		if err := tcpListener.Start(); err != nil {
			return err
		}
		fmt.Println("tcp gate listen on", tcpListener.Address())
		defer tcpListener.Stop()

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}
