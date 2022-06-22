package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"hotwave/services/auth/handler"
	"hotwave/services/auth/proto"
	"hotwave/transport"
)

var (
	Name       string = "unknown"
	Version    string = "unknown"
	GitCommit  string = "unknown"
	BuildAt    string = "unknown"
	BuildBy    string = runtime.Version()
	RunnningOS string = runtime.GOOS + "/" + runtime.GOARCH
)

func shortVersion() string {
	return Version
}

func longVersion() string {
	buf := bytes.NewBuffer(nil)
	fmt.Println(buf, "project:", Name)
	fmt.Println(buf, "version:", Version)
	fmt.Println(buf, "git commit:", GitCommit)
	fmt.Println(buf, "build at:", BuildAt)
	fmt.Println(buf, "build by:", BuildBy)
	fmt.Fprintln(buf, "Running OS/Arch:", RunnningOS)
	return buf.String()
}

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

		// s := utilSignal.WaitShutdown()
		// logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}

var calltable = transport.ExtractProtoFile(proto.File_auth_proto, &handler.Handler{})

func httpServer(handler interface{}) {
	go func() {
		http.HandleFunc("/", transport.ServerGRPCMethodForHttp(handler, calltable))
		err := http.ListenAndServe(":8088", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}
