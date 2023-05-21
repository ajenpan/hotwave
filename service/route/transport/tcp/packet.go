package tcp

import (
	"encoding/binary"
	"errors"
)

// Codec constants.
const (
	//16MB
	MaxPacketSize = 1<<(3*8) - 1
)

// TODO: rand in packet

var ErrWrongPacketHeadLen = errors.New("wrong packet head len")
var ErrWrongPacketType = errors.New("wrong packet type")
var ErrPacketSizeExcced = errors.New("packet size exceed")
var ErrParseHead = errors.New("parse head error")
var ErrDisconn = errors.New("socket disconnected")

// -<PacketType>-|-<HeadLen>-|-<BodyLen>-|-<Body>-
// -1------------|-1---------|-4---------|--------

type PacketType = uint8

const (
	// inner
	PacketTypInnerStartAt_ PacketType = iota
	PacketTypHeartbeat
	PacketTypAck
	PacketTypInnerEndAt_
)
const (
	// user
	PacketTypStartAt_ PacketType = iota + 10
	PacketTypRout
	PacketTypRoutDeliver
	PacketTypRoutErr
	PacketTypEndAt_
)

type PacketHead [1 + 1 + 4]byte

func (hr *PacketHead) GetType() uint8 {
	return hr[0]
}
func (hr *PacketHead) GetHeadLen() uint8 {
	return hr[1]
}
func (hr *PacketHead) GetBodyLen() uint32 {
	return binary.LittleEndian.Uint32(hr[2:6])
}

func (hr *PacketHead) SetType(t uint8) {
	hr[0] = t
}
func (hr *PacketHead) SetHeadLen(t uint8) {
	hr[1] = t
}
func (hr *PacketHead) SetBodyLen(l uint32) {
	binary.LittleEndian.PutUint32(hr[2:6], l)
}

func (hr *PacketHead) Reset() {
	for i := 0; i < len(hr); i++ {
		hr[i] = 0
	}
}

type PackFrame struct {
	PacketHead
	Head []byte
	Body []byte
}

func (p *PackFrame) Reset() {
	p.PacketHead.Reset()
	p.Head = p.Head[:0]
	p.Body = p.Body[:0]
}

func (p *PackFrame) Clone() *PackFrame {
	return &PackFrame{
		PacketHead: p.PacketHead,
		Head:       p.Head[:],
		Body:       p.Body[:],
	}
}

// |-askid-|
// |-4-----|
// type RequestHead []byte
// type ResponseHead []byte
// type AsyncHead []byte

type RoutDeliverHead []byte
type RoutMsgTyp uint8

const (
	RoutTypAsync = iota
	RoutTypRequest
	RoutTypResponse
)

func NewRoutDeliverHead() RoutDeliverHead {
	return make([]byte, 25)
}

func (h RoutDeliverHead) GetTargetUID() uint64 {
	return binary.LittleEndian.Uint64(h[0:8])
}

func (h RoutDeliverHead) GetSrouceUID() uint64 {
	return binary.LittleEndian.Uint64(h[8:16])
}

func (h RoutDeliverHead) GetAskID() uint32 {
	return binary.LittleEndian.Uint32(h[16:20])
}

func (h RoutDeliverHead) GetMsgID() uint32 {
	return binary.LittleEndian.Uint32(h[20:24])
}

func (h RoutDeliverHead) GetMsgTyp() uint8 {
	return h[24]
}

func (h RoutDeliverHead) SetTargetUID(u uint64) {
	binary.LittleEndian.PutUint64(h[0:8], u)
}

func (h RoutDeliverHead) SetSrouceUID(u uint64) {
	binary.LittleEndian.PutUint64(h[8:16], u)
}

func (h RoutDeliverHead) SetAskID(id uint32) {
	binary.LittleEndian.PutUint32(h[16:20], id)
}

func (h RoutDeliverHead) SetMsgID(id uint32) {
	binary.LittleEndian.PutUint32(h[20:24], id)
}

func (h RoutDeliverHead) SetMsgTyp(typ uint8) {
	h[24] = typ
}

type RoutErrHead RoutDeliverHead
