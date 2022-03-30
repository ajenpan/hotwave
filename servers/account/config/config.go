package config

type Config struct {
	UserDBDSN         string
	TokenExpireMinute int
	JwtKey            string

	PasswdErrDetectSecond int
	PasswdErrLimit        int
	PasswdErrBlockSecond  int
}

var DefaultConf = &Config{
	JwtKey:            "my_secret_key",
	TokenExpireMinute: 12,
	UserDBDSN:         "root:123456@tcp(localhost:3306)/usercenter?charset=utf8mb4&parseTime=True&loc=Local",
}
