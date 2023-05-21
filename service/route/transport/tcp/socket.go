package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type SocketStat int32

const (
	Disconnected SocketStat = iota
	Connected    SocketStat = iota
)

type OnMessageFunc func(*Socket, *PackFrame)
type OnConnStatFunc func(*Socket, SocketStat)
type NewIDFunc func() string

type SocketOptions struct {
	ID string
}

type SocketOption func(*SocketOptions)

var staticIdx uint64

func nextID() string {
	idx := atomic.AddUint64(&staticIdx, 1)
	if idx == 0 {
		idx = atomic.AddUint64(&staticIdx, 1)
	}
	return fmt.Sprintf("tcp_%v_%v", idx, time.Now().Unix())
}

func NewSocket(conn net.Conn, opts SocketOptions) *Socket {
	if opts.ID == "" {
		opts.ID = nextID()
	}

	ret := &Socket{
		id:       opts.ID,
		conn:     conn,
		timeOut:  120 * time.Second,
		chSend:   make(chan *PackFrame, 10),
		chClosed: make(chan bool),
		state:    Connected,
		packetpool: sync.Pool{
			New: func() interface{} {
				return &PackFrame{}
			},
		},
	}
	return ret
}

type Socket struct {
	conn     net.Conn   // low-level conn fd
	state    SocketStat // current state
	id       string
	uid      uint64
	chSend   chan *PackFrame // push message queue
	chClosed chan bool

	timeOut time.Duration

	lastSendAt uint64
	lastRecvAt uint64

	store sync.Map

	askididx uint32

	packetpool sync.Pool
}

func (s *Socket) GetAskID() uint32 {
	ret := atomic.AddUint32(&s.askididx, 1)
	if ret == 0 {
		ret = atomic.AddUint32(&s.askididx, 1)
	}
	return ret
}

func (s *Socket) ID() string {
	return s.id
}

func (s *Socket) UID() uint64 {
	return atomic.LoadUint64(&s.uid)
}

func (s *Socket) SetUID(uid uint64) {
	atomic.StoreUint64(&s.uid, uid)
}

func (s *Socket) SendPacket(p *PackFrame) error {
	if atomic.LoadInt32((*int32)(&s.state)) == int32(Disconnected) {
		return ErrDisconn
	}
	if len(p.Body) > MaxPacketSize {
		return ErrPacketSizeExcced
	}
	s.chSend <- p
	return nil
}

func (s *Socket) Close() {
	if s == nil {
		return
	}
	stat := atomic.SwapInt32((*int32)(&s.state), int32(Disconnected))
	if stat == int32(Disconnected) {
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
func (s *Socket) RemoteAddr() net.Addr {
	if s == nil {
		return nil
	}
	return s.conn.RemoteAddr()
}

func (s *Socket) LocalAddr() net.Addr {
	if s == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

// retrun socket work status
func (s *Socket) Status() SocketStat {
	if s == nil {
		return Disconnected
	}
	return SocketStat(atomic.LoadInt32((*int32)(&s.state)))
}

func (s *Socket) writeWork() {
	for p := range s.chSend {
		s.writePacket(p)
	}
}

func (s *Socket) newPacket() *PackFrame {
	return &PackFrame{}
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

func (s *Socket) readPacket(p *PackFrame) error {
	if s.Status() == Disconnected {
		return ErrDisconn
	}

	var err error

	if s.timeOut > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.timeOut))
	}

	_, err = io.ReadFull(s.conn, p.PacketHead[:])
	if err != nil {
		return err
	}

	//TODO: use buffer pool impove this performance

	headlen := uint32(p.GetHeadLen())
	if headlen > 0 {
		p.Head = make([]byte, headlen)
		_, err = io.ReadFull(s.conn, p.Head)
		if err != nil {
			return err
		}
	}
	bodylen := p.GetBodyLen()
	if bodylen > 0 {
		p.Body = make([]byte, bodylen)
		_, err = io.ReadFull(s.conn, p.Body)
		if err != nil {
			return err
		}
	}

	atomic.StoreUint64(&s.lastRecvAt, uint64(time.Now().Unix()))
	return nil
}

func (s *Socket) writePacket(p *PackFrame) error {
	if s.Status() == Disconnected {
		return ErrDisconn
	}

	var err error

	if len(p.Body) >= MaxPacketSize {
		return ErrPacketSizeExcced
	}

	p.SetHeadLen(uint8(len(p.Head)))
	p.SetBodyLen(uint32(len(p.Body)))

	_, err = writeAll(s.conn, p.PacketHead[:])
	if err != nil {
		return err
	}

	if len(p.Head) > 0 {
		_, err = writeAll(s.conn, p.Head[:])
		if err != nil {
			return err
		}
	}

	if len(p.Body) > 0 {
		_, err = writeAll(s.conn, p.Body)
		if err != nil {
			return err
		}
	}

	atomic.StoreUint64(&s.lastSendAt, uint64(time.Now().Unix()))
	return nil
}

func (m *Socket) MetaLoad(key string) (interface{}, bool) {
	return m.store.Load(key)
}

func (m *Socket) MetaStore(key string, value interface{}) {
	m.store.Store(key, value)
}

func (m *Socket) MetaDelete(key string) {
	m.store.Delete(key)
}
