package niuniu

import (
	"math/rand"
	"testing"

	"hotwave/service/battle"
	"hotwave/service/battle/noop"
)

func createTestLogic() (battle.GameLogic, error) {
	nnloigc := CreateLogic()
	conf := &Config{}

	table := noop.NewGameTable()

	err := nnloigc.OnInit(table, conf)
	if err != nil {
		return nil, err
	}
	return nnloigc, nil
}

func createTestPlayer() []battle.Player {
	ret := []battle.Player{}

	count := rand.Int31n(3) + 1
	for i := 0; i < int(count); i++ {
		ret = append(ret, &noop.GamePlayer{
			SeatID: int32(i),
			Score:  1000,
		})
	}
	return ret
}

func TestStartGame(t *testing.T) {
	l, err := createTestLogic()
	if err != nil {
		t.Fatal(err)
	}
	players := createTestPlayer()

	err = l.OnStart(players)
	if err != nil {
		t.Fatal(err)
	}
}
