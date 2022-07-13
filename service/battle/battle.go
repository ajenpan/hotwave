package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type GameTable interface {
	SendMessageToPlayer(Player, proto.Message)
	BroadcastMessage(proto.Message)
	PublishEvent(proto.Message)
	ReportGameOver()
}

type GameStatus int16

type GameLogic interface {
	OnInit(desk GameTable, conf interface{}) error
	OnStart(players []Player) error
	OnTick(time.Duration)
	OnReset()
	OnMessage(p Player, topic string, data []byte)
	OnEvent(topic string, event proto.Message)
}
