package main

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "hotwave/logger"
	gwClient "hotwave/service/gateway/client"
	gwProto "hotwave/service/gateway/proto"
	utilSignal "hotwave/utils/signal"
)

func GRPCClient() {
	log.Info("start")
	conn, err := grpc.Dial("localhost:10000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := gwClient.GRPCClient{
		NodeID:   "node1",
		NodeType: "node1",
		GrpcConn: conn,
	}
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	// client := gwProto.NewGatewayClient(conn)
	// md := metadata.New(map[string]string{"nodeid": "123", "nodename": "node1"})
	// ctx := metadata.NewOutgoingContext(context.Background(), md)
	// proxyc, err := client.Proxy(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	go func() {
		return
		for {
			err := client.SendMessage(&gwProto.ToUserMessage{})
			if err != nil {
				log.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	utilSignal.WaitShutdown()
}
