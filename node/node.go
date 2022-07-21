package node

import (
	"time"

	"google.golang.org/protobuf/proto"

	"hotwave/event"
	log "hotwave/logger"
	"hotwave/utils/calltable"
)

type NodeBase struct {
	NodeID      string
	NodeName    string
	NodeVersion string

	publisher event.Publisher
}

func (n *NodeBase) PublishEvent(msg proto.Message) {
	raw, err := proto.Marshal(msg)
	if err != nil {
		log.Error("proto.Marshal error:", err)
	}
	event := &event.Event{
		Topic:     string(proto.MessageName(msg)),
		Data:      raw,
		FromNode:  n.NodeID,
		Timestamp: time.Now().Unix(),
	}

	n.publisher.Publish(event)
}

func (n *NodeBase) SubEvent(nodeid string, topics []string, ct *calltable.CallTable) {

}
