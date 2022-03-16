package main

import (
	"fmt"
	"os"
	"time"

	"hotwave/frame"
	"hotwave/logger"
	"hotwave/servers/gateway/gate/tcp"
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
	core := frame.NewFrame(
		frame.Name(Name),
		frame.Version(Version),
	)

	gate := handler.NewGater()
	proto.RegisterGatewayServer(core, gate)

	if err := core.Start(); err != nil {
		return err
	}

	tcpListener := tcp.NewServer(&tcp.ServerOptions{
		Adapter:          gate,
		HeatbeatInterval: time.Second * 10,
		Address:          ":0",
	})
	if err := tcpListener.Start(); err != nil {
		return err
	}
	fmt.Println("tcp gate listen on", tcpListener.Address())
	tk := time.NewTicker(time.Second * 20)

	go func() {
		defer tk.Stop()
		for v := range tk.C {
			logger.Infof("tick: %v", v)
			svrs := core.GetService("gateway")
			if svrs == nil {
				continue
			}
			if len(svrs.Nodes) == 0 {
				continue
			}
			for _, node := range svrs.Nodes {
				fmt.Println(node.Id)
			}
		}
	}()
	// time.Sleep(time.Second * 10)
	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	tk.Stop()
	core.Stop()
	tcpListener.Stop()
	return nil
}
