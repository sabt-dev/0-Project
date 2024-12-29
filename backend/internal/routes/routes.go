package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/api"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/utils"
)

func Routes(router *gin.Engine) {
	r := router.Group("/api/" + os.Getenv("API_VERSION"))

	r.POST("/auth/register", api.Register)
	r.POST("/auth/login", api.Login)
	r.GET("/auth/verifyemail/:code", utils.VerifyEmail)
	r.GET("/users/me", middleware.RequireAuthToken, api.GetUser)
	r.GET("/auth/logout", middleware.RequireAuthToken, api.Logout)
}
