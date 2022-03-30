package cache

import (
	"context"

	"hotwave/servers/account/database/models"
)

type Noop struct {
}

func (Noop) StoreUser(ctx context.Context, user *AuthCacheInfo) error { return nil }

func (Noop) FetchUser(ctx context.Context, uid int64) *AuthCacheInfo {
	return &AuthCacheInfo{
		User: &models.User{Id: uid},
	}
}

func (Noop) FetchUserByName(ctx context.Context, uname string) *AuthCacheInfo {
	return &AuthCacheInfo{
		User: &models.User{Name: uname},
	}
}
