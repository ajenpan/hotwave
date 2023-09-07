package common

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// iss	Issuer	发行方
// sub	Subject	 主题
// aud	Audience	 受众
// exp	Expiration Time	过期时间
// nbf	Not Before	早于该定义的时间的JWT不能被接受处理
// iat	Issued At	JWT发行时的时间戳
// jti	JWT ID	JWT的唯一标识
// uid	用户ID
// rid	角色ID

func GenerateToken(pk *rsa.PrivateKey, uid int64, uname string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["iss"] = "hotwave"
	claims["sub"] = "auth"
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()
	claims["aud"] = uname
	claims["uid"] = strconv.FormatInt(uid, 10)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(pk)
}

func VerifyToken(pk *rsa.PublicKey, tokenRaw string) (int64, string, error) {
	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (interface{}, error) {
		return pk, nil
	})
	if err != nil {
		return 0, "", err
	}
	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}
	uname := claims["aud"]
	uidstr := claims["uid"]
	uid, _ := strconv.ParseInt(uidstr.(string), 10, 64)
	return uid, uname.(string), nil
}
