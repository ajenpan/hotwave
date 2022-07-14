package tcpgate

import (
	"hotwave/node/transport"
	"hotwave/node/transport/tcp"
)

type tcpUserSocket struct {
	*tcp.TcpSocket
	uid int64
}

func (s *tcpUserSocket) Send(msg *transport.Message) error {

	//todo:

	return s.TcpSocket.Send(msg)
}
