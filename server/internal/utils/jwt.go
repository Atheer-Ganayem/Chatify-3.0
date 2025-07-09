package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(token string) (string, error) {
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		panic("AUTH_SECRET is a required env variable.")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	userId := string(claims["id"].(string))
	if userId == "" {
		return "", errors.New("Invalid token")
	}

	return userId, nil
}