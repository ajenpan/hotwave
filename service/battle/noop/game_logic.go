package noop

import (
	"time"

	"google.golang.org/protobuf/proto"

	"hotwave/service/battle"
)

func NewGameLogic() battle.GameLogic {
	return &GameLogic{}
}

type GameLogic struct {
}

func (gl *GameLogic) OnInit(battle.GameTable, []battle.Player, interface{}) error {
	return nil
}
func (gl *GameLogic) OnStart() error {
	return nil
}
func (gl *GameLogic) OnTick(time.Duration) {

}
func (gl *GameLogic) OnReset() {

}
func (gl *GameLogic) OnMessage(battle.Player, string, []byte) {

}
func (gl *GameLogic) OnEvent(string, proto.Message) {

}
