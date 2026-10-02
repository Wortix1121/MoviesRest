package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

func NewJWTManager(secretKey []byte, ttl time.Duration) TokenManager {
	return &JWTManager{
		secretKey: secretKey,
		ttl:       ttl,
	}

}

func (j *JWTManager) GenerateToken(userID int64) (string, error) {
	const op = "pkg.auth.GenerateToken"

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signed, nil
}

func (j *JWTManager) ValidateToken(tokenStr string) (int64, error) {
	const op = "pkg.auth.ValidateToken"

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("%s: invalid token", op)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("%s: invalid claims type", op)
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("%s: user_id claim missing or invalid", op)
	}

	return int64(userIDFloat), nil
}
