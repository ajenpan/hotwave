package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type SeatID = int32

type Player interface {
	GetSeatID() int32
	GetScore() int64 //game jetton
	IsRobot() bool
}

type GameTable interface {
	SendMessageToPlayer(Player, proto.Message)
	BroadcastMessage(proto.Message)
	PublishEvent(proto.Message)

	ReportGameStart()
	ReportGameOver()
}

type GameStatus int16

type GameLogic interface {
	OnInit(desk GameTable, conf interface{}) error
	OnPlayerJoin([]Player) error
	OnStart() error
	OnTick(time.Duration)
	OnReset()
	OnMessage(p Player, topic string, data []byte)
	OnEvent(topic string, event proto.Message)
}
