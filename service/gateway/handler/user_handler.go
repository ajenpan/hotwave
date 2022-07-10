package handler

import (
	"github.com/google/uuid"

	"hotwave/service/gateway/auth"
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport/tcp"
)

func (g *Gateway) OnLoginRequestV1(socket *tcp.Socket, in *protocal.LoginRequest) {
	var user *auth.UserSession
	switch c := in.Checker.(type) {
	case *protocal.LoginRequest_Account:
		user = g.Authc.AccountAuth(c.Account.Account, c.Account.Passwd)
	case *protocal.LoginRequest_Session:
		// user = g.authc.SeesionLogin(c.Session.Account, c.Session.Session)
	case *protocal.LoginRequest_Jwt:
		user = g.Authc.TokenAuth(c.Jwt)
	}

	if user == nil {
		return
	}

	user.Socket = socket
	socket.Meta.Store("user", user)

	g.users.Store(user.UID(), user)

	out := &protocal.LoginResponse{
		Sessionid: uuid.NewString(),
	}
	user.SendPB(out)
}

func (g *Gateway) OnEcho(socket *auth.UserSession, in *protocal.EchoRequest) {
	out := &protocal.EchoResponse{
		Data: in.Data,
	}
	socket.Send(out)
}
