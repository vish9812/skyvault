package auth

import (
	"errors"
	"skyvault/pkg/common"
	"skyvault/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid jwt token")

type JWT struct {
	app            *common.App
	jwtKey         []byte
	expirationTime time.Duration
}

func NewAuthJWT(app *common.App) *JWT {
	return &JWT{
		app:            app,
		jwtKey:         []byte(app.Config.AUTH_JWT_KEY),
		expirationTime: time.Duration(app.Config.AUTH_JWT_TOKEN_TIMEOUT_MIN) * time.Minute,
	}
}

type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func (a *JWT) Generate(userID int64, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "skyvault",
			Audience:  []string{"skyvault"},
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(a.expirationTime)),
			Subject:   email,
			ID:        utils.UUID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtKey)
}

// Claims validates the token and returns the claims
func (a *JWT) Claims(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	}, jwt.WithAudience("skyvault"), jwt.WithIssuer("skyvault"), jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithJSONNumber(), jwt.WithLeeway(2*time.Minute), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return nil, errors.New("failed to parse with claims")
	}
	if !token.Valid {
		return nil, common.NewValidationError(ErrInvalidToken)
	}
	return claims, nil
}
