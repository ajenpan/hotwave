package table

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"hotwave/service/battle"
	"hotwave/service/battle/noop"
	"hotwave/service/battle/proto"
)

type logicwarper struct {
	battle.GameLogic
	ontick func(time.Duration)
}

func (l *logicwarper) OnTick(d time.Duration) {
	l.ontick(d)
}

func TestTableTicker(t *testing.T) {
	tk := time.Duration(0)

	logic := noop.NewGameLogic()
	logic = &logicwarper{
		GameLogic: logic,
		ontick: func(d time.Duration) {
			t.Logf("tick %v", d)
			tk += d
		},
	}

	d := NewTable(&proto.BattleConfigure{})
	if err := logic.OnInit(d, nil); err != nil {
		t.Fatal(err)
	}

	err := d.Start(logic)

	if err != nil {
		t.Fatal(err)
	}

	sec := time.Duration(rand.Int31n(10) + 10)
	time.Sleep(time.Second * sec)

	if tk < (sec-1)*time.Second || tk > (sec+1)*time.Second {
		t.Fatal("tick error:", tk)
	}
	fmt.Println("tick", tk)
}
