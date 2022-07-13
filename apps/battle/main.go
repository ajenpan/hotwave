package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"hotwave/logger"
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

		// h := battleHandler.New()

		// svr := tcpsvr.NewServer(&tcpsvr.ServerOptions{
		// 	Address: ":10010",
		// 	OnConn: func(s *tcpsvr.Socket, ss tcpsvr.SocketStat) {
		// 		logger.Infof("socket:%s conn:%s", s.ID(), tcpsvr.SocketStatString(ss))
		// 	},
		// 	OnMessage: func(s *tcpsvr.Socket, p *tcpsvr.Packet) {
		// 		msg := &proto.BattleMessageWrap{}
		// 		h.OnBattleMessage(context.Background(), msg)
		// 	},
		// })

		// err := svr.Start()
		// if err != nil {
		// 	return err
		// }
		// defer svr.Stop()

		// core := frame.New(
		// 	frame.Name(Name),
		// 	frame.Version(Version),
		// 	frame.Address(":10010"),
		// )

		// h, err := handler.New(config.DefaultConf)
		// if err != nil {
		// 	panic(err)
		// }
		// reflection.Register(core.GrpcServer())
		// if err := core.Start(); err != nil {
		// 	return err
		// }
		// defer core.Stop()
		// httpServer(h)

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}
