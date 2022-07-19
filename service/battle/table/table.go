package table

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"hotwave/event"
	evproto "hotwave/event/proto"
	log "hotwave/logger"
	"hotwave/service/battle"
	pb "hotwave/service/battle/proto"
)

type TableOption struct {
	ID        string
	Publisher event.Publisher
	Conf      *pb.BattleConfigure
}

func NewTable(opt TableOption) *Table {
	if opt.ID == "" {
		opt.ID = uuid.NewString()
	}

	ret := &Table{
		TableOption: opt,
		CreateAt:    time.Now(),
	}

	ret.action = make(chan func(), 100)

	return ret
}

type Table struct {
	TableOption

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time

	logic battle.GameLogic

	IsPlaying bool

	// watchers    sync.Map
	// evenReport

	rwlock  sync.RWMutex
	players sync.Map

	action chan func()

	ticker *time.Ticker
}

func (d *Table) Init(logic battle.GameLogic, players []*Player, logicConf interface{}) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	if d.logic != nil {
		d.logic.OnReset()
	}

	battlePlayers := make([]battle.Player, len(players))
	for i, p := range players {
		// store player
		d.players.Store(p.Uid, p)

		battlePlayers[i] = p
	}

	err := logic.OnInit(d, battlePlayers, logicConf)
	if err != nil {
		return err
	}

	d.logic = logic

	return nil
}

func (d *Table) Start() error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	go func(jobque chan func()) {
		for job := range jobque {
			job()
		}

	}(d.action)

	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.ticker = time.NewTicker(1 * time.Second)
	go func(ticker *time.Ticker) {
		latest := time.Now()
		for now := range ticker.C {
			// now := time.Now()
			sub := now.Sub(latest)
			latest = now

			d.action <- func() {
				if d.logic != nil {
					d.logic.OnTick(sub)
				}
			}
		}
	}(d.ticker)

	d.logic.OnStart()

	return nil
}

func (d *Table) Close() {
	if d.ticker != nil {
		d.ticker.Stop()
	}
	close(d.action)
}

func (d *Table) SendMessageToPlayer(p battle.Player, msg proto.Message) {
	rp := p.(*Player)
	rp.SendMessage(msg)

	log.Infof("send message to player: %v, %s: %v", rp.Uid, string(proto.MessageName(msg)), msg)
}

func (d *Table) OnWatcherJoin() {
	d.action <- func() {
		//todo:
	}
}

func (d *Table) BroadcastMessage(msg proto.Message) {
	log.Infof("BroadcastMessage: %s: %v", string(proto.MessageName(msg)), msg)

	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Player); ok && p != nil {
			p.SendMessage(msg)
		}
		return true
	})
}

func (d *Table) PublishEvent(event proto.Message) {
	log.Infof("PublishEvent: %s: %v", string(proto.MessageName(event)), event)

	if d.Publisher == nil {
		return
	}
	//TODO:
	warp := &evproto.EventMessage{
		Topic:     string(proto.MessageName(event)),
		Timestamp: time.Now().Unix(),
	}
	d.Publisher.Publish(warp)
}

func (d *Table) ReportGameStart() {
	if d.IsPlaying {
		log.Error("table is playing")
		return
	}
	d.IsPlaying = true
	d.StartAt = time.Now()

	d.PublishEvent(&pb.BattleStartEvent{})
}

func (d *Table) ReportGameOver() {
	if !d.IsPlaying {
		log.Error("table is not playing")
		return
	}

	d.IsPlaying = false
	d.OverAt = time.Now()

	d.PublishEvent(&pb.BattleOverEvent{})
}

func (d *Table) GetPlayer(uid int64) *Player {
	if p, has := d.players.Load(uid); has {
		return p.(*Player)
	}
	return nil
}

func (d *Table) OnPlayerMessage(uid int64, topic string, iraw []byte) {
	// here is not safe
	// msg := proto.Clone(fmsg).(*pb.BattleMessageWrap)

	d.action <- func() {
		p := d.GetPlayer(uid)
		if p != nil && d.logic != nil {
			d.logic.OnMessage(p, topic, iraw)
		}
	}
}
