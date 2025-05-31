package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // ❗ เปลี่ยนเป็น ENV จริงในโปรดักชัน

func GenerateToken(staffID uint, hospitalID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"staff_id":    staffID,
		"hospital_id": hospitalID,
		"role":        role,
		"exp":         time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
