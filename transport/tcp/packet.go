package tcp

import (
	"errors"
	"io"
)

// Codec constants.
const (
	//8MB 8388607
	MaxPacketSize = 0x7FFFFF
)

var ErrWrongPacketType = errors.New("wrong packet type")
var ErrBodySizeWrong = errors.New("packet body size error")
var ErrHeadSizeWrong = errors.New("packet head size error")
var ErrParseHead = errors.New("parse packet error")
var ErrDisconn = errors.New("socket disconnected")

// -<PacketType>-|-<BodyLen>-|-<Body>-
// -1------------|-3---------|--------

type PacketType = uint8

const (
	// inner
	PacketTypInnerStartAt_ PacketType = iota
	PacketTypAck
	PacketTypHeartbeat
	PacketTypeEcho
	PacketTypInnerEndAt_
)

type Packet interface {
	io.ReaderFrom
	io.WriterTo
}

func Uint24(b []uint8) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func PutUint24(b []uint8, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = uint8(v)
	b[1] = uint8(v >> 8)
	b[2] = uint8(v >> 16)
}

const PackMetaLen = 4

type THVPacketHead []uint8

func NewPackMeta() THVPacketHead {
	return make([]uint8, PackMetaLen)
}

func (hr THVPacketHead) GetType() uint8 {
	return hr[0]
}

func (hr THVPacketHead) GetBodyLen() uint32 {
	return Uint24(hr[1:4])
}

func (hr THVPacketHead) SetType(t uint8) {
	hr[0] = t
}

func (hr THVPacketHead) SetBodyLen(l uint32) {
	PutUint24(hr[1:4], l)
}

func (hr THVPacketHead) Reset() {
	for i := 0; i < len(hr); i++ {
		hr[i] = 0
	}
}

func NewEmptyTHVPacket() *THVPacket {
	return &THVPacket{
		Head: NewPackMeta(),
	}
}

func NewPackFrame(t uint8, b []uint8) *THVPacket {
	p := NewEmptyTHVPacket()
	p.SetType(t)
	p.SetBody(b)
	return p
}

type THVPacket struct {
	Head THVPacketHead
	Body []uint8
}

func (p *THVPacket) ReadFrom(reader io.Reader) (int64, error) {
	var err error
	metalen, err := io.ReadFull(reader, p.Head)
	if err != nil {
		return 0, err
	}

	bodylen := p.Head.GetBodyLen()
	if bodylen > 0 {
		p.Body = make([]byte, bodylen)
		_, err = io.ReadFull(reader, p.Body)
		if err != nil {
			return 0, err
		}
	}
	return int64(metalen + int(bodylen)), nil
}

func (p *THVPacket) WriteTo(writer io.Writer) (int64, error) {
	ret := int64(0)

	n, err := writer.Write(p.Head)
	ret += int64(n)
	if err != nil {
		return ret, err
	}

	if len(p.Body) > 0 {
		n, err = writer.Write(p.Body)
		ret += int64(n)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func (p *THVPacket) Name() string {
	return "tcp-binary"
}

func (p *THVPacket) Reset() {
	p.Head.Reset()
	p.Body = p.Body[:0]
}

func (p *THVPacket) Clone() *THVPacket {
	return &THVPacket{
		Head: p.Head[:],
		Body: p.Body[:],
	}
}

func (p *THVPacket) SetType(t uint8) {
	p.Head.SetType(t)
}

func (p *THVPacket) GetType() uint8 {
	return p.Head.GetType()
}

func (p *THVPacket) SetBody(b []uint8) {
	p.Body = b
	p.Head.SetBodyLen(uint32(len(b)))
}

func (p *THVPacket) GetBody() []uint8 {
	return p.Body
}
