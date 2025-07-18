package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sabt-dev/0-Project/internal/config"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
)

func RequireAuthToken(c *gin.Context) {
	// get the token from the header
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.AppConfig.JWTSecret), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check if the token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error": "Token has expired",
			})
		}

		// find the user_id from the token in "sub"
		userID, err := uuid.Parse(claims["sub"].(string))
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        var user models.User
        result := initializers.DB.First(&user, "id = ?", userID)
        if result.Error != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

		// check if the user has the same IP-address
		//if claims["ip"] != c.ClientIP() {
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		// attach the user to the context
		c.Set("userData", user)

		// continue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
