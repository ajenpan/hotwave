package grpcgate

import (
	"hotwave/node/transport"
	gwProto "hotwave/service/gateway/proto"
)

type UserSocket struct {
	client         *ClientGate
	UID            int64
	RemoteSocketID string
}

func (s *UserSocket) Send(msg *transport.Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	warp := &gwProto.ToUserMessage{
		ToUid:      s.UID,
		ToSocketid: s.RemoteSocketID,
		Data:       data,
	}

	return s.client.sendMessage(warp)
}
