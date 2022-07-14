package tcp

import (
	"net"
	"time"

	mnet "hotwave/utils/net"
)

type TcpListenerOption func(*TcpListenerOptions)

type TcpListenerOptions struct {
	Address string
	Timeout time.Duration
}

func NewTcpListener(opts TcpListenerOptions) (*TcpListener, error) {
	ret := &TcpListener{
		opts: opts,
		die:  make(chan bool),
	}

	fn := func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	}

	l, err := mnet.Listen(opts.Address, fn)
	if err != nil {
		return nil, err
	}
	ret.listener = l

	return ret, nil
}

type TcpListener struct {
	opts     TcpListenerOptions
	die      chan bool
	listener net.Listener
}

func (s *TcpListener) Close() error {
	return s.listener.Close()
}

func (s *TcpListener) Accept(fn func(*TcpSocket)) error {
	var tempDelay time.Duration = 0

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		socket := NewTcpSocket(conn, TcpSocketOptions{
			Timeout: s.opts.Timeout,
		})

		//recv syn
		p := &Packet{}
		if err := socket.readPacket(p); err != nil {
			continue
		}
		if p.Typ != PacketTypeSyn {
			continue
		}

		// send ack
		// p.Raw = []byte("")
		p.Typ = PacketTypeAck
		if err := socket.writePacket(p); err != nil {
			continue
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					socket.Close()
				}
			}()
			fn(socket)
		}()
	}
}

func (s *TcpListener) Addr() string {
	return s.listener.Addr().String()
}
