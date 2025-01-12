package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/config"
	"github.com/sabt-dev/0-Project/internal/handlers"
	"github.com/sabt-dev/0-Project/internal/middleware"
)

func Routes(router *gin.Engine) {
	public := router.Group("/api/" + config.AppConfig.APIv) // Public routes

	public.POST("/auth/register", handlers.Register)
	public.POST("/auth/login", handlers.Login)
	public.POST("/auth/request-password-reset", handlers.RequestPasswordReset)
    public.PUT("/auth/reset-password", handlers.ResetPassword)
	public.GET("/auth/verify-email", handlers.VerifyUserEmail)
	public.POST("/auth/refresh-token", handlers.RefreshToken)

	protected := router.Group("/api/" + config.AppConfig.APIv) // Protected/Private routes

	protected.GET("/users/me", middleware.RequireAuthToken, handlers.GetUser)
	protected.GET("/auth/logout", middleware.RequireAuthToken, handlers.Logout)
}
