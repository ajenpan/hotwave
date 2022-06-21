package handler

import (
	"context"
	"fmt"
	"sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	// "github.com/google/uuid"
	"hotwave/logger"
	"hotwave/services/battle"
	"hotwave/services/battle/proto"
	"hotwave/services/battle/table"
	gateproto "hotwave/services/gateway/proto"
)

type Handler struct {
	// proto.UnimplementedBattleServer
	// gateproto.UnimplementedGateAdapterServer

	desks        sync.Map
	LogicCreator *battle.GameLogicCreator
}

func New() *Handler {
	return &Handler{
		LogicCreator: &battle.GameLogicCreator{},
	}
}

func (h *Handler) CreateBattle(ctx context.Context, in *proto.CreateBattleRequest) (*proto.CreateBattleResponse, error) {
	out := &proto.CreateBattleResponse{}
	// first get game logic
	// var creator battle.GameLogicCreator

	// if c, ok := h.LoigcCreators.Load(in.GameName); !ok {
	// 	return fmt.Errorf("not found game by name:%s", in.GameName)
	// } else {
	// 	creator, ok = c.(battle.GameLogicCreator)
	// 	if !ok {
	// 		return fmt.Errorf("")
	// 	}
	// }
	logic, err := h.LogicCreator.CreateLogic(in.GameName, "")
	//get from pool?

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

func (h *Handler) UserMessage(ctx context.Context, in *gateproto.AsyncMessageWraper) (*gateproto.SteamClosed, error) {
	logger.Info("UserMessage", in.UserId, in.Topic)
	return &gateproto.SteamClosed{}, nil
}

func (h *Handler) OnBattleMessage(ctx context.Context, in *proto.BattleMessageWrap) {
	d := h.getDesk(in.BattleId)
	if d == nil {
		// fmt.Errorf("desk not found battleid")
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
