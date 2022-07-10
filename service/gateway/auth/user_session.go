package auth

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	protocal "hotwave/service/gateway/proto"
	"hotwave/transport/tcp"
)

type UserSession struct {
	uid   int64
	uname string

	Socket *tcp.Socket
}

func (u *UserSession) Send(data interface{}) error {
	if u.Socket == nil {
		return fmt.Errorf("session is nil")
	}
	switch data := data.(type) {
	case proto.Message:
		return u.SendPB(data)
	case []byte:
		return u.Socket.Send(data)
	default:
		return fmt.Errorf("data type %T not support", data)
	}
}

func (u *UserSession) SendPB(msg proto.Message) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	wrap := &protocal.ServerMessage{
		Name: string(proto.MessageName(msg)),
		Body: body,
	}
	raw, err := proto.Marshal(wrap)
	if err != nil {
		return err
	}
	return u.Socket.Send(raw)
}

func (u *UserSession) UID() int64 {
	return u.uid
}
func (u *UserSession) ID() string {
	return u.uname
}
func (u *UserSession) String() string {
	return "UserSession"
}
func (u *UserSession) Close() error {
	u.Socket.Close()
	return nil
}
