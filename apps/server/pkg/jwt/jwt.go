package jwt

import (
	"fmt"
	"time"

	"fullstack-engineering-lab/server/internal/config"

	goJwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	goJwt.RegisteredClaims
}

func GenerateToken(userID uint, cfg *config.JWTConfig) (string, error) {
	expireDuration := time.Duration(cfg.ExpireMinutes) * time.Minute
	claims := Claims{
		UserID: userID,
		RegisteredClaims: goJwt.RegisteredClaims{
			ExpiresAt: goJwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  goJwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := goJwt.ParseWithClaims(tokenString, &Claims{}, func(token *goJwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*goJwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
