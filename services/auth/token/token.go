package token

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"hotwave/services/auth/database/models"
)

type Token = jwt.Token

type TokenGenerater interface {
	Generate(ctx context.Context, UserInfo *models.Users) (token *Token)
	SignedString(ctx context.Context, token *Token) string
	Parse(ctx context.Context, raw string) *Token
}
