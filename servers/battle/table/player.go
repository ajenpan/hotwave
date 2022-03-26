package table

import (
	protobuf "google.golang.org/protobuf/proto"

	"hotwave/servers/battle"
	pb "hotwave/servers/battle/proto"
	"hotwave/user"
)

type player struct {
	*pb.PlayerInfo
	user.User
	desk *Table
}

func NewPlayer(u user.User) battle.Player {
	return &player{
		User: u,
	}
}

func (p *player) GetScore() int64 {
	return int64(p.Score)
}

func (p *player) GetUserID() int64 {
	return p.Uid
}

func (p *player) GetSeatID() battle.SeatID {
	return battle.SeatID(p.SeatId)
}

func (p *player) IsRobot() bool {
	return false
}

func (p *player) SendMessage(protobuf.Message) error {
	return nil
}
