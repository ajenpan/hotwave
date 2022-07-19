package noop

import (
	"google.golang.org/protobuf/proto"

	"hotwave/service/battle"
)

func NewGameTable() *GameDesk {
	return &GameDesk{}
}

type GameDesk struct {
}

func (gd *GameDesk) SendMessageToPlayer(battle.Player, proto.Message) {

}

func (gd *GameDesk) BroadcastMessage(proto.Message) {

}

func (gd *GameDesk) PublishEvent(proto.Message) {

}

func (gd *GameDesk) ReportGameOver() {
}
