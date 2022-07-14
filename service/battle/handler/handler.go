package handler

import (
	"context"
	"fmt"
	"sync"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/marshal"
	"hotwave/service/battle"
	"hotwave/service/battle/proto"
	"hotwave/service/battle/table"
	"hotwave/transport"
	"hotwave/utils/calltable"
)

type Handler struct {
	battles sync.Map
	users   sync.Map

	LogicCreator *battle.GameLogicCreator

	CT        *calltable.CallTable
	marshaler *marshal.ProtoMarshaler
}

func New() *Handler {
	h := &Handler{
		LogicCreator: &battle.GameLogicCreator{},
		marshaler:    &marshal.ProtoMarshaler{},
	}
	return h
}

func (h *Handler) CreateBattle(ctx context.Context, in *proto.CreateBattleRequest) (*proto.CreateBattleResponse, error) {
	out := &proto.CreateBattleResponse{}
	// first get game logic
	// var creator battle.GameLogicCreator
	logic, err := h.LogicCreator.CreateLogic(in.GameName)
	if err != nil {
		return nil, err
	}
	d := table.NewTable(in.BattleConf)
	if err = logic.OnInit(d, in.GameConf); err != nil {
		return out, err
	}

	d.Start(logic)

	h.battles.Store(d.ID, d)
	out.BattleId = d.ID
	return out, nil
}

func (h *Handler) WatcherJoinBattle(ctx context.Context, in *proto.WatcherJoinBattleRequest) (*proto.WatcherJoinBattleResponse, error) {
	out := &proto.WatcherJoinBattleResponse{}
	d := h.geBattleById(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}

	d.OnWatcherJoin()
	return out, nil
}

func (h *Handler) OnEvent(topc string, msg protobuf.Message) {

}

func (h *Handler) OnUserConnStat(uid int64, ss transport.SessionStat) {

}

func (h *Handler) OnBattleMessage(uid int64, msg *proto.BattleMessageWrap) {
	b := h.geBattleById(msg.BattleId)
	if b == nil {
		return
	}
	b.OnPlayerMessage(msg.Uid, msg.Topic, msg.Data)
}

func (h *Handler) geBattleById(battleId string) *table.Table {
	if raw, ok := h.battles.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}

func (h *Handler) geBattleByUid(uid int64) *table.Table {
	if raw, ok := h.users.Load(uid); ok {
		return raw.(*table.Table)
	}
	return nil
}
