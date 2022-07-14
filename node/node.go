package node

import (
	"google.golang.org/protobuf/proto"

	"hotwave/event"
)

//TODO:
type NodeBase struct {
	NodeID      string
	NodeType    string
	NodeVersion string
}

func (n *NodeBase) PublishEvent(msg proto.Message) {

}

func (n *NodeBase) OnEvent(msg *event.Event) {

}
