package auth

import (
	"crypto/rsa"

	"hotwave/service/auth/common"
)

type LocalAuth struct {
	PK *rsa.PublicKey
}

func (a *LocalAuth) TokenAuth(token string) (*UserInfo, error) {

	uid, uname, err := common.VerifyToken(a.PK, token)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Uid:   uid,
		Uname: uname,
	}, nil
}

var c = int64(0)

func (a *LocalAuth) AccountAuth(account string, password string) *UserInfo {
	c++
	return &UserInfo{
		Uid:   c,
		Uname: account,
	}
}
