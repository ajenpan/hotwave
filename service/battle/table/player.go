package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/service/battle"
	bf "hotwave/service/battle"
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
	tableid string

	sender func(msgname string, raw []byte) error
}

func (p *Player) Score() int64 {
	return p.PlayerInfo.Score
}

func (p *Player) UserID() int64 {
	return p.PlayerInfo.Uid
}

func (p *Player) SeatID() battle.SeatID {
	return battle.SeatID(p.PlayerInfo.SeatId)
}

func (p *Player) Role() bf.RoleType {
	if p.PlayerInfo.IsRobot {
		return bf.RoleType_Robot
	}
	return bf.RoleType_Player
}

func (p *Player) SendMessage(protobuf.Message) error {

	return nil
}

func (p *Player) Send(msgname string, raw []byte) error {
	return p.sender(msgname, raw)
}
