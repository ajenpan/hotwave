package main

import (
	"context"
	"fmt"
	"time"

	log "hotwave/logger"
	"hotwave/service/route/client"
	routeMsg "hotwave/service/route/proto"
	"hotwave/service/route/transport/tcp"
	utilSignal "hotwave/utils/signal"
)

func main() {

	loginReq := &routeMsg.AccountLoginRequest{
		Account:  "test",
		Password: "123456",
	}
	loginResp := &routeMsg.AccountLoginResponse{}
	c := client.NewTcpClient("localhost:14321")
	c.AuthFunc = func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := c.SyncCall(0, ctx, loginReq, loginResp)
		ok := err == nil && loginResp.Errcode == 0
		return ok
	}

	c.OnLoginFunc = func(c *client.TcpClient, stat client.LoginStat) {
		if stat == client.LoginStat_Success {
			fmt.Println("login success")
			req := &routeMsg.EchoRequest{
				Msg: "hello",
			}
			client.SendRequestWithCB(c, 0, context.Background(), req, func(err error, c *client.TcpClient, resp *routeMsg.EchoResponse) {
				if err != nil {
					fmt.Println("send request failed:", err)
				} else {
					fmt.Println("recv resp:", resp.Msg)
				}
			})
		} else {
			fmt.Println("login failed")
		}
	}

	c.OnMessageFunc = func(c *client.TcpClient, p *tcp.PackFrame) {

	}
	c.Reconnect()

	s := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", s.String())
}
