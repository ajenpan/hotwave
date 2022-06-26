package hotwave

import (
	"hotwave/proto"
)

type NoopAdpater struct{}

func (a *NoopAdpater) OnUserMessage(User, *proto.UserMessageWraper)  {}
func (a *NoopAdpater) OnNodeEvent(string, *proto.EventMessageWraper) {}
