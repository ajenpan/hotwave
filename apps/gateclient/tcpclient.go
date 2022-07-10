package main

import (
	"google.golang.org/protobuf/proto"

	log "hotwave/logger"
	gwProto "hotwave/service/gateway/proto"
	"hotwave/transport/tcp"
	utilSignal "hotwave/utils/signal"
)

func SendMsg(client *tcp.Client, msg proto.Message) {
	raw, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}

	warp := &gwProto.ClientMessage{
		Name: string(proto.MessageName(msg)),
		Body: raw,
	}

	raw, err = proto.Marshal(warp)
	if err != nil {
		log.Error(err)
		return
	}
	if err := client.Send(raw); err != nil {
		log.Error(err)
	}
}

func TCPClient() {
	client := tcp.NewClient(&tcp.ClientOptions{
		RemoteAddress: "localhost:10010",
		OnMessage: func(s *tcp.Client, p *tcp.Packet) {
			log.Info("OnMessage:", s.ID(), ", len:", p.RawLen, ", typ:", p.Typ)
		},
		OnConnStat: func(s *tcp.Client, state tcp.SocketStat) {
			log.Info("OnConnStat:", tcp.SocketStatString(state), " id:", s.ID())
		},
	})

	if err := client.Connect(); err != nil {
		panic(err)
	}

	SendMsg(client, &gwProto.LoginRequest{
		Checker: &gwProto.LoginRequest_Account{
			Account: &gwProto.AccountInfo{
				Account: "root",
				Passwd:  "123456",
			},
		},
	})

	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
}
