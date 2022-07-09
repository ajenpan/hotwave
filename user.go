package hotwave

import "google.golang.org/protobuf/proto"

type User interface {
	UID() int64
	SendMessage(proto.Message) error
}