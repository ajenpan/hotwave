package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotwave/frame"
	"hotwave/logger"
	"hotwave/servers/gateway/gate/tcp"
	"hotwave/servers/gateway/handler"
	"hotwave/servers/gateway/proto"
)

var Name = "gateway"
var Version string = "unknow"
var GitCommit string = "unknow"
var BuildAt string = "unknow"
var BuildBy string = "unknow"

func main() {
	err := realmain()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(-1)
	}
}

func realmain() error {
	var err error
	gate := handler.NewGater()

	frame.DefaultOptions.Name = Name
	frame.DefaultOptions.Version = Version
	frame.DefaultOptions.Adpater = gate
	core, err := frame.NewCore()
	if err != nil {
		return err
	}

	proto.RegisterGatewayServer(core, gate)

	if err = core.Start(); err != nil {
		return err
	}

	tcpListener := tcp.NewServer(&tcp.ServerOptions{
		Adapter:          gate,
		HeatbeatInterval: time.Second * 10,
		Address:          ":18080",
	})
	if err = tcpListener.Start(); err != nil {
		return err

	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	s := <-signals
	logger.Infof("recv signal: %v", s.String())
	core.Stop()
	tcpListener.Stop()
	return nil
}
