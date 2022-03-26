package table

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"

	"hotwave/servers/battle"
	pb "hotwave/servers/battle/proto"
)

func NewTable(c *pb.BattleConfigure) *Table {
	ret := &Table{
		ID:   uuid.NewString(),
		conf: c,
	}

	ret.action = make(chan func(), 100)

	// runtime.SetFinalizer(ret, func(d *Desk) {

	// })

	return ret
}

type Table struct {
	ID      string
	conf    *pb.BattleConfigure
	players sync.Map

	StartAt time.Time
	OverAt  time.Time

	logic battle.GameLogic

	isPlaying bool

	// watchers    sync.Map
	// evenReport

	rwlock sync.RWMutex
	action chan func()

	ticker *time.Ticker
}

func (d *Table) Start(l battle.GameLogic) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	if d.logic != nil {
		d.logic.OnReset()
	}

	d.logic = l

	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.ticker = time.NewTicker(1 * time.Second)
	go func() {
		latest := time.Now()
		for range d.ticker.C {
			now := time.Now()
			sub := now.Sub(latest)
			latest = now

			d.action <- func() {
				if d.logic != nil {
					d.logic.OnTick(sub)
				}
			}
		}
	}()

	return nil
}

func (d *Table) Close() {
	if d.ticker != nil {
		d.ticker.Stop()
	}

	close(d.action)
}

func (d *Table) SendMessage(p battle.Player, msg proto.Message) {

}

func (d *Table) OnWatcherJoin() {
	d.action <- func() {

	}
}

func (d *Table) BroadcastMessage(msg proto.Message) {

	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*player); ok && p != nil {
			d.SendMessage(p, msg)
		}
		return true
	})
}

func (d *Table) EmitEvent(event proto.Message) {

}

func (d *Table) ReportGameStart() {
	d.isPlaying = true
	d.StartAt = time.Now()
}

func (d *Table) ReportGameOver() {
	d.isPlaying = false
	d.StartAt = time.Now()
}

func (d *Table) JoinPlayer() {}

func (d *Table) getPlayer(uid int64) *player {
	if p, has := d.players.Load(uid); has {
		return p.(*player)
	}
	return nil
}

func (d *Table) OnBattleMessage(ctx context.Context, msg *pb.BattleMessageWrap) {
	d.action <- func() {
		//TODO: in a channle
		p := d.getPlayer(msg.Uid)
		if p != nil && d.logic != nil {
			d.logic.OnMessage(p, msg.Topic, msg.Data)
		}
	}
}
