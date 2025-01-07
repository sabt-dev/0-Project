package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sabt-dev/0-Project/internal/config"
	"github.com/sabt-dev/0-Project/internal/models"
)

func GenerateToken(user *models.User, c *gin.Context) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"adr": c.ClientIP(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}