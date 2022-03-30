package battle

import (
	protobuf "google.golang.org/protobuf/proto"
	// "hotwave/user"
)

type Player interface {
	// user.User
	GetSeatID() SeatID
	GetScore() int64 //game jetton
	IsRobot() bool

	SendMessage(protobuf.Message) error
}
