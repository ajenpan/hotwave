package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/reflection"

	frame "hotwave"
	"hotwave/logger"
	"hotwave/services/auth/config"
	"hotwave/services/auth/handler"
	"hotwave/services/auth/proto"
	utilHandle "hotwave/util/handle"
	utilSignal "hotwave/util/signal"
)

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Run() error {
	Name = "auth"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(longVersion())
	}

	app := cli.NewApp()
	app.Version = Version
	app.Name = Name

	app.Action = func(c *cli.Context) error {
		core := frame.New(
			frame.Name(Name),
			frame.Version(Version),
			frame.Address(":10010"),
		)

		h, err := handler.New(config.DefaultConf)
		if err != nil {
			panic(err)
		}

		reflection.Register(core.GrpcServer())

		if err := core.Start(); err != nil {
			return err
		}
		defer core.Stop()

		httpServer(h)

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}

var calltable = utilHandle.ExtractProtoFile(proto.File_auth_proto, &handler.Handler{})

func httpServer(handler interface{}) {
	go func() {
		http.HandleFunc("/", utilHandle.ServerGRPCMethodForHttp(handler, calltable))
		err := http.ListenAndServe(":8088", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}
