package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"hotwave/transport"
)

type OnMessageFunc func(transport.Session, []byte)
type OnConnStatFunc func(transport.Session, transport.SessionStat)

type SocketStat = transport.SessionStat

type SocketOptions struct {
	ID string
}

type SocketOption func(*SocketOptions)

var sid int64 = 0

func newSocketID() string {
	return fmt.Sprintf("%d_%d", atomic.AddInt64(&sid, 1), time.Now().Unix())
}

func NewSocket(conn net.Conn, opts SocketOptions) *Socket {
	if opts.ID == "" {
		opts.ID = newSocketID()
	}

	ret := &Socket{
		id:       opts.ID,
		conn:     conn,
		timeOut:  120 * time.Second,
		chSend:   make(chan *Packet, 10),
		chClosed: make(chan bool),
		state:    transport.Connected,
	}
	return ret
}

type Socket struct {
	transport.SyncMapSocketMeta

	// sync.RWMutex // export

	conn     net.Conn              // low-level conn fd
	state    transport.SessionStat // current state
	id       string
	chSend   chan *Packet // push message queue
	chClosed chan bool
	timeOut  time.Duration
}

func (s *Socket) ID() string {
	return s.id
}

func (s *Socket) SendPacket(p *Packet) error {
	if atomic.LoadInt32((*int32)(&s.state)) == int32(transport.Disconnected) {
		return errors.New("sendPacket failed, the socket is disconnected")
	}
	if p.Typ == 0 {
		return fmt.Errorf("packet typ is 0")
	}
	s.chSend <- p
	return nil
}

func (s *Socket) SendError(err error) {
	if err == nil {
		return
	}
	s.SendPacket(&Packet{
		PacketHead: PacketHead{
			Typ: PacketTypeError,
		},
		Raw: []byte(err.Error()),
	})
}

func (s *Socket) Send(iraw interface{}) error {
	raw, ok := iraw.([]byte)
	if !ok {
		return errors.New("send data must be []byte")
	}
	return s.SendPacket(&Packet{
		PacketHead: PacketHead{
			Typ: PacketTypePacket,
		},
		Raw: raw,
	})
}

func (s *Socket) Close() {
	if s == nil {
		return
	}
	stat := atomic.SwapInt32((*int32)(&s.state), int32(transport.Disconnected))
	if stat == int32(transport.Disconnected) {
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
		return transport.Disconnected
	}
	return transport.SessionStat(atomic.LoadInt32((*int32)(&s.state)))
}

// String, implementation for Stringer interface
func (s *Socket) String() string {
	return fmt.Sprintf("id:%s, remoteaddr:%s", s.ID(), s.conn.RemoteAddr().String())
}

func (s *Socket) writeWork() {
	for p := range s.chSend {
		err := s.writePacket(p)
		if err != nil {
			// log.Warn(err)
			fmt.Println(err)
		}
	}
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
	if s.Status() == transport.Disconnected {
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
		//TODO: use buffer pool impove this performance
		p.Raw = make([]byte, p.RawLen)
		_, err = io.ReadFull(s.conn, p.Raw)
		return err
	}
	return nil
}

func (s *Socket) writePacket(p *Packet) error {
	if s.Status() == transport.Disconnected {
		return errors.New("recv packet failed, the socket is disconnected")
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