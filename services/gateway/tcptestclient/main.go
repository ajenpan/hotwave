package main

import (
	"fmt"
	"os"
	"strings"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/services/gateway/gate"
	gatetcp "hotwave/services/gateway/gate/tcp"
	"hotwave/services/gateway/proto"
	utilsig "hotwave/util/signal"
)

func SendMsg(s gate.Session, msg protobuf.Message) error {
	data, err := protobuf.Marshal(msg)
	if err != nil {
		return err
	}
	msgName := string(protobuf.MessageName(msg))
	msgName = strings.TrimSuffix(msgName, "Request")
	warp := &proto.ClientMessageWraper{
		Method: msgName,
		Body:   data,
	}

	fmt.Println(warp.String())

	if s == nil {
		return fmt.Errorf("session is nil")
	}
	return s.Send(warp)
}

func main() {

	SendMsg(nil, &proto.LoginRequest{
		Checker: &proto.LoginRequest_AccountInfo{
			AccountInfo: &proto.AccountInfo{
				Account: "test",
				Passwd:  "123456",
			},
		},
	})

	args := os.Args[1:]
	address := "localhost:10086"
	if len(args) == 1 {
		address = args[0]
	}

	fmt.Println("connect to", address)

	c := gatetcp.NewClient(&gatetcp.ClientOptions{
		RemoteAddress: address,
		OnMessage: func(s *gatetcp.Client, p *gatetcp.Packet) {
			fmt.Printf("session:%s recv msg \n", s.ID())
		},
		OnConnStat: func(s *gatetcp.Client, state gatetcp.SocketStat) {
			fmt.Printf("session:%s, connect state:%v \n", s.ID(), state)
			switch state {
			case gatetcp.SocketStatConnected:
			case gatetcp.SocketStatDisconnected:
				os.Exit(1)
			}
		},
	})

	var err error

	if err = c.Connect(); err != nil {
		fmt.Println("connect error:", err)
		return
	}

	err = SendMsg(c, &proto.LoginRequest{
		Checker: &proto.LoginRequest_AccountInfo{
			AccountInfo: &proto.AccountInfo{
				Account: "test",
				Passwd:  "123456",
			},
		},
	})

	if err != nil {
		fmt.Println("send error:", err)
		return
	}

	s := utilsig.WaitShutdown()
	fmt.Println("shutdown on:", s.String())
}
