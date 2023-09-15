package auth

import "hotwave/service/auth/common"

// type UserInfo struct {
// 	Uid   uint64
// 	Uname string
// 	Role  string
// }
type UserInfo = common.UserClaims
type Auth interface {
	TokenAuth(token string) (*UserInfo, error)
	// AccountAuth(account string, password string) *UserInfo
}
