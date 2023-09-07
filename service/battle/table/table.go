package table

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"hotwave/event"
	log "hotwave/logger"

	bf "hotwave/service/battle"
	pb "hotwave/service/battle/proto"
)

type TableOption struct {
	ID             string
	EventPublisher event.Publisher
	Conf           *pb.BattleConfigure
}

func NewTable(opt TableOption) *Table {
	if opt.ID == "" {
		opt.ID = uuid.NewString()
	}

	ret := &Table{
		TableOption: &opt,
		CreateAt:    time.Now(),
	}

	ret.action = make(chan func(), 100)

	return ret
}

type Table struct {
	*TableOption

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time

	battle bf.Logic

	// watchers    sync.Map
	// evenReport

	rwlock  sync.RWMutex
	players sync.Map

	action chan func()

	ticker *time.Ticker

	battleStatus bf.GameStatus
}

func (d *Table) Init(logic bf.Logic, players []*Player, logicConf interface{}) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	if d.battle != nil {
		d.battle.OnReset()
	}

	battlePlayers := make([]bf.Player, len(players))
	for i, p := range players {
		// store player
		d.players.Store(p.Uid, p)

		battlePlayers[i] = p
	}

	if err := logic.OnInit(d, logicConf); err != nil {
		return err
	}

	if err := logic.OnPlayerJoin(battlePlayers); err != nil {
		return err
	}

	d.battle = logic

	switch conf := d.Conf.StartCondition.(type) {
	case *pb.BattleConfigure_Delayed:
		if conf.Delayed > 0 {
			log.Infof("start table after %d seconds", conf.Delayed)
			time.AfterFunc(time.Duration(conf.Delayed)*time.Second, func() {
				err := d.Start()
				if err != nil {
					log.Error(err)
				}
			})
		}
	}
	return nil
}

func (d *Table) pushAction(f func()) {
	d.action <- f
}

func (d *Table) AfterFunc(td time.Duration, f func()) {
	time.AfterFunc(td, func() {
		d.pushAction(f)
	})
}

func (d *Table) Start() error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	go func() {
		safecall := func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic: %v", err)
				}
			}()
			f()
		}

		for job := range d.action {
			safecall(job)
		}
	}()

	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.ticker = time.NewTicker(1 * time.Second)
	go func(ticker *time.Ticker) {
		latest := time.Now()
		for now := range ticker.C {
			sub := now.Sub(latest)
			latest = now

			d.pushAction(func() {
				if d.battle != nil {
					d.battle.OnTick(sub)
				}
			})
		}
	}(d.ticker)

	return d.battle.OnStart()
}

func (d *Table) Close() {
	if d.ticker != nil {
		d.ticker.Stop()
	}
	close(d.action)
}

func (d *Table) ReportBattleStatus(s bf.GameStatus) {
	if d.battleStatus == s {
		return
	}

	statusBefore := d.battleStatus
	d.battleStatus = s

	event := &pb.BattleStatusChangeEvent{
		StatusBefore: int32(statusBefore),
		StatusNow:    int32(s),
		BattleId:     d.ID,
	}

	d.PublishEvent(event)

	switch s {
	case bf.BattleStatus_Idle:
	case bf.BattleStatus_Start:
		d.reportGameStart()
	case bf.BattleStatus_Over:
		d.reportGameOver()
	}
}

func (d *Table) ReportBattleEvent(topic string, event proto.Message) {
	d.PublishEvent(event)
}

func (d *Table) SendMessageToPlayer(p bf.Player, msg proto.Message) {

	rp, ok := p.(*Player)
	if !ok {
		log.Error("player is not Player")
		return
	}

	err := rp.SendMessage(msg)
	if err != nil {
		log.Errorf("send message to player: %v, %s: %v", rp.Uid, string(proto.MessageName(msg)), msg)
	} else {
		log.Debugf("send message to player: %v, %s: %v", rp.Uid, string(proto.MessageName(msg)), msg)
	}
}

func (d *Table) BroadcastMessage(msg proto.Message) {
	msgname := string(proto.MessageName(msg))
	log.Debugf("BroadcastMessage: %s: %v", msgname, msg)

	raw, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}

	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Player); ok && p != nil {
			err := p.Send(msgname, raw)
			if err != nil {
				log.Error(err)
			}
		}
		return true
	})
}

func (d *Table) IsPlaying() bool {
	return d.battleStatus == bf.BattleStatus_Start
}

func (d *Table) reportGameStart() {
	d.StartAt = time.Now()
}

func (d *Table) reportGameOver() {
	d.OverAt = time.Now()
}

func (d *Table) GetPlayer(uid int64) *Player {
	if p, has := d.players.Load(uid); has {
		return p.(*Player)
	}
	return nil
}

func (d *Table) PublishEvent(eventmsg proto.Message) {
	if d.EventPublisher == nil {
		return
	}

	log.Infof("PublishEvent: %s: %v", string(proto.MessageName(eventmsg)), eventmsg)

	raw, err := proto.Marshal(eventmsg)
	if err != nil {
		log.Error(err)
		return
	}
	warp := &event.Event{
		Topic:     string(proto.MessageName(eventmsg)),
		Timestamp: time.Now().Unix(),
		Data:      raw,
	}
	d.EventPublisher.Publish(warp)
}

func (d *Table) OnPlayerMessage(uid int64, msgid int, iraw []byte) {
	d.action <- func() {
		p := d.GetPlayer(uid)
		if p != nil && d.battle != nil {
			d.battle.OnPlayerMessage(p, msgid, iraw)
		}
	}
}
