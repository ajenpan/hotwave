package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	frame "hotwave"
	"hotwave/logger"
	"hotwave/services/battle/proto"
	utilSignal "hotwave/util/signal"
)

var Name string = string(proto.File_servers_battle_proto_battle_proto.Package())
var Version string = "unknow"
var GitCommit string = "unknow"
var BuildAt string = "unknow"
var BuildBy string = "unknow"

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

	app.Action = func(c *cli.Context) error {
		core := frame.New(
			frame.Name(Name),
			frame.Version(Version),
			frame.Address(":10001"),
		)

		go func() {
			w, err := core.Options().Registry.Watch()
			if err != nil {
				logger.Info("watch error:", err)
				return
			}

			for {
				res, err := w.Next()
				if err != nil {
					logger.Errorf("watch error: %s", err)
					break
				}
				if len(res.Service.Nodes) == 0 {
					logger.Infof("watch change: %s %s", res.Action, res.Service.Name)
					continue
				}
				node := res.Service.Nodes[0]
				logger.Infof("watch res: %s %s", res.Action, node.Id)
			}
		}()

		// h := handler.New()
		// proto.RegisterBattleServer(core, h)
		// gateproto.RegisterGateAdapterServer(core, h)

		core.Start()
		defer core.Stop()

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
