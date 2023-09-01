package auth

type UserInfo struct {
	Uid   int64
	Uname string
	Role  string
}

type Auth interface {
	TokenAuth(token string) (*UserInfo, error)
	// AccountAuth(account string, password string) *UserInfo
}
