package handler

import (
	"context"
	"testing"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/servers/battle"
	"hotwave/servers/battle/noop"
	pb "hotwave/servers/battle/proto"
	gatewayproto "hotwave/servers/gateway/proto"
)

func createTestHandler() *Handler {
	h := New()
	h.LogicCreator = &battle.GameLogicCreator{}
	h.LogicCreator.Store("noop", noop.NewGameLogic)
	return h
}
func TestCreateBattle(t *testing.T) {
	h := createTestHandler()

	ctx := context.Background()
	in := &pb.CreateBattleRequest{
		GameName: "noop",
	}
	// out := &pb.CreateBattleResponse{}
	h.CreateBattle(ctx, in)
}
func newMessageWarp(msg protobuf.Message) *gatewayproto.AsyncMessageWraper {
	ret := &gatewayproto.AsyncMessageWraper{
		// Gateway: "",
		// MsgName: string(protobuf.MessageName(msg)),
		UserId: 1,
	}
	ret.Body, _ = protobuf.Marshal(msg)
	return ret
}

func TestOnUserMessage(t *testing.T) {
	h := createTestHandler()

	warper := newMessageWarp(&pb.CreateBattleRequest{})

	h.OnUserAsyncMessage(warper)

	// md := pb.File_servers_battle_proto_battle_proto.Messages().ByName(protoreflect.Name(warper.Name))
	// msg := md.Options().ProtoReflect().New()
	// msg:= pb.File_servers_battle_proto_battle_proto.Options().ProtoReflect().New()

	// protobuf.Unmarshal(warper.Body, msg)

	// mt,_:= preg.GlobalTypes.FindMessageByName(protoreflect.FullName(md.FullName()))
	// mt.New()

	// h.OnUserMessage()
	// h.OnUserMessage()
}
