package niuniu

import (
	"math/rand"
	"testing"
	"time"

	"hotwave/service/battle/table"
)

func createTestPlayer(c int) []*table.Player {
	ret := []*table.Player{}

	for i := 0; i < c; i++ {
		ret = append(ret, &table.Player{})
	}
	return ret
}

func TestStartGame(t *testing.T) {
	ta := table.NewTable(table.TableOption{})
	var err error

	nnloigc := CreateNiuniu()
	count := rand.Int31n(3) + 1
	players := createTestPlayer(int(count))
	err = ta.Init(nnloigc, players, &Config{
		3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = ta.Start()
	if err != nil {
		t.Fatal(err)
	}

	if ta.IsPlaying() != true {
		t.Fatal("game is not playing")
	}

	time.Sleep(30 * time.Second)

	if ta.IsPlaying() == true {
		t.Fatal("game is not playing")
	}
}
