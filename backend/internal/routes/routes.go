package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/api"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"os"
)

func Routes(r *gin.Engine) {
	var version string = os.Getenv("API_VERSION")

	r.POST("/api/"+version+"/auth/register", api.Register)
	r.POST("/api/"+version+"/auth/login", api.Login)
	r.GET("/api/"+version+"/auth/verifyemail/:code", api.VerifyEmail)
	r.GET("/api/"+version+"/users/me", middleware.RequireAuthToken, api.GetUser)
	r.GET("/api/"+version+"/auth/logout", middleware.RequireAuthToken, api.Logout)
}
