package noop

import "hotwave/servers/battle"

type GamePlayer struct {
	SeatID battle.SeatID
	Score  int64
	Robot  bool
}

func (p *GamePlayer) GetSeatID() battle.SeatID {
	return p.SeatID
}
func (p *GamePlayer) GetScore() int64 {
	return p.Score
}
func (p *GamePlayer) IsRobot() bool {
	return p.Robot
}
