package tcpsvr

import (
	"encoding/binary"

	"hotwave/transport/tcp"
	"hotwave/utils/calltable"
	"hotwave/utils/marshal"
)

type TcpSvr struct {
	CT      *calltable.CallTable[uint32]
	Marshal marshal.Marshaler
}

func (s *TcpSvr) Start() error {

	return nil
}

func (s *TcpSvr) Stop() error {

	return nil
}

func (svr *TcpSvr) OnMessage(session *tcp.Socket, packet *tcp.THVPacket) {
	body := packet.GetBody()
	ctype := packet.GetType()

	if len(body) < 4 {
		return
	}

	msgid := binary.LittleEndian.Uint32(body)
	method := svr.CT.Get(msgid)
	if method == nil {
		return
	}

	req := method.NewRequest()
	svr.Marshal.Unmarshal(body[4:], req)

	if ctype == 4 {
		res := method.Call(req)
		packet.Reset()
		respi := res[0]

		respraw, err := svr.Marshal.Marshal(respi)
		if err != nil {
			return
		}

		respHead := make([]byte, 4)
		binary.LittleEndian.PutUint32(respHead, msgid)
		packet.Body = append(respHead, respraw...)
		session.SendPacket(packet)
	}
}

func (svr *TcpSvr) OnConn(session *tcp.Socket, ss tcp.SocketStat) {

}
