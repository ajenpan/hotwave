package user

import (
	"hotwave/node/usergater"
	protocal "hotwave/service/lobby/proto"
)

type UserInfo struct {
	UserId    int64
	UserProps map[int32]*protocal.PropsInfo
	Socket    usergater.UserSocket
}
