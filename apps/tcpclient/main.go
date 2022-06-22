package main

import (
	"os"
	"time"

	log "hotwave/logger"
	"hotwave/transport/tcpsvr"
	utilSignal "hotwave/util/signal"
)

func reconnectFunc(s *tcpsvr.Client) {
	time.Sleep(2 * time.Second)

	log.Info("reconnect to ", s.Opt.RemoteAddress)
	err := s.Connect()
	if err != nil {
		log.Error(err)
		go reconnectFunc(s)
	}
}

func main() {
	if len(os.Args) != 3 {
		return
	}
	remote := os.Args[1]
	port := os.Args[2]

	remoteAddr := remote + ":" + port
	client := tcpsvr.NewClient(&tcpsvr.ClientOptions{
		RemoteAddress: remoteAddr,
		OnMessage: func(s *tcpsvr.Client, p *tcpsvr.Packet) {
			p2 := tcpsvr.CopyPacket(p)

			go func() {
				log.Info("OnMessage:", s.ID(), ", len:", p2.RawLen, ", typ:", p2.Typ)
				time.Sleep(1 * time.Second)
				err := s.Send(p2)
				if err != nil {
					log.Info("response message failed:", err)
				}
			}()
		},
		OnConnStat: func(s *tcpsvr.Client, state tcpsvr.SocketStat) {
			log.Info("OnConnStat:", tcpsvr.SocketStatString(state), s.ID())
			if state == tcpsvr.SocketStatConnected {
				err := s.Send(tcpsvr.NewPacket(tcpsvr.PacketTypeEcho, []byte("hello world")))
				if err != nil {
					log.Info("send failed", err)
				}
			} else {
				go reconnectFunc(s)
			}
		},
	})
	defer client.Close()

	go reconnectFunc(client)

	s := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", s.String())
}
