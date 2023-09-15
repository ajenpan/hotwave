package handler

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	protobuf "google.golang.org/protobuf/proto"

	"hotwave/event"
	log "hotwave/logger"
	battle "hotwave/service/battle"
	"hotwave/service/battle/proto"
	"hotwave/service/battle/table"
	"hotwave/transport/tcp"
	"hotwave/utils/calltable"
	"hotwave/utils/marshal"
)

type Handler struct {
	battles sync.Map

	LogicCreator *battle.GameLogicCreator
	ct           *calltable.CallTable[uint32]
	marshal      marshal.Marshaler
	Publisher    event.Publisher

	createCounter int32
}

func New() *Handler {
	h := &Handler{
		LogicCreator: &battle.GameLogicCreator{},
	}
	h.ct = calltable.ExtractAsyncMethodByMsgID(proto.File_proto_battle_server_proto.Messages(), h)
	return h
}

func (h *Handler) CreateBattle(ctx context.Context, in *proto.StartBattleRequest) (*proto.StartBattleResponse, error) {
	logic, err := h.LogicCreator.CreateLogic(in.GameName)
	if err != nil {
		return nil, err
	}

	atomic.AddInt32(&h.createCounter, 1)

	battleid := uuid.NewString() + fmt.Sprintf("-%d", h.createCounter)

	d := table.NewTable(table.TableOption{
		ID:             battleid,
		Conf:           in.BattleConf,
		EventPublisher: h.Publisher,
		FinishReporter: func() {
			h.onBattleFinished(battleid)
		},
	})

	players, err := table.NewPlayers(in.PlayerInfos)
	if err != nil {
		return nil, err
	}

	err = d.Init(logic, players, in.BattleConf)
	if err != nil {
		return nil, err
	}

	h.battles.Store(battleid, d)

	out := &proto.StartBattleResponse{
		BattleId: d.ID,
	}
	return out, nil
}

// func (h *Handler) StopBattle(ctx context.Context, in *proto.StopBattleRequest) (*proto.StopBattleResponse, error) {
// 	out := &proto.StopBattleResponse{}
// 	d := h.getBattleById(in.BattleId)
// 	if d == nil {
// 		return out, fmt.Errorf("battle not found")
// 	}
// 	d.Close()
// 	h.battles.Delete(in.BattleId)
// 	return out, nil
// }

type userInfoKey struct{}

var UserInfoKey = &userInfoKey{}

func GetUserInfo(ctx context.Context) *tcp.UserInfo {
	return ctx.Value(UserInfoKey).(*tcp.UserInfo)
}

func WithUserInfo(ctx context.Context, uinfo *tcp.UserInfo) context.Context {
	return context.WithValue(ctx, UserInfoKey, uinfo)
}

func (h *Handler) onBattleFinished(battleid string) {
	d := h.getBattleById(battleid)
	if d == nil {
		return
	}
	d.Close()
	h.battles.Delete(battleid)
}

func (h *Handler) OnEvent(topc string, msg protobuf.Message) {

}

func (h *Handler) bingBattle() error {

	return nil
}

func (h *Handler) JoinBattle(ctx context.Context, in *proto.JoinBattleRequest) (*proto.JoinBattleResponse, error) {
	out := &proto.JoinBattleResponse{
		BattleId:   in.BattleId,
		SeatId:     in.SeatId,
		ReadyState: in.ReadyState,
	}

	d := h.getBattleById(in.BattleId)
	if d == nil {
		return nil, fmt.Errorf("battle not found")
	}
	uinfo := GetUserInfo(ctx)
	d.OnPlayerReady(uinfo.Uid, in.ReadyState)
	return out, nil
}

func (h *Handler) OnBattleMessageWrap(s *tcp.Socket, msg *proto.LoigcMessageWrap) {
	b := h.getBattleById(msg.BattleId)
	if b == nil {
		return
	}
	b.OnPlayerMessage(s.Uid, (msg.Msgid), msg.Data)
}

func (h *Handler) getBattleById(battleId string) *table.Table {
	if raw, ok := h.battles.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}

func (h *Handler) OnConn(s *tcp.Socket, ss tcp.SocketStat) {
	log.Info("OnConn:", int(ss))
	if ss == tcp.Disconnected {
		s.MetaLoad("uid")
	}
}

func (h *Handler) OnMessage(s *tcp.Socket, ss *tcp.THVPacket) {
	ctype := ss.GetType()
	if ctype <= tcp.PacketTypeInnerEndAt_ {
		return
	}

	if ctype == 4 {
		body := ss.GetBody()
		if len(body) < 4 {
			log.Errorf("invalid message, body len: %d", len(body))
			return
		}

		msgid := binary.LittleEndian.Uint32(body)
		method := h.ct.Get(msgid)
		if method == nil {
			return
		}

		req := method.GetRequest()
		defer method.PutRequest(req)

		err := h.marshal.Unmarshal(body[4:], req)
		if err != nil {
			log.Errorf("marshal msgid:%d,error:%w", msgid, err)
			return
		}

		res := method.Call(req)
		if len(res) == 0 {
			return
		}
		if res[0].IsNil() {
			return
		}
		err, ok := res[0].Interface().(error)
		if ok && err != nil {
			log.Errorf("call msgid:%d,error:%w", msgid, err)
			return
		}

		// ss.Reset()
		// respi := res[0]
		// respraw, err := h.marshal.Marshal(respi)
		// if err != nil {
		// 	return
		// }
		// respHead := make([]byte, 4)
		// binary.LittleEndian.PutUint32(respHead, msgid+1)
		// ss.Body = append(respHead, respraw...)
		// s.SendPacket(ss)
	}
}
