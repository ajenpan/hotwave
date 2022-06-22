package handler

import (
	"context"
	"fmt"
	"sync"

	protobuf "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	// "github.com/google/uuid"
	logger "hotwave/logger"
	"hotwave/marshal"
	"hotwave/services/battle"
	"hotwave/services/battle/proto"
	"hotwave/services/battle/table"
	gateproto "hotwave/services/gateway/proto"
	utilHandler "hotwave/transport"
)

type Handler struct {
	desks        sync.Map
	LogicCreator *battle.GameLogicCreator

	ct        *utilHandler.CallTable
	marshaler *marshal.ProtoMarshaler
}

func New() *Handler {
	h := &Handler{
		LogicCreator: &battle.GameLogicCreator{},
		marshaler:    &marshal.ProtoMarshaler{},
	}
	h.ct = utilHandler.ExtractAsyncMethod("battle", proto.File_servers_battle_proto_battle_proto.Messages(), h)
	return h
}

func (h *Handler) CreateBattle(ctx context.Context, in *proto.CreateBattleRequest) (*proto.CreateBattleResponse, error) {
	out := &proto.CreateBattleResponse{}
	// first get game logic
	// var creator battle.GameLogicCreator

	logic, err := h.LogicCreator.CreateLogic(in.GameName, in.GameConf)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return out, err
	}

	// in.GameName
	d := table.NewTable(in.BattleConf)

	d.Start(logic)

	h.desks.Store(d.ID, d)

	out.BattleId = d.ID
	return out, nil
}

//TODO:
func (h *Handler) PlyaerJoinBattle(ctx context.Context, in *proto.EmptyMessage) (*proto.EmptyMessage, error) {
	return &proto.EmptyMessage{}, nil
}

func (h *Handler) WatcherJoinBattle(ctx context.Context, in *proto.WatcherJoinBattleRequest) (*proto.WatcherJoinBattleResponse, error) {
	out := &proto.WatcherJoinBattleResponse{}
	d := h.getDesk(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}

	//TODO:
	d.OnWatcherJoin()
	// d.OnWatcherJoin(in.PlayerInfos)
	return out, nil
}

func (h *Handler) OnUserAsyncMessage(msg *gateproto.AsyncMessageWraper) {
	// msg.UserId
	md := proto.File_servers_battle_proto_battle_proto.Messages().ByName(protoreflect.Name(msg.Topic))

	md.Options().ProtoReflect().New()
	// proto.File_servers_battle_proto_battle_proto.Messages().ByName(protoreflect.Name(msg.Name)).new
}

func (h *Handler) OnMessage(ctx context.Context, msgname string, msg protobuf.Message) {

}

func (h *Handler) OnEvent(topc string, msg protobuf.Message) {

}

func (h *Handler) OnUserConnStat() {

}

func (h *Handler) OnUserMessage(uid int64, topic string, raw []byte) {
	logger.Info("UserMessage", uid, topic)

	method := h.ct.Get(topic)
	if method == nil {
		return
	}

	// method.Call()
}

func (h *Handler) OnBattleMessage(ctx context.Context, in *proto.BattleMessageWrap) {
	d := h.getDesk(in.BattleId)
	if d == nil {
		logger.Error("desk not found battleid")
		return
	}
	d.OnBattleMessage(ctx, in)
}

// func (h *Handler) BattleMessage(ctx context.Context, stream proto.Battle_BattleMessageStream) error {
// 	defer stream.Close()
// 	for {
// 		in, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		d := h.getDesk(in.BattleId)
// 		if d == nil {
// 			return fmt.Errorf("desk not found battleid")
// 		}
// 		d.OnBattleMessage(ctx, in)
// 	}
// 	return nil
// }

func (h *Handler) getDesk(battleId string) *table.Table {
	if raw, ok := h.desks.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}
