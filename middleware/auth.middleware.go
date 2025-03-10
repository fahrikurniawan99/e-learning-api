package middleware

import (
	"elearning_api/repository"
	"elearning_api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	DecodeToken() gin.HandlerFunc
}

type AuthMiddlewareImplement struct {
	TokenRepository repository.TokenRepository
}

func NewAuthMiddleware(tokenRepository repository.TokenRepository) AuthMiddleware {
	return &AuthMiddlewareImplement{
		TokenRepository: tokenRepository,
	}
}

func (am *AuthMiddlewareImplement) DecodeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		accessTokenFind, err := am.TokenRepository.FindAccessToken(authHeader)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when find access token : " + err.Error()})
			c.Abort()
			return
		}

		if accessTokenFind.ID == 0 {
			tokenLog, err := am.TokenRepository.FindAccessTokenLog(authHeader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when find access token log : " + err.Error()})
				c.Abort()
				return
			}

			if tokenLog.ID == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "your account has been used on another device"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.VerifyToken(tokenStr, false)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error when verify token : " + err.Error()})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}
