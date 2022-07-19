package client

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
	"hotwave/utils/calltable"
)

type HandleDeliver struct {
	H       interface{}
	SvrName string
	NodeId  string
	CT      *calltable.CallTable
}

func (h *HandleDeliver) OnMessage(u transport.Session, msg *protocal.GateMessage) error {
	method := h.CT.Get(msg.Name)
	if method == nil {
		return fmt.Errorf("method %s not found", msg.Name)
	}
	req := method.NewRequest().(proto.Message)
	proto.Unmarshal(msg.Body, req)

	callResult := method.Call(h.H, u, req)
	if len(callResult) != 2 {
		return fmt.Errorf("method %s return no result", msg.Name)
	}
	var callErr error
	var callResp proto.Message

	if !callResult[0].IsNil() {
		callResp = callResult[0].Interface().(proto.Message)
	}
	if !callResult[1].IsNil() {
		callErr = callResult[1].Interface().(error)
	}
	wrap := &protocal.GateMessage{
		Name: msg.Name,
	}
	if callResp != nil {
		fmt.Println(callErr)
	}
	if callResp != nil {
		wrap.Body, _ = proto.Marshal(callResp)
	}
	return nil
}
