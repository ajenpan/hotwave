package matcher

import (
	"hotwave/service/lobby/user"
	"sync"
)

type MatchQue struct {
	expext int32
}

type OnMatchedFunc func(Matcher, []int64)
type OnTimeoutFunc func(Matcher)

type Matcher interface {
	String() string
}

type StaticMatcher struct {
	rwlock sync.RWMutex
	users  []*user.UserInfo
}

//TODO:
func (sm *StaticMatcher) Add(u *user.UserInfo) {
	sm.rwlock.Lock()
	defer sm.rwlock.Unlock()

	sm.users = append(sm.users, u)

	const expert = 4

	if len(sm.users) >= expert {
		//create battle

	}
}
