package battle

import (
	protobuf "google.golang.org/protobuf/proto"
	// "hotwave/user"
)

type Player interface {
	// user.User
	GetSeatID() int32
	GetScore() int64 //game jetton
	IsRobot() bool

	SendMessage(protobuf.Message) error
}
