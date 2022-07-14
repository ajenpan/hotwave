package tcp

import (
	"fmt"
	"io"
	"net"
	"time"

	"hotwave/node/transport"
)

type TcpSocketOptions struct {
	Timeout time.Duration
}

type TcpSocketOption func(*TcpSocketOptions)

func NewTcpSocket(conn net.Conn, opts TcpSocketOptions) *TcpSocket {
	ret := &TcpSocket{
		TcpSocketOptions: opts,
		conn:             conn,
	}
	return ret
}

type TcpSocket struct {
	TcpSocketOptions
	conn net.Conn // low-level conn fd
}

func (s *TcpSocket) SendError(err error) {
	if err == nil {
		return
	}
	s.writePacket(&Packet{
		PacketHead: PacketHead{
			Typ: PacketTypeError,
		},
		Raw: []byte(err.Error()),
	})
}

func (s *TcpSocket) Send(msg *transport.Message) error {
	raw, err := msg.Encode()
	if err != nil {
		return err
	}
	return s.writePacket(&Packet{
		PacketHead: PacketHead{
			Typ: PacketTypePacket,
		},
		Raw: raw,
	})
}

func (s *TcpSocket) Recv(msg *transport.Message) error {
	p := &Packet{}
	err := s.readPacket(p)
	if err != nil {
		return err
	}

	if p.Typ == PacketTypePacket {
		return msg.Decode(p.Raw)
	} else if p.Typ == PacketTypePing {
		p.Typ = PacketTypePong
		s.writePacket(p)
	} else if p.Typ == PacketTypeError {
		s.OnRecvError(p.Raw)
	}
	return s.Recv(msg)
}

func (s *TcpSocket) OnRecvError(raw []byte) {
	// TODO:
}

func (s *TcpSocket) Close() error {
	if s == nil {
		return nil
	}

	if s.conn != nil {
		err := s.conn.Close()
		return err
	}
	return nil
}

// returns the remote network address.
func (s *TcpSocket) RemoteAddr() string {
	if s == nil {
		return ""
	}
	return s.conn.RemoteAddr().String()
}

func (s *TcpSocket) LocalAddr() string {
	if s == nil {
		return ""
	}
	return s.conn.LocalAddr().String()
}

func (s *TcpSocket) String() string {
	return fmt.Sprintf("%s-%s", s.RemoteAddr(), s.LocalAddr())
}

func writeAll(conn net.Conn, raw []byte) (int, error) {
	writelen := 0
	rawSize := len(raw)

	for writelen < rawSize {
		n, err := conn.Write(raw[writelen:])
		writelen += n
		if err != nil {
			return writelen, err
		}
	}

	return writelen, nil
}

func (s *TcpSocket) readPacket(p *Packet) error {
	var err error

	headRaw := make([]byte, p.HeadLen())

	if s.Timeout > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.Timeout))
	}

	_, err = io.ReadFull(s.conn, headRaw)
	if err != nil {
		return err
	}

	err = p.PacketHead.Decode(headRaw)
	if err != nil {
		return err
	}

	if p.RawLen > 0 {
		//TODO: use buffer pool impove this performance
		p.Raw = make([]byte, p.RawLen)
		_, err = io.ReadFull(s.conn, p.Raw)
		return err
	}
	return nil
}

func (s *TcpSocket) writePacket(p *Packet) error {
	var err error
	p.RawLen = int32(len(p.Raw))

	if p.RawLen >= MaxPacketSize {
		return ErrPacketSizeExcced
	}

	if s.Timeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.Timeout))
	}

	headRaw := make([]byte, p.HeadLen())
	if err := p.PacketHead.Encode(headRaw); err != nil {
		return err
	}

	_, err = writeAll(s.conn, headRaw)
	if err != nil {
		return err
	}

	if p.RawLen > 0 {
		_, err = writeAll(s.conn, p.Raw)
		if err != nil {
			return err
		}
	}
	return nil
}
