package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type SeatID int16

var InvalidSeatID = SeatID(-1)

type GameDesk interface {
	SendMessageToPlayer(Player, proto.Message)
	BroadcastMessage(proto.Message)
	EmitEvent(proto.Message)

	ReportGameOver()
}

type GameStatus int16

type GameLogic interface {
	OnInit(desk GameDesk, conf interface{}) error
	OnStart(players []Player) error
	OnTick(time.Duration)
	OnReset()
	OnMessage(p Player, topic string, data []byte)
	OnEvent(topic string, event proto.Message)
}
