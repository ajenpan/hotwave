package handler

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"hotwave/event"
	log "hotwave/logger"
	"hotwave/service/common"
	"hotwave/service/gateway/auth"
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
	"hotwave/utils/calltable"
)

type GatewayOption struct {
	AuthClient auth.Auth
	Publisher  event.Publisher
}

func NewGateway(opt GatewayOption) (*Gateway, error) {
	g := &Gateway{
		GatewayOption: opt,
	}
	g.CT = calltable.ExtractAsyncMethod(protocal.File_service_gateway_proto_gateway_client_proto.Messages(), g)
	return g, nil
}

type Gateway struct {
	protocal.UnimplementedGatewayServer
	GatewayOption

	users     sync.Map //map[int64]usersession
	sockets   sync.Map
	Servers   sync.Map
	Allowlist sync.Map
	CT        *calltable.CallTable
}

func (g *Gateway) AddAllowlistByMsg(msg proto.Message, msgs ...proto.Message) {
	g.Allowlist.Store(string(proto.MessageName(msg)), true)
	for _, v := range msgs {
		g.Allowlist.Store(string(proto.MessageName(v)), true)
	}
}

func (g *Gateway) AddAllowlistByName(msgname string, msgnames ...string) {
	g.Allowlist.Store(msgname, true)
	for _, v := range msgnames {
		g.Allowlist.Store(v, true)
	}
}

func (g *Gateway) pushishEvent(msg proto.Message) {
	if g.Publisher == nil {
		return
	}

	log.Info("[Gateway.Event] ", string(proto.MessageName(msg)))

	data, _ := proto.Marshal(msg)
	g.Publisher.Publish(&event.Event{
		Topic:     string(proto.MessageName(msg)),
		FromNode:  "gateway",
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

type SocketSendWarper struct {
	transport.Session
}

func NewSendWarper(s transport.Session) transport.Session {
	return &SocketSendWarper{
		Session: s,
	}
}

func (u *SocketSendWarper) String() string {
	return "SocketSendWarper"
}

func (u *SocketSendWarper) Send(data interface{}) error {
	if u.Session == nil {
		return fmt.Errorf("session is nil")
	}
	switch data := data.(type) {
	case []byte:
		return u.Session.Send(data)
	case proto.Message:
		log.Info("use sendpb install of send proto")
		return u.SendPB(data)
	default:
		return fmt.Errorf("data type %T not support", data)
	}
}

func (u *SocketSendWarper) SendPB(msg proto.Message) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	wrap := &protocal.GateMessage{
		Name: string(proto.MessageName(msg)),
		Body: body,
	}
	raw, err := proto.Marshal(wrap)
	if err != nil {
		return err
	}
	return u.Session.Send(raw)
}

func (g *Gateway) OnGateConnStat(socket transport.Session, status transport.SessionStat) {
	switch status {
	case transport.Connected:
		g.sockets.Store(socket.ID(), NewSendWarper(socket))
	case transport.Disconnected:
		g.sockets.Delete(socket.ID())
		rawUser, has := socket.MetaLoad("userinfo")
		if !has {
			return
		}
		user := rawUser.(*auth.UserInfo)
		g.users.Delete(user.Uid)
		socket.MetaDelete("userinfo")
		g.pushishEvent(&protocal.UserDisconnect{
			Uid: user.Uid,
		})
	}
}

func (g *Gateway) OnGateMessage(session transport.Session, iraw interface{}) {
	raw, ok := iraw.([]byte)
	if !ok {
		return
	}
	if wrap, has := g.sockets.Load(session.ID()); !has {
		return
	} else {
		session = wrap.(transport.Session)
	}

	msg := &protocal.GateMessage{}
	err := proto.Unmarshal(raw, msg)
	if err != nil {
		return
	}
	name := (protoreflect.FullName(msg.Name))
	fullName := string(name)
	serverName := msg.Server

	rawUser, hasUserInfo := session.MetaLoad("user")
	var uid = int64(0)
	if !hasUserInfo {
		_, in := g.Allowlist.Load(fullName)
		if !in {
			log.Warn("user not login and not in allowlist: ", fullName)
			return
		}
	} else {
		user := rawUser.(*auth.UserInfo)
		uid = user.Uid
	}

	if serverName == "gateway" {
		if mthod := g.CT.Get(fullName); mthod != nil {
			err := common.CallHelper(mthod, session, msg.Body)
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		//TODO:
		s, has := g.Servers.Load(serverName)
		if has {
			s := s.(*grpcSvrSession)
			warp := &protocal.ToServerMessage{
				FromUid:      uid,
				FromSocketid: session.ID(),
				Name:         msg.Name,
				Data:         msg.Body,
			}
			s.Send(warp)
		}
	}
}

func (g *Gateway) ToUserMessage(s transport.Session, msg *protocal.ToUserMessage) {
	var isocket interface{}

	if msg.ToUid != 0 {
		isocket, _ = g.users.Load(msg.ToUid)
	}
	if isocket == nil {
		isocket, _ = g.sockets.Load(msg.ToSocketid)
	}

	// socket still nil
	if isocket == nil {
		log.Warn("user not found: ", msg.ToUid, msg.ToSocketid)
		return
	}

	wrap := &protocal.GateMessage{
		Name: msg.Name,
		Body: msg.Data,
	}

	raw, err := proto.Marshal(wrap)
	if err != nil {
		log.Error(err)
		return
	}

	if err = isocket.(transport.Session).Send(raw); err != nil {
		log.Error("socket send err: ", err)
		//todo: report send failed
	}
}

type grpcSvrSession struct {
	conn protocal.Gateway_ProxyServer
	transport.SessionMeta

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

func (s *grpcSvrSession) Close() {

}
func (s *grpcSvrSession) RemoteAddr() string {
	return ""
}
func (s *grpcSvrSession) LocalAddr() string {
	return ""
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
		return fmt.Errorf("no metadata, nodename and nodeid is required")
	}

	nodename := strings.Join(md.Get("nodename"), "")
	nodeid := strings.Join(md.Get("nodeid"), "")

	if nodeid == "" {
		nodeid = uuid.NewString()
	}

	id := nodename + "-" + nodeid

	conn := &grpcSvrSession{
		conn:     s,
		id:       id,
		nodename: nodename,
		nodeid:   nodeid,
	}

	if _, loaded := rs.Servers.LoadOrStore(nodename, conn); loaded {
		return fmt.Errorf("node already exists")
	}

	log.Info("router-grpc-conn connect: ", id)

	var recvErr error
	for {
		in, err := s.Recv()
		if err != nil {
			if err != io.EOF {
				recvErr = err
			}
			break
		}
		rs.ToUserMessage(conn, in)
	}

	log.Info("router-grpc-conn disconnect: ", id)
	rs.Servers.Delete(nodename)
	return recvErr
}

func (rs *Gateway) AddGateAllowList(ctx context.Context, in *protocal.AddGateAllowListRequest) (*protocal.AddGateAllowListResponse, error) {
	for _, v := range in.Names {
		log.Info("add allowlist: ", v)
		rs.Allowlist.Store(v, true)
	}
	return &protocal.AddGateAllowListResponse{}, nil
}
