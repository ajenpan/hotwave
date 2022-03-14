package transfer

import (
	"hotwave/frame"
	protocol "hotwave/servers/gateway/proto"
)

type Adpater interface {
	OnUserMessage(frame.User, *protocol.UserMessageWraper)
}
type Transfer struct {
	protocol.UnimplementedGateAdpaterServer

	// Adapter
}

func (t *Transfer) UserMessage(svr *protocol.UserMessageWraper) error {

	return nil
}
