package jwtauth

import (
	"context"
	"errors"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"skyvault/pkg/validate"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer   = "skyvault"
	audience = "skyvault"
	leeway   = time.Minute
)

var signingMethod = jwt.SigningMethodHS256

var _ auth.Claims = (*Claims)(nil)

type Claims struct {
	ProfileID int64 `json:"profileId"`
	jwt.RegisteredClaims
}

func (c *Claims) GetProfileID() int64 {
	return c.ProfileID
}

type Config struct {
	TokenTimeoutMin int
	Key             []byte
}

var _ auth.Authenticator = (*JWTAuth)(nil)

type JWTAuth struct {
	cfg Config
}

func NewJWTAuth(cfg Config) *JWTAuth {
	return &JWTAuth{cfg: cfg}
}

func (a *JWTAuth) GenerateToken(ctx context.Context, profileID int64) (string, error) {
	now := time.Now().UTC()
	expirationTime := time.Duration(a.cfg.TokenTimeoutMin) * time.Minute

	claims := &Claims{
		ProfileID: profileID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expirationTime)),
		},
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	tokenStr, err := token.SignedString(a.cfg.Key)
	if err != nil {
		return "", apperror.NewAppError(err, "JWTAuth.GenerateToken:SignedString")
	}
	return tokenStr, nil
}

func (a *JWTAuth) ValidateToken(ctx context.Context, tokenStr string) (auth.Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return a.cfg.Key, nil
		},
		jwt.WithAudience(audience),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithJSONNumber(),
		jwt.WithLeeway(leeway),
		jwt.WithValidMethods([]string{signingMethod.Name}),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperror.NewAppError(apperror.ErrAuthTokenExpired, "JWTAuth.ValidateToken:ParseWithClaims.TokenExpired")
		}

		return nil, apperror.NewAppError(apperror.ErrAuthInvalidToken, "JWTAuth.ValidateToken:ParseWithClaims")
	}
	if !token.Valid {
		return nil, apperror.NewAppError(apperror.ErrAuthInvalidToken, "JWTAuth.ValidateToken:InvalidToken")
	}
	return claims, nil
}

func (a *JWTAuth) ValidateCredentials(ctx context.Context, credentials map[auth.CredKey]any) error {
	passwordHash := *(credentials[auth.CredKeyPasswordHash].(*string))
	password := *(credentials[auth.CredKeyPassword].(*string))

	if p, err := validate.PasswordLen(password); err != nil {
		return apperror.NewAppError(err, "JWTAuth.ValidateCredentials:ValidatePassword")
	} else {
		password = p
	}

	ok, err := utils.SamePassword(passwordHash, password)
	if err != nil {
		return apperror.NewAppError(err, "JWTAuth.ValidateCredentials:IsValidPassword")
	}

	if !ok {
		return apperror.NewAppError(apperror.ErrAuthInvalidCredentials, "JWTAuth.ValidateCredentials")
	}

	return nil
}
