package handler

import (
	"github.com/google/uuid"

	"hotwave/service/gateway/auth"
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
)

func (g *Gateway) OnLoginGateRequest(socket transport.Session, in *protocal.LoginGateRequest) error {
	var user *auth.UserInfo
	switch c := in.Checker.(type) {
	case *protocal.LoginGateRequest_Account:
		user = g.AuthClient.AccountAuth(c.Account.Account, c.Account.Passwd)
	case *protocal.LoginGateRequest_Session:
	case *protocal.LoginGateRequest_Jwt:
		user = g.AuthClient.TokenAuth(c.Jwt)
	}

	out := &protocal.LoginGateResponse{}

	if user == nil {
		out.Flag = protocal.LoginGateResponse_UnknowError
		return socket.Send(out)
	}

	socket.MetaStore("userinfo", user)
	g.users.Store(user.Uid, user)

	out.Sessionid = uuid.NewString()

	err := socket.Send(out)

	g.pushishEvent(&protocal.UserConnect{
		Uid: user.Uid,
	})
	return err
}

func (g *Gateway) OnEchoRequest(socket transport.Session, in *protocal.EchoRequest) error {
	out := &protocal.EchoResponse{
		Data: in.Data,
	}
	return socket.Send(out)
}
