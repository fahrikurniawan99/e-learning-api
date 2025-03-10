package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

var (
	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
)

func GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func VerifyToken(tokenStr string, isRefresh bool) (*jwt.Token, error) {
	secret := accessSecret
	if isRefresh {
		secret = refreshSecret
	}
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
