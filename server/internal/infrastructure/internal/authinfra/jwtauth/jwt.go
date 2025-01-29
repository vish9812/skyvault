package jwtauth

import (
	"context"
	"errors"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer   = "skyvault"
	audience = "skyvault"
	leeway   = 2 * time.Minute
)

var signingMethod = jwt.SigningMethodHS256

var _ auth.Claims = (*Claims)(nil)

type Claims struct {
	ProfileID int64  `json:"profileId"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func (c *Claims) GetProfileID() int64 {
	return c.ProfileID
}

func (c *Claims) GetEmail() string {
	return c.Email
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

func (a *JWTAuth) GenerateToken(ctx context.Context, profileID int64, email string) (string, error) {
	now := time.Now().UTC()
	expirationTime := time.Duration(a.cfg.TokenTimeoutMin) * time.Minute

	claims := &Claims{
		ProfileID: profileID,
		Email:     email,
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
			return nil, apperror.NewAppError(apperror.ErrAuthTokenExpired, "JWTAuth.ValidateToken:ParseWithClaims:TokenExpired")
		}

		return nil, apperror.NewAppError(apperror.ErrAuthInvalidToken, "JWTAuth.ValidateToken:ParseWithClaims")
	}
	if !token.Valid {
		return nil, apperror.NewAppError(apperror.ErrAuthInvalidToken, "JWTAuth.ValidateToken:InvalidToken")
	}
	return claims, nil
}

func (a *JWTAuth) ValidateCredentials(ctx context.Context, credentials map[auth.CredsKeys]any) error {
	passwordHashPtr, ok := credentials[auth.CredsKeysPasswordHash]
	if !ok {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "JWTAuth.ValidateCredentials:PasswordHash")
	}

	passwordHash, ok := passwordHashPtr.(string)
	if !ok {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "JWTAuth.ValidateCredentials:PasswordHash:InvalidType")
	}

	passwordPtr, ok := credentials[auth.CredsKeysPassword]
	if !ok {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "JWTAuth.ValidateCredentials:Password")
	}

	password, ok := passwordPtr.(string)
	if !ok {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "JWTAuth.ValidateCredentials:Password:InvalidType")
	}

	ok, err := utils.IsValidPassword(passwordHash, password)
	if err != nil {
		return apperror.NewAppError(err, "JWTAuth.ValidateCredentials:IsValidPassword")
	}

	if !ok {
		return apperror.NewAppError(apperror.ErrAuthInvalidCredentials, "JWTAuth.ValidateCredentials")
	}

	return nil
}
