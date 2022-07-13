package table

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	log "hotwave/logger"
	"hotwave/service/battle"
	pb "hotwave/service/battle/proto"
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

func (d *Table) SendMessageToPlayer(p battle.Player, msg proto.Message) {
	err := p.SendMessage(msg)
	if err != nil {
		log.Error(err)
	}
}

func (d *Table) OnWatcherJoin() {
	d.action <- func() {
		//todo:
	}
}

func (d *Table) BroadcastMessage(msg proto.Message) {
	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*player); ok && p != nil {
			d.SendMessageToPlayer(p, msg)
		}
		return true
	})
}

func (d *Table) PublishEvent(event proto.Message) {
	//TODO:
}

func (d *Table) ReportGameStart() {
	d.isPlaying = true
	d.StartAt = time.Now()
}

func (d *Table) ReportGameOver() {
	d.isPlaying = false
	d.StartAt = time.Now()
}

func (d *Table) PlayerJoin() {

}

func (d *Table) getPlayer(uid int64) *player {
	if p, has := d.players.Load(uid); has {
		return p.(*player)
	}
	return nil
}

func (d *Table) OnBattleMessage(ctx context.Context, fmsg *pb.BattleMessageWrap) {
	// here is not safe
	msg := proto.Clone(fmsg).(*pb.BattleMessageWrap)

	d.action <- func() {
		p := d.getPlayer(msg.Uid)
		if p != nil && d.logic != nil {
			d.logic.OnMessage(p, msg.Topic, msg.Data)
		}
	}
}
