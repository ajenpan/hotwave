package handle

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/proto"

	msg "hotwave/service/route/proto"
	"hotwave/service/route/transport/tcp"
	"hotwave/utils/calltable"
)

func NewRouter() (*Router, error) {
	ret := &Router{
		uid: 9999,
	}
	ct := calltable.ExtractAsyncMethodByMsgID(msg.File_proto_route_proto.Messages(), ret)
	if ct == nil {
		return nil, fmt.Errorf("ExtractProtoFile failed")
	}
	ret.ct = ct
	return ret, nil
}

type Router struct {
	user    sync.Map
	session sync.Map

	uid uint64

	PublicKey *rsa.PublicKey
	ct        *calltable.CallTable[int]
}

type UserInfo struct {
	UID      uint64
	UserName string
	Group    int64
	Role     string

	LoginAt time.Time
}

const uinfoKey string = "uinfo"
const errcntKey string = "errcnt"

type tcpSocketKeyT struct{}
type tcpPacketKeyT struct{}

var tcpSocketKey = tcpSocketKeyT{}
var tcpPacketKey = tcpPacketKeyT{}

func VerifyToken(pk *rsa.PublicKey, tokenRaw string) (uint64, string, string, error) {
	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (interface{}, error) {
		return pk, nil
	})
	if err != nil {
		return 0, "", "", err
	}
	if !token.Valid {
		return 0, "", "", fmt.Errorf("invalid token")
	}
	uname := claims["sub"]
	uidstr := claims["uid"]
	role := claims["role"]
	uid, err := strconv.ParseUint(uidstr.(string), 10, 64)
	return uid, uname.(string), role.(string), err
}

func GetSocketUserInfo(s *tcp.Socket) *UserInfo {
	if s == nil {
		return nil
	}
	if v, ok := s.MetaLoad(uinfoKey); ok {
		return v.(*UserInfo)
	}
	return nil
}

func addSocketErrCnt(s *tcp.Socket) int {
	if v, ok := s.MetaLoad(errcntKey); ok {
		cnt := v.(int)
		cnt++
		s.MetaStore(errcntKey, cnt)
		return cnt
	}
	s.MetaStore(errcntKey, 1)
	return 1
}

func dealSocketErrCnt(s *tcp.Socket) {
	cnt := addSocketErrCnt(s)
	fmt.Printf("socket:%v, uid:%v, errcnt:%v", s.ID(), s.UID(), cnt)
}

func (r *Router) OnMessage(s *tcp.Socket, p *tcp.PackFrame) {
	ptype := p.GetType()
	if ptype == tcp.PacketTypRoutDeliver {
		head := tcp.RoutDeliverHead(p.Head)
		targetid := head.GetTargetUID()
		if targetid == 0 || targetid == r.uid {
			r.OnCall(s, head, p)
			return
		}
		var err error

		suid := s.UID()
		if suid != 0 {
			tsocket := r.GetSocketByUID(targetid)
			if tsocket != nil {
				head.SetSrouceUID(suid)
				err = tsocket.SendPacket(p)
			} else {
				err = fmt.Errorf("targetid:%d not found", targetid)
			}
		} else {
			err = fmt.Errorf("not login")
		}

		if err != nil {
			p.SetType(tcp.PacketTypRoutErr)
			p.Body = []byte(err.Error())
			err = s.SendPacket(p)
			if err != nil {
				fmt.Println("send error failed:", err)
			}
		}
	}
}

func (r *Router) OnConn(s *tcp.Socket, stat tcp.SocketStat) {
	fmt.Println("OnConn", s.ID(), ",stat:", int(stat))

	if stat == tcp.Disconnected {
		r.session.Delete(s.ID())
		uinfo := GetSocketUserInfo(s)
		if uinfo != nil {
			r.user.Delete(uinfo.UID)
		}
	} else {
		r.session.Store(s.ID(), s)
	}
}

func (r *Router) SendMessage(s *tcp.Socket, askid uint32, msgtyp uint8, m proto.Message) error {
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
	head.SetSrouceUID(r.uid)
	head.SetTargetUID(s.UID())
	head.SetAskID(askid)
	head.SetMsgTyp(msgtyp)

	p := &tcp.PackFrame{
		Head: head,
		Body: raw,
	}
	p.SetType(tcp.PacketTypRoutDeliver)
	return s.SendPacket(p)
}

func (r *Router) GetSocketByUID(uid uint64) *tcp.Socket {
	if v, ok := r.user.Load(uid); ok {
		return v.(*tcp.Socket)
	}
	return nil
}

func (r *Router) OnCall(s *tcp.Socket, head tcp.RoutDeliverHead, p *tcp.PackFrame) {
	var err error
	msgid := int(head.GetMsgID())
	askid := head.GetAskID()
	method := r.ct.Get(msgid)
	if method == nil {
		fmt.Println("not found method,msgid:", msgid)
		dealSocketErrCnt(s)
		return
	}

	reqRaw := method.NewRequest()
	if reqRaw == nil {
		fmt.Println("not found request,msgid:", msgid)
		return
	}

	req := reqRaw.(proto.Message)
	err = proto.Unmarshal(p.Body, req)

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.WithValue(context.Background(), tcpSocketKey, s)
	ctx = context.WithValue(ctx, tcpPacketKey, p)

	result := method.Call(r, ctx, req)

	if len(result) != 2 {
		return
	}
	respI := result[0].Interface()
	if respI != nil {
		resp, ok := respI.(proto.Message)
		if !ok {
			return
		}
		respMsgTyp := head.GetMsgTyp()
		if respMsgTyp == tcp.RoutTypRequest {
			respMsgTyp = tcp.RoutTypResponse
		}
		r.SendMessage(s, askid, respMsgTyp, resp)
		fmt.Printf("oncall sid:%v,uid:%v,msgid:%v,askid:%v,req:%v,resp:%v\n", s.ID(), s.UID(), msgid, askid, req, resp)
		return
	}

	resperrI := result[1].Interface()
	if resperrI != nil {
		resperr, ok := resperrI.(error)
		if !ok {
			return
		}

		fmt.Println("resperr:", resperr)
		dealSocketErrCnt(s)
	}
}

func GetSocketFromCtx(ctx context.Context) *tcp.Socket {
	if v, ok := ctx.Value(tcpSocketKey).(*tcp.Socket); ok {
		return v
	}
	return nil
}

func GetPacketFromCtx(ctx context.Context) *tcp.PackFrame {
	if v, ok := ctx.Value(tcpPacketKey).(*tcp.PackFrame); ok {
		return v
	}
	return nil
}

func (r *Router) OnLoginRequest(ctx context.Context, req *msg.LoginRequest) (*msg.LoginResponse, error) {
	resp := &msg.LoginResponse{
		Errcode: msg.LoginResponse_unkown_err,
	}

	uid, uname, role, err := VerifyToken(r.PublicKey, req.Token)
	if err != nil {
		resp.Errcode = msg.LoginResponse_invalid_token
		return resp, nil
	}

	uinfo := &UserInfo{
		UID:      uid,
		UserName: uname,
		Role:     role,
	}
	s := GetSocketFromCtx(ctx)

	s.MetaStore(uinfoKey, uinfo)
	s.SetUID(uid)
	r.user.Store(uid, s)
	resp.Errcode = msg.LoginResponse_ok
	return resp, nil
}

func (r *Router) OnAdminLoginRequest(ctx context.Context, req *msg.AdminLoginRequest) (*msg.AdminLoginResponse, error) {
	resp := &msg.AdminLoginResponse{
		Errcode: msg.AdminLoginResponse_unkown_err,
	}

	uinfo := &UserInfo{
		UID:      uint64(req.Uid),
		UserName: req.Uname,
		Role:     req.Role,
	}

	s := GetSocketFromCtx(ctx)

	s.MetaStore(uinfoKey, uinfo)
	s.SetUID(uinfo.UID)
	r.user.Store(uinfo.UID, s)

	resp.Errcode = msg.AdminLoginResponse_ok

	return resp, nil
}

func (r *Router) OnEchoRequest(ctx context.Context, req *msg.EchoRequest) (*msg.EchoResponse, error) {
	resp := &msg.EchoResponse{
		Msg: req.Msg,
	}
	return resp, nil
}
