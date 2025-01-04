package auth

// import (
// 	"skyvault/pkg/common"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type JWT struct {
// 	app *common.App
// }

// func NewJWT(app *common.App) *JWT {
// 	return &JWT{app: app}
// }

// type Claims struct {
// 	UserID string `json:"user_id"`
// 	jwt.RegisteredClaims
// }

// func (j *JWT) generateToken(userID string) (string, error) {
// 	expirationTime := time.Now().Add(time.Duration(j.app.Config.AUTH_JWT_TOKEN_TIMEOUT_MIN) * time.Minute)
// 	claims := &Claims{
// 		UserID: userID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwtSecret)
// }

// func ValidateJWT(tokenStr string) (*Claims, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, jwt.ErrSignatureInvalid
// 	}

// 	return claims, nil
// }
