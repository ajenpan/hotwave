package auth

import (
	"crypto/rsa"

	"hotwave/service/auth/common"
)

type AuthClient struct {
	PK *rsa.PublicKey
}

func (a *AuthClient) TokenAuth(token string) *UserInfo {
	uc, err := common.VerifyToken(a.PK, token)
	if err != nil {
		return nil
	}
	return uc
}

func (a *AuthClient) AccountAuth(account string, password string) *UserInfo {
	return nil
}
