package tcp

import (
	"errors"
	"fmt"
)

// Codec constants.
const (
	//16MB
	MaxPacketSize = 1<<(3*8) - 1
	HeadLength    = 4
)

// TODO: rand in packet

var ErrWrongPacketHeadLen = errors.New("wrong packet head len")
var ErrWrongPacketType = errors.New("wrong packet type")
var ErrPacketSizeExcced = errors.New("packet size exceed")
var ErrParseHead = errors.New("parse head error")

// Encode create a packet.packet from  the raw bytes slice and then encode to network bytes slice
// -<Type>-|-<length>-|-<data>-
// -1------|-3--------|--------
// 1 byte packet type, 3 bytes packet data length(big end), and data segment

type PacketType = uint8

const (
	PacketTypeStatAt_   PacketType = iota
	PacketTypeHeartbeat PacketType = iota
	PacketTypeAck       PacketType = iota
	PacketTypePacket    PacketType = iota
	PacketTypeEcho      PacketType = iota
	PacketTypeError     PacketType = iota
	PacketTypeEndAt_    PacketType = iota
)

type PacketHead struct {
	Typ    PacketType
	RawLen int32
}

func (p *PacketHead) HeadLen() int {
	return HeadLength
}

func (p *PacketHead) Encode() []byte {
	head := make([]byte, p.HeadLen())
	head[0] = p.Typ
	copy(head[1:HeadLength], intToBytes(p.RawLen))
	return head
}

func (p *PacketHead) Decode(headRaw []byte) error {
	if len(headRaw) != p.HeadLen() {
		return fmt.Errorf("head len is wrong")
	}
	p.Typ = uint8(headRaw[0])
	if p.Typ <= PacketTypeStatAt_ || p.Typ > PacketTypeEndAt_ {
		return fmt.Errorf("packet type is invalid")
	}

	p.RawLen = bytesToInt(headRaw[1:])
	if p.RawLen == -1 || p.RawLen > MaxPacketSize {
		return fmt.Errorf("pakcet head decode error")
	}
	return nil
}

// Packet represents a network Packet.
type Packet struct {
	PacketHead
	Raw []byte
}

func NewAckPacket() *Packet {
	return &Packet{
		PacketHead: PacketHead{
			Typ: PacketTypeAck,
		},
	}
}

//String represents the packet's in text mode.
func (p *Packet) String() string {
	return fmt.Sprintf("type:%d, len:%d, raw:%X", p.Typ, len(p.Raw), string(p.Raw))
}

func (p *Packet) RawData() []byte {
	return p.Raw
}

func (p *Packet) PacketType() PacketType {
	return p.Typ
}

func (p *Packet) Reset() {
	p.Typ = PacketTypeStatAt_
	p.RawLen = 0
	p.Raw = nil
}

var HeartbeatPakcet = &Packet{
	PacketHead: PacketHead{
		Typ:    PacketTypeHeartbeat,
		RawLen: 0,
	},
}

// Decode packet data length byte to int(Big end)
func bytesToInt(b []byte) int32 {
	if len(b) != 3 {
		return -1
	}

	result := int32(0)
	for _, v := range b {
		result = (result << 8) + int32(v)
	}
	return result
}

// Encode packet data length to bytes(Big end)
func intToBytes(n int32) []byte {
	buf := make([]byte, 3)
	buf[0] = byte((n >> 16) & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte(n & 0xFF)
	return buf
}

// // A Codec reads and decodes network data slice
// type decode struct {
// 	buf *bytes.Buffer

// 	size int  // last packet length
// 	typ  byte // last packet type
// }

// func newDecoder() *decode {
// 	return &decode{
// 		buf:  bytes.NewBuffer(nil),
// 		size: -1,
// 	}
// }

// func (c *decode) forward() error {
// 	header := c.buf.Next(HeadLength)
// 	c.typ = header[0]

// 	c.size = bytesToInt(header[1:])

// 	// packet length limitation
// 	if c.size > MaxPacketSize {
// 		return ErrPacketSizeExcced
// 	}
// 	return nil
// }

// // Decode decode the network bytes slice to packet.packet(s)
// // TODO(Warning): shared slice
// func (c *decode) Decode(data []byte) ([]*Packet, error) {
// 	c.buf.Write(data)

// 	var (
// 		packets []*Packet
// 		err     error
// 	)
// 	// check length
// 	if c.buf.Len() < HeadLength {
// 		return nil, err
// 	}

// 	// first time
// 	if c.size < 0 {
// 		if err = c.forward(); err != nil {
// 			return nil, err
// 		}
// 	}

// 	for c.size <= c.buf.Len() {
// 		p := &Packet{Typ: (c.typ), Raw: c.buf.Next(c.size)}
// 		packets = append(packets, p)

// 		// more packet
// 		if c.buf.Len() < HeadLength {
// 			c.size = -1
// 			break
// 		}

// 		if err = c.forward(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return packets, nil
// }
