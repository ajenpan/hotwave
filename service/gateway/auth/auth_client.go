package auth

import (
	"context"
	"crypto/rsa"

	log "hotwave/logger"
	authproto "hotwave/service/auth/proto"
	"hotwave/service/common"
)

type AuthClient struct {
	PK         *rsa.PublicKey
	authClient authproto.AuthClient
}

func (a *AuthClient) TokenAuth(token string) *UserInfo {
	uid, uname, err := common.VerifyToken(a.PK, token)
	if err != nil {
		return nil
	}
	return &UserInfo{
		Uid:   uid,
		Uname: uname,
	}
}

func (a *AuthClient) AccountAuth(account string, password string) *UserInfo {
	req := &authproto.LoginRequest{
		Uname:  account,
		Passwd: password,
	}
	resp, err := a.authClient.Login(context.Background(), req)
	if err != nil {
		log.Errorf("auth client login error:%v", err)
		return nil
	}
	ret := &UserInfo{
		Uid:   resp.UserInfo.Uid,
		Uname: resp.UserInfo.Uname,
	}
	return ret
}
