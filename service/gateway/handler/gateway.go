package handler

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	log "hotwave/logger"
	"hotwave/service/gateway/auth"
	protocal "hotwave/service/gateway/proto"
	"hotwave/session"
	"hotwave/transport/tcp"
)

type Gateway struct {
	protocal.UnimplementedGatewayServer

	users sync.Map //map[int64]usersession

	Servers sync.Map
	Nodes   sync.Map

	Authc *auth.LocalAuth
}

func (g *Gateway) OnClientConnStat(socket *tcp.Socket, status tcp.SocketStat) {
	switch status {
	case tcp.SocketStatConnected:

	case tcp.SocketStatDisconnected:
		rawUser, has := socket.Meta.Load("user")
		if !has {
			return
		}
		user := rawUser.(*auth.UserSession)
		g.users.Delete(user.UID())
		socket.Meta.Delete("user")
	}
}

func (g *Gateway) OnClientMessage(session *tcp.Socket, raw []byte) {
	msg := &protocal.ClientMessage{}
	err := proto.Unmarshal(raw, msg)
	if err != nil {
		return
	}

	fullName := protoreflect.FullName(msg.Name)
	serverName := string(fullName.Parent())

	rawUser, has := session.Meta.Load("user")
	if !has {
		if fullName == "gateway.LoginRequest" {
			in := &protocal.LoginRequest{}
			err := proto.Unmarshal(msg.Body, in)
			if err != nil {
				return
			}
			g.OnLoginRequestV1(session, in)
		}
	} else {
		user := rawUser.(*auth.UserSession)
		if serverName == "gateway" {

		} else {
			//TODO : batter router
			s, has := g.Servers.Load(serverName)
			if has {
				s := s.(*grpcSvrSession)
				warp := &protocal.ToServerMessage{
					Uid:  user.UID(),
					Name: msg.Name,
					Data: msg.Body,
				}
				s.Send(warp)
			}
		}
	}
}

func (g *Gateway) OnGateMessage(s session.Session, msg *protocal.ToClientMessage) {
	rawUser, has := g.users.Load(msg.Uid)
	if !has {
		log.Warn("user not found: ", msg.Uid)
		return
	}
	user := rawUser.(*auth.UserSession)
	wrap := &protocal.ServerMessage{
		Name: msg.Name,
		Body: msg.Data,
	}
	raw, err := proto.Marshal(wrap)
	if err != nil {
		log.Error(err)
		return
	}
	if err = user.Socket.Send(raw); err != nil {
		log.Error("socket send err: ", err)
	}
}

type grpcSvrSession struct {
	conn protocal.Gateway_ProxyServer

	nodename string
	nodeid   string
	id       string
}

func (s *grpcSvrSession) ID() string {
	return s.id
}

func (s *grpcSvrSession) String() string {
	return "router-grpc-conn"
}

func (s *grpcSvrSession) Close() error {
	return nil
}

func (s *grpcSvrSession) Send(msg interface{}) error {
	switch msg := msg.(type) {
	case *protocal.ToServerMessage:
		if s.conn != nil {
			return s.conn.Send(msg)
		}
	}
	return fmt.Errorf("unknown message type")
}

func (rs *Gateway) Proxy(s protocal.Gateway_ProxyServer) error {
	md, ok := metadata.FromIncomingContext(s.Context())
	if !ok {
		return fmt.Errorf("no metadata")
	}

	nodename := strings.Join(md.Get("nodename"), "")
	nodeid := strings.Join(md.Get("nodeid"), "")

	id := nodename + "-" + nodeid

	conn := &grpcSvrSession{
		conn:     s,
		id:       id,
		nodename: nodename,
		nodeid:   nodeid,
	}
	log.Info("router-grpc-conn connect: ", id)
	rs.Nodes.Store(nodeid, conn)
	rs.Servers.Store(nodename, conn)
	var recvErr error

	for {
		in, err := s.Recv()
		if err != nil {
			if err != io.EOF {
				recvErr = err
			}
			break
		}
		rs.OnGateMessage(conn, in)
	}

	log.Info("router-grpc-conn disconnect: ", id)

	rs.Nodes.Delete(nodeid)
	rs.Nodes.Delete(nodename)
	return recvErr
}
