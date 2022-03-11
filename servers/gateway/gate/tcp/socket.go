package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"hotwave/servers/gateway/gate/codec"
)

type OnMessageFunc func(*Socket, *Packet)
type OnConnStatFunc func(*Socket, SocketStat)

type SocketStat = int32

const (
	SocketStatConnected    SocketStat = 1
	SocketStatDisconnected SocketStat = 2
)

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
		chSend:  make(chan *Packet, 10),
		state:   SocketStatConnected,
	}
	return ret
}

type Socket struct {
	conn   net.Conn   // low-level conn fd
	state  SocketStat // current state
	id     string
	chSend chan *Packet // push message queue

	timeOut time.Duration
	Meta    sync.Map
}

func (s *Socket) ID() string {
	return s.id
}

func ConverPacket(msg *codec.Message) *Packet {
	packet := &Packet{}
	switch msg.Type {
	case codec.Request:
		packet.Typ = PacketTypeRequest
	case codec.Response:
		packet.Typ = PacketTypeResponse
	case codec.Event:
		fallthrough
	case codec.Async:
		packet.Typ = PacketTypePacket
	default:
	}
	packet.Raw = msg.Body
	packet.PacketHead.RawLen = int32(len(msg.Body))
	return packet
}
func ConverMessage(p *Packet) *codec.Message {
	msg := &codec.Message{}

	switch p.Typ {
	case PacketTypeRequest:
	case PacketTypeResponse:
	case PacketTypePacket:
	}
	msg.Body = p.Raw
	return msg
}

func (a *Socket) Send(p *codec.Message) error {
	if atomic.LoadInt32(&a.state) == SocketStatDisconnected {
		return fmt.Errorf("send packet failed, the socket is disconnected")
	}
	a.chSend <- ConverPacket(p)
	return nil
}

func (a *Socket) Recv() (*Packet, error) {
	if atomic.LoadInt32(&a.state) == SocketStatDisconnected {
		return nil, fmt.Errorf("recv packet failed, the socket is disconnected")
	}

	p := &Packet{}
	if err := a.readPacket(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (a *Socket) Close() {
	s := atomic.SwapInt32(&a.state, SocketStatDisconnected)
	if s == SocketStatDisconnected {
		return
	}

	if a.conn != nil {
		a.conn.Close()
	}
}

// returns the remote network address.
func (a *Socket) RemoteAddr() string {
	return a.conn.RemoteAddr().String()
}

func (a *Socket) LocalAddr() string {
	return a.conn.LocalAddr().String()
}

//retrun socket work status
func (a *Socket) Status() SocketStat {
	return a.state
}

// String, implementation for Stringer interface
func (a *Socket) String() string {
	return fmt.Sprintf("id:%s, remoteaddr:%s", a.ID(), a.conn.RemoteAddr().String())
}

func (a *Socket) writeWork() {
	for p := range a.chSend {
		a.writePacket(p)
	}
}

func (a *Socket) UID() int64 {
	v, has := a.Meta.Load("UID")
	if !has {
		return 0
	}
	return v.(int64)
}

func (a *Socket) SetUID(uid int64) {
	a.Meta.Store("UID", uid)
}

func (a *Socket) SetMeta(k string, v interface{}) {
	a.Meta.Store(k, v)
}

func (a *Socket) GetMeta(k string) (interface{}, bool) {
	return a.Meta.Load(k)
}

// func (a *Socket) readWorkd() error {
// 	p := &Packet{}
// 	for {
// 		p.Reset()
// 		if err := a.readPacket(p); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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

func (a *Socket) readPacket(p *Packet) error {
	var err error
	headRaw := make([]byte, p.HeadLen())

	if a.timeOut > 0 {
		a.conn.SetReadDeadline(time.Now().Add(a.timeOut))
	}

	_, err = io.ReadFull(a.conn, headRaw)
	if err != nil {
		return err
	}

	err = p.PacketHead.Decode(headRaw)
	if err != nil {
		return err
	}

	p.Raw = make([]byte, p.RawLen)

	_, err = io.ReadFull(a.conn, p.Raw)
	return err
}

func (a *Socket) writePacket(p *Packet) error {
	var err error

	head := p.PacketHead.Encode()
	_, err = writeAll(a.conn, head)
	if err != nil {
		return err
	}

	_, err = writeAll(a.conn, p.Raw)
	if err != nil {
		return err
	}
	return err
}

// func (a *tcpSocket) read(p *Packet) error {
// 	// read loop
// 	readBuf := make([]byte, 2048)

// 	for {
// 		n, err := a.conn.Read(readBuf)
// 		if n <= 0 || err != nil {
// 			log.Println(fmt.Sprintf("Conn read error: %s, session will be closed immediately", err.Error()))
// 			return
// 		}

// 		packets, err := a.decoder.Decode(readBuf[:n])
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}

// 		//reflash the conn's active time
// 		atomic.StoreInt64(&a.lastAt, time.Now().Unix())

// 		if a.opt.OnPacket == nil {
// 			continue
// 		}

// 		for _, v := range packets {
// 			if v.Typ == HeartbeatPakcet.Typ {
// 				continue
// 			}
// 			a.opt.OnPacket(a, v)
// 		}
// 	}
// }
