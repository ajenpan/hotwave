package client

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"hotwave/service/route/transport/tcp"
	"hotwave/utils/calltable"
)

type LoginStat int

const (
	LoginStat_Success    LoginStat = iota
	LoginStat_Fail       LoginStat = iota
	LoginStat_Disconnect LoginStat = iota
)

type TcpClient struct {
	*tcp.Client

	cb sync.Map

	isAuth   bool
	AuthFunc func() bool

	eventque chan func()

	OnMessageFunc func(c *TcpClient, head tcp.RoutDeliverHead, p *tcp.PackFrame)
	OnLoginFunc   func(c *TcpClient, stat LoginStat)
	AutoRecconect bool

	reconnectTimeDelay time.Duration
}

func (c *TcpClient) Reconnect() {
	err := c.Connect()
	if err != nil && c.AutoRecconect {
		fmt.Println("start to reconnect")
		time.AfterFunc(c.reconnectTimeDelay, func() {
			c.Reconnect()
		})
	}
}

func (c *TcpClient) SetCallback(askid uint32, f func(*tcp.Socket, *tcp.PackFrame)) {
	c.cb.Store(askid, f)
}

func (c *TcpClient) RemoveCallback(askid uint32) {
	c.cb.Delete(askid)
}

func (c *TcpClient) MakeRequestPacket(target uint64, req proto.Message) (*tcp.PackFrame, uint32, error) {
	msgid := calltable.GetMessageMsgID(req.ProtoReflect().Descriptor())
	if msgid == 0 {
		return nil, 0, fmt.Errorf("not found msgid:%v", msgid)
	}

	raw, err := proto.Marshal(req)
	if err != nil {
		return nil, 0, err
	}
	askid := c.GetAskID()
	head := tcp.NewRoutDeliverHead()
	head.SetAskID(askid)
	head.SetMsgID(uint32(msgid))
	head.SetTargetUID(target)
	head.SetMsgTyp(tcp.RoutTypRequest)

	ret := &tcp.PackFrame{
		Head: head,
		Body: raw,
	}
	ret.SetType(tcp.PacketTypRoutDeliver)
	return ret, askid, nil
}

func SendRequestWithCB[T proto.Message](c *TcpClient, target uint64, ctx context.Context, req proto.Message, cb func(error, *TcpClient, T)) {
	go func() {
		var tresp T
		rsep := reflect.New(reflect.TypeOf(tresp).Elem()).Interface().(T)
		err := c.SyncCall(target, ctx, req, rsep)
		cb(err, c, rsep)
	}()
}

func (c *TcpClient) SyncCall(target uint64, ctx context.Context, req proto.Message, resp proto.Message) error {
	var err error

	packet, askid, err := c.MakeRequestPacket(target, req)
	if err != nil {
		return err
	}

	res := make(chan error, 1)

	c.SetCallback(askid, func(c *tcp.Socket, packet *tcp.PackFrame) {
		var err error
		defer func() {
			res <- err
		}()

		if err = proto.Unmarshal(packet.Body, resp); err != nil {
			return
		}
	})

	err = c.SendPacket(packet)

	if err != nil {
		c.RemoveCallback(askid)
		return err
	}

	select {
	case err = <-res:
		return err
	case <-ctx.Done():
		// dismiss callback
		c.SetCallback(askid, func(c *tcp.Socket, packet *tcp.PackFrame) {})
		return ctx.Err()
	}
}

func (r *TcpClient) AsyncCall(target uint64, m proto.Message) error {
	raw, err := proto.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal failed:%v", err)
	}

	msgid := calltable.GetMessageMsgID(m.ProtoReflect().Descriptor())
	if msgid == 0 {
		return fmt.Errorf("not found msgid:%v", msgid)
	}

	head := tcp.NewRoutDeliverHead()
	head.SetMsgID(uint32(msgid))
	head.SetSrouceUID(r.UID())
	head.SetTargetUID(target)
	head.SetMsgTyp(tcp.RoutTypAsync)
	p := &tcp.PackFrame{
		Head: head,
		Body: raw,
	}
	p.SetType(tcp.PacketTypRoutDeliver)
	return r.SendPacket(p)
}

func NewTcpClient(remoteAddr string) *TcpClient {
	ret := &TcpClient{
		eventque:      make(chan func(), 100),
		AutoRecconect: true,
	}

	go func() {
		for f := range ret.eventque {
			f()
		}
	}()

	c := tcp.NewClient(&tcp.ClientOptions{
		RemoteAddress: remoteAddr,
		OnMessage: func(s *tcp.Socket, p *tcp.PackFrame) {
			ptype := p.GetType()
			if ptype == tcp.PacketTypRoutDeliver {
				head := tcp.RoutDeliverHead(p.Head)
				askid := head.GetAskID()
				msgtyp := head.GetMsgTyp()

				if askid != 0 && msgtyp == tcp.RoutTypResponse {
					if v, ok := ret.cb.Load(askid); ok {
						f := v.(func(*tcp.Socket, *tcp.PackFrame))
						ret.RemoveCallback(askid)
						ret.eventque <- func() { f(s, p) }
						return
					}
				}

				ret.eventque <- func() {
					ret.OnMessageFunc(ret, head, p)
				}
				return
			} else if ptype == tcp.PacketTypRoutErr {
				//TODO:
			}

			fmt.Println("not support packet type:", ptype)
		},
		OnConnStat: func(s *tcp.Socket, ss tcp.SocketStat) {
			ret.isAuth = false
			if ss == tcp.Disconnected {
				ret.eventque <- func() {
					if ret.OnLoginFunc != nil {
						ret.OnLoginFunc(ret, LoginStat_Disconnect)
					}
				}

				ret.Reconnect()

			} else {
				fmt.Println("connected:", remoteAddr)
				if ret.AuthFunc != nil {
					go func() {
						stat := LoginStat_Success
						isauth := ret.AuthFunc()
						if !isauth {
							stat = LoginStat_Fail
						}
						ret.eventque <- func() {
							if ret.OnLoginFunc != nil {
								ret.OnLoginFunc(ret, stat)
							}
						}
					}()
				}
			}
		},
	})
	ret.Client = c
	return ret
}
