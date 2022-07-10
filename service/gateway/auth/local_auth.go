package auth

import (
	"crypto/rsa"
	"hotwave/service/common"
)

type LocalAuth struct {
	PK *rsa.PublicKey
}

func (a *LocalAuth) TokenAuth(token string) *UserSession {
	uid, uname, err := common.VerifyToken(a.PK, token)
	if err != nil {
		return nil
	}
	return &UserSession{
		uid:   uid,
		uname: uname,
	}
}

func (a *LocalAuth) AccountAuth(account string, password string) *UserSession {
	return &UserSession{
		uid:   1,
		uname: account,
	}
}
