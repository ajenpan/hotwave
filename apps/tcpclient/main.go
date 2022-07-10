package main

import (
	"os"
	"time"

	log "hotwave/logger"
	"hotwave/transport/tcp"
	utilSignal "hotwave/utils/signal"
)

func reconnectFunc(s *tcp.Client) {
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
	client := tcp.NewClient(&tcp.ClientOptions{
		RemoteAddress: remoteAddr,
		OnMessage: func(s *tcp.Client, p *tcp.Packet) {
			p2 := tcp.CopyPacket(p)
			go func() {
				log.Info("OnMessage:", s.ID(), ", len:", p2.RawLen, ", typ:", p2.Typ)
				time.Sleep(1 * time.Second)
				err := s.SendPacket(p2)
				if err != nil {
					log.Info("response message failed:", err)
				}
			}()
		},
		OnConnStat: func(s *tcp.Client, state tcp.SocketStat) {
			log.Info("OnConnStat:", tcp.SocketStatString(state), " id:", s.ID())
			if state == tcp.SocketStatConnected {
				err := s.SendPacket(tcp.NewPacket(tcp.PacketTypeEcho, []byte("hello world")))
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
