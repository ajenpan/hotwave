package tcpsvr

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	log "hotwave/logger"
)

type OnMessageFunc func(*Socket, *Packet)
type OnConnStatFunc func(*Socket, SocketStat)

type SocketStat = int32

const (
	SocketStatConnected    SocketStat = 1
	SocketStatDisconnected SocketStat = 2
)

func SocketStatString(s SocketStat) string {
	switch s {
	case SocketStatConnected:
		return "connected"
	case SocketStatDisconnected:
		return "disconnected"
	}
	return "unknown"
}

type SocketOptions struct {
}

type SocketOption func(*SocketOptions)

var sid int64 = 0

func NewSocketID() string {
	return fmt.Sprintf("%d_%d", atomic.AddInt64(&sid, 1), time.Now().Unix())
}

func NewSocket(conn net.Conn, opts ...SocketOption) *Socket {
	ret := &Socket{
		id:      NewSocketID(),
		conn:    conn,
		timeOut: 10 * time.Second,
		// timeOut:  0,
		chSend:   make(chan *Packet, 10),
		chClosed: make(chan bool),
		state:    SocketStatConnected,
	}
	return ret
}

type Socket struct {
	sync.RWMutex // export

	conn     net.Conn   // low-level conn fd
	state    SocketStat // current state
	id       string
	chSend   chan *Packet // push message queue
	chClosed chan bool
	timeOut  time.Duration
	Meta     sync.Map
}

func (s *Socket) ID() string {
	return s.id
}

func (s *Socket) Send(p *Packet) error {
	if atomic.LoadInt32(&s.state) == SocketStatDisconnected {
		return errors.New("sendPacket failed, the socket is disconnected")
	}
	if p.Typ == 0 {
		return fmt.Errorf("packet typ is 0")
	}
	s.chSend <- p
	return nil
}

func (s *Socket) Recv() (*Packet, error) {
	p := &Packet{}
	if err := s.readPacket(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Socket) Close() {
	if s == nil {
		return
	}
	stat := atomic.SwapInt32(&s.state, SocketStatDisconnected)
	if stat == SocketStatDisconnected {
		return
	}

	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	close(s.chSend)
	close(s.chClosed)
}

// returns the remote network address.
func (s *Socket) RemoteAddr() string {
	if s == nil {
		return ""
	}
	return s.conn.RemoteAddr().String()
}

func (s *Socket) LocalAddr() string {
	if s == nil {
		return ""
	}
	return s.conn.LocalAddr().String()
}

//retrun socket work status
func (s *Socket) Status() SocketStat {
	if s == nil {
		return SocketStatDisconnected
	}
	return s.state
}

// String, implementation for Stringer interface
func (s *Socket) String() string {
	return fmt.Sprintf("id:%s, remoteaddr:%s", s.ID(), s.conn.RemoteAddr().String())
}

func (s *Socket) writeWork() {
	for p := range s.chSend {
		err := s.writePacket(p)
		if err != nil {
			log.Warn(err)
		}
	}
}

func (s *Socket) SetMeta(k string, v interface{}) {
	s.Meta.Store(k, v)
}

func (s *Socket) GetMeta(k string) (interface{}, bool) {
	return s.Meta.Load(k)
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

func (s *Socket) readPacket(p *Packet) error {
	if atomic.LoadInt32(&s.state) == SocketStatDisconnected {
		return errors.New("recv packet failed, the socket is disconnected")
	}

	var err error

	headRaw := make([]byte, p.HeadLen())

	if s.timeOut > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.timeOut))
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
		p.Raw = make([]byte, p.RawLen)
		_, err = io.ReadFull(s.conn, p.Raw)
		return err
	}
	return nil
}

func (s *Socket) writePacket(p *Packet) error {
	if atomic.LoadInt32(&s.state) == SocketStatDisconnected {
		return errors.New("writePacket failed, the socket is disconnected")
	}
	var err error
	p.RawLen = int32(len(p.Raw))

	if p.RawLen >= MaxPacketSize {
		return ErrPacketSizeExcced
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
