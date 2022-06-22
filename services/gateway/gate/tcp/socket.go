package tcp

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/services/gateway/proto"
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
		chSend:  make(chan *Packet, 10),
		state:   SocketStatConnected,
	}
	return ret
}

type Socket struct {
	sync.RWMutex // export

	conn   net.Conn   // low-level conn fd
	state  SocketStat // current state
	id     string
	chSend chan *Packet // push message queue

	timeOut      time.Duration
	Meta         sync.Map
	serialNumber uint64
}

func (s *Socket) ID() string {
	return s.id
}

func ConverPacket(msg *proto.ClientMessageWraper) *Packet {
	raw, err := protobuf.Marshal(msg)
	if err != nil {
		return nil
	}

	packet := &Packet{
		Raw: raw,
		PacketHead: PacketHead{
			Typ:    PacketTypePacket,
			RawLen: int32(len(raw)),
		},
	}
	return packet
}

func ConverMessage(p *Packet) (*proto.ClientMessageWraper, error) {
	msg := &proto.ClientMessageWraper{}
	err := protobuf.Unmarshal(p.Raw, msg)
	return msg, err
}

func (a *Socket) SendWrap(msg *proto.ClientMessageWraper) error {
	return a.sendPacket(ConverPacket(msg))
}

func (a *Socket) sendPacket(p *Packet) error {
	if atomic.LoadInt32(&a.state) == SocketStatDisconnected {
		return fmt.Errorf("sendPacket failed, the socket is disconnected")
	}
	a.chSend <- p
	return nil
}

func (a *Socket) Send(msg protobuf.Message) error {
	raw, err := protobuf.Marshal(msg)
	if err != nil {
		return err
	}

	methodName := string(protobuf.MessageName(msg))
	wrap := &proto.ClientMessageWraper{
		Body: raw,
	}

	if strings.HasSuffix(methodName, "Response") {
		wrap.Method = methodName[:len(methodName)-len("Response")]
		wrap.Typ = proto.ClientMessageWraper_Response
	} else {
		wrap.Method = methodName
		wrap.Typ = proto.ClientMessageWraper_Async
	}

	return a.SendWrap(wrap)
}

func (a *Socket) Recv() (*Packet, error) {
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
		a.conn = nil
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

func (a *Socket) SerialNumber() uint64 {
	return atomic.AddUint64(&a.serialNumber, 1)
}

func (a *Socket) UID() uint64 {
	v, has := a.Meta.Load("UID")
	if !has {
		return 0
	}
	return v.(uint64)
}

func (a *Socket) SetUID(uid uint64) {
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
	if atomic.LoadInt32(&a.state) == SocketStatDisconnected {
		return fmt.Errorf("recv packet failed, the socket is disconnected")
	}

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
	if atomic.LoadInt32(&a.state) == SocketStatDisconnected {
		return fmt.Errorf("writePacket failed, the socket is disconnected")
	}
	var err error
	p.RawLen = int32(len(p.Raw))

	head := make([]byte, p.HeadLen())

	p.PacketHead.Encode(head)
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
