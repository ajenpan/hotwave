package gate

import "hotwave/transport/tcpsvr"

func NewTCPGate(addr string) *TCPGate {
	ret := &TCPGate{}
	ret.svr = tcpsvr.NewServer(&tcpsvr.ServerOptions{
		Address:   addr,
		OnMessage: ret.OnMessage,
		OnConn:    ret.OnConnStat,
	})
	return ret
}

type TCPGate struct {
	svr *tcpsvr.Server
}

type TCPSession struct {
}

func (g *TCPGate) Start() {
	// svr.Start()

}

func (g *TCPGate) OnMessage(s *tcpsvr.Socket, p *tcpsvr.Packet) {
	if p.Typ != tcpsvr.PacketTypePacket {
		return
	}

}

func (g *TCPGate) OnConnStat(s *tcpsvr.Socket, ss SocketStat) {
	switch ss {
	case tcpsvr.SocketStatConnected:

	case tcpsvr.SocketStatDisconnected:

	}
}
