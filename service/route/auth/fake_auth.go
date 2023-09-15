package auth

import "sync/atomic"

type FakeAuth struct {
	c uint64
}

func (a *FakeAuth) TokenAuth(token string) (*UserInfo, error) {
	return &UserInfo{
		UID:   a.nextID(),
		UName: token,
	}, nil
}

func (a *FakeAuth) AccountAuth(account string, password string) *UserInfo {
	return &UserInfo{
		UID:   a.nextID(),
		UName: account,
		Role:  "user",
	}
}

func (a *FakeAuth) nextID() uint64 {
	ret := atomic.AddUint64(&a.c, 1)
	if ret == 0 {
		ret = atomic.AddUint64(&a.c, 1)
	}
	return ret
}
