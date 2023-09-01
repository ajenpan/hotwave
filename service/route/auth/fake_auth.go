package auth

import "sync/atomic"

type FakeAuth struct {
	c uint32
}

func (a *FakeAuth) TokenAuth(token string) (*UserInfo, error) {
	c++
	return &UserInfo{
		Uid:   c,
		Uname: token,
	}, nil
}

func (a *FakeAuth) AccountAuth(account string, password string) *UserInfo {
	c++
	return &UserInfo{
		Uid:   c,
		Uname: account,
		Role:  "user",
	}
}

func (a *FakeAuth) nextID() uint32 {
	ret := atomic.AddUint32(&a.c, 1)
	if ret == 0 {
		ret = atomic.AddUint32(&a.c, 1)
	}
	return ret
}
