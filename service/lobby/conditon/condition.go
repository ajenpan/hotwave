package conditon

import "hotwave/service/lobby/user"

type Condition interface {
	Check(*user.UserInfo) bool
}
