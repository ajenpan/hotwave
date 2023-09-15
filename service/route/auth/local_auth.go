package auth

import (
	"crypto/rsa"

	"hotwave/service/auth/common"
)

type LocalAuth struct {
	PK *rsa.PublicKey
}

func (a *LocalAuth) TokenAuth(token string) (*UserInfo, error) {
	return common.VerifyToken(a.PK, token)
}

func (a *LocalAuth) AccountAuth(account string, password string) *UserInfo {
	return nil
}
