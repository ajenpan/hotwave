package usergater

import "hotwave/node/transport"

type User interface {
	UID() int64 // return user id
}

type UserSocket interface {
	User
	transport.Socket
}
