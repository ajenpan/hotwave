package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type SeatID int32
type GameStatus int32
type RoleType int32

const (
	BattleStatus_Idle GameStatus = iota
	BattleStatus_Start
	BattleStatus_Over
)

const (
	RoleType_Player RoleType = iota
	RoleType_Robot
)

type Player interface {
	SeatID() SeatID
	Score() int64 //game jetton
	Role() RoleType
}

type Table interface {
	SendMessageToPlayer(Player, proto.Message)
	BroadcastMessage(proto.Message)

	ReportBattleStatus(GameStatus)
	ReportBattleEvent(topic string, event proto.Message)
}

type Logic interface {
	OnInit(t Table, conf interface{}) error
	OnPlayerJoin([]Player) error
	OnStart() error
	OnTick(time.Duration)
	OnReset()
	OnPlayerMessage(p Player, msgid uint32, data []byte)
	OnCommand(topic string, data []byte)
}
