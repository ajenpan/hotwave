package main

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	log "hotwave/logger"
	auproto "hotwave/service/auth/proto"
	gwproto "hotwave/service/gateway/proto"
	"hotwave/transport/tcp"
	utilSignal "hotwave/utils/signal"
)

func SendMsg(client *tcp.Client, msg proto.Message) {
	raw, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}

	warp := &gwproto.GateClientMessage{
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

func RecvMsg[T any](recvchan chan *gwproto.GateServerMessage) (*T, error) {
	t := new(T)
	msg, ok := <-recvchan
	if !ok {
		return t, fmt.Errorf("recv chan closed")
	}
	if err := proto.Unmarshal(msg.Body, any(t).(proto.Message)); err != nil {
		log.Error(err)
		return t, err
	}
	return t, nil
}

func TCPClient() {
	recvchan := make(chan *gwproto.GateServerMessage, 10)
	client := tcp.NewClient(&tcp.ClientOptions{
		RemoteAddress: "localhost:10010",
		OnMessage: func(s *tcp.Client, p *tcp.Packet) {
			warp := &gwproto.GateServerMessage{}
			if err := proto.Unmarshal(p.Raw, warp); err != nil {
				log.Error(err)
				return
			}
			log.Info("recv msg:", warp.Name)
			recvchan <- warp
		},
		OnConnStat: func(s *tcp.Client, state tcp.SocketStat) {
			log.Info("OnConnStat:", state, " id:", s.ID())
		},
	})

	if err := client.Connect(); err != nil {
		panic(err)
	}

	recvMsg := func(recv proto.Message) {
		msg, ok := <-recvchan
		if !ok {
			return
		}
		if err := proto.Unmarshal(msg.Body, recv); err != nil {
			log.Error(err)
			return
		}
	}

	SendMsg(client, &auproto.LoginRequest{
		Uname:  "test",
		Passwd: "123456",
	})

	loginResp := &auproto.LoginResponse{}
	recvMsg(loginResp)
	if loginResp.Flag != 0 {
		log.Error("login failed:", loginResp.Msg)
		return
	}

	RecvMsg[gwproto.GateClientMessage](recvchan)

	SendMsg(client, &gwproto.LoginGateRequest{
		// Checker: &gwproto.LoginGateRequest_Account{
		// 	Account: &gwproto.AccountInfo{
		// 		Account: "root",
		// 		Passwd:  "123456",
		// 	},
		// },
		Checker: &gwproto.LoginGateRequest_Jwt{
			Jwt: loginResp.AssessToken,
		},
	})

	loginGateResp, _ := RecvMsg[gwproto.LoginGateResponse](recvchan)
	log.Info("sid", loginGateResp.Sessionid)

	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
}
