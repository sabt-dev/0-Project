package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/handlers"
	"github.com/sabt-dev/0-Project/internal/middleware"
)

func Routes(router *gin.Engine) {
	r := router.Group("/api/" + os.Getenv("API_VERSION"))

	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.POST("/auth/request-password-reset", handlers.RequestPasswordReset)
    r.POST("/auth/reset-password", handlers.ResetPassword)
	r.GET("/auth/verify-email", handlers.VerifyUserEmail)
	r.GET("/users/me", middleware.RequireAuthToken, handlers.GetUser)
	r.GET("/auth/logout", middleware.RequireAuthToken, handlers.Logout)
}
