package usergater

import "hotwave/node/transport"

type User interface {
	UID() int64 // return user id
	SetUID(int64)
}

type UserSocket interface {
	User
	transport.Socket
}

type UserGater interface {
	Addr() string
	Close() error
	Accept(func(UserSocket)) error
}
