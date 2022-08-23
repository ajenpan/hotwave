package gate

import "hotwave/transport"

func NewTcpGater() *TcpGater {
	return nil
}

type TcpGater struct {
}

func (g *TcpGater) Start() error {

	return nil
}

func (g *TcpGater) Stop() {

}

func (g *TcpGater) OnStat(socket transport.Session, status transport.SessionStat) {

}

func (g *TcpGater) OnGateMessage(session transport.Session, raw []byte) {

}
