package cache

import (
	"context"
	"time"

	"hotwave/servers/account/database/models"
)

type AuthCacheInfo struct {
	User         *models.User
	AssessToken  string
	RefreshToken string
	LoginAt      string
	LoginIP      string
}

type AuthCache interface {
	StoreUser(ctx context.Context, user *AuthCacheInfo, expireAt time.Duration) error
	DeleteUser(ctx context.Context, uid int64)
	FetchUser(ctx context.Context, uid int64) *AuthCacheInfo
	FetchUserByName(ctx context.Context, uname string) *AuthCacheInfo
	FetchUserByToken(ctx context.Context, token string) *AuthCacheInfo
}
