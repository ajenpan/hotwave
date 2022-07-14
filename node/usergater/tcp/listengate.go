package tcpgate

import (
	"hotwave/marshal"
	"hotwave/node/transport"
	"hotwave/node/transport/tcp"
	"hotwave/node/usergater"
)

type TcpListenGater struct {
	listen tcp.TcpListener
}

func (gater *TcpListenGater) Accept(fn func(usergater.UserSocket)) error {
	gater.listen.Accept(func(socket *tcp.TcpSocket) {
		for {
			msg := &transport.Message{
				Marshaler: &marshal.ProtoMarshaler{},
			}
			err := socket.Recv(msg)
			if err != nil {
				break
			}
		}
	})
	return nil
}

func (gater *TcpListenGater) Close() error {
	return gater.listen.Close()
}
