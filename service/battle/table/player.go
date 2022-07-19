package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	pb "hotwave/service/battle/proto"
)

func NewPlayer(p *pb.PlayerInfo) *Player {
	return &Player{
		PlayerInfo: protobuf.Clone(p).(*pb.PlayerInfo),
	}
}

func NewPlayers(infos []*pb.PlayerInfo) ([]*Player, error) {
	ret := make([]*Player, len(infos))
	for i, info := range infos {
		ret[i] = NewPlayer(info)
	}

	// check seatid
	for _, v := range ret {
		if v.SeatId == 0 {
			return nil, fmt.Errorf("seat id is 0")
		}
		if v.Uid == 0 {
			return nil, fmt.Errorf("uid is 0")
		}
	}

	return ret, nil
}

type Player struct {
	*pb.PlayerInfo
	// table *Table
}

func (p *Player) GetScore() int64 {
	return p.PlayerInfo.Score
}

func (p *Player) GetUserID() int64 {
	return p.PlayerInfo.Uid
}

func (p *Player) GetSeatID() int32 {
	return p.PlayerInfo.SeatId
}

func (p *Player) IsRobot() bool {
	return p.PlayerInfo.IsRobot
}

func (p *Player) SendMessage(protobuf.Message) error {

	return nil
}
