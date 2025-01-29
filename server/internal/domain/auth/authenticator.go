package auth

import (
	"context"
	"skyvault/pkg/common"
)

type Claims interface {
	GetProfileID() int64
	GetEmail() string
}

func GetClaimsFromContext(ctx context.Context) Claims {
	return ctx.Value(common.CtxKeyClaims).(Claims)
}

func GetProfileIDFromContext(ctx context.Context) int64 {
	return GetClaimsFromContext(ctx).GetProfileID()
}

func GetEmailFromContext(ctx context.Context) string {
	return GetClaimsFromContext(ctx).GetEmail()
}

type CredsKeys string

const (
	CredsKeysPassword     CredsKeys = "password"
	CredsKeysPasswordHash CredsKeys = "password_hash"
)

type Authenticator interface {
	GenerateToken(ctx context.Context, profileID int64, email string) (string, error)

	// App Errors:
	// - apperror.ErrInvalidToken
	// - apperror.ErrTokenExpired
	ValidateToken(ctx context.Context, token string) (Claims, error)

	// App Errors:
	// - apperror.ErrInvalidValue
	// - apperror.ErrInvalidCredentials
	ValidateCredentials(ctx context.Context, credentials map[CredsKeys]any) error
}

type AuthenticatorFactory interface {
	// App Errors:
	// - apperror.ErrInvalidProvider
	GetAuthenticator(provider Provider) (Authenticator, error)
}
