package auth

import (
	"crypto/rsa"

	"hotwave/service/common"
)

type LocalAuth struct {
	PK *rsa.PublicKey
}

func (a *LocalAuth) TokenAuth(token string) *UserInfo {
	if len(token) < 10 {
		return nil
	}

	uid, uname, err := common.VerifyToken(a.PK, token)
	if err != nil {
		return nil
	}

	return &UserInfo{
		Uid:   uid,
		Uname: uname,
	}
}

func (a *LocalAuth) AccountAuth(account string, password string) *UserInfo {
	return &UserInfo{
		Uid:   1,
		Uname: account,
	}
}
