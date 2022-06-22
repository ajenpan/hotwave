package table

import (
	protobuf "google.golang.org/protobuf/proto"

	pb "hotwave/services/battle/proto"
)

type player struct {
	*pb.PlayerInfo
	// table *Table
}

// func NewPlayer(u user.User) battle.Player {
// 	return &player{
// 		User: u,
// 	}
// }

func (p *player) GetScore() int64 {
	return int64(p.Score)
}

func (p *player) GetUserID() int64 {
	return p.Uid
}

func (p *player) GetSeatID() int32 {
	return p.SeatId
}

func (p *player) IsRobot() bool {
	return false
}

func (p *player) SendMessage(protobuf.Message) error {
	return nil
}
