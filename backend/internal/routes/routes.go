package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/api"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/utils"
)

func Routes(r *gin.Engine) {
	var version string = os.Getenv("API_VERSION")

	r.POST("/api/"+version+"/auth/register", api.Register)
	r.POST("/api/"+version+"/auth/login", api.Login)
	r.GET("/api/"+version+"/auth/verifyemail/:code", utils.VerifyEmail)
	r.GET("/api/"+version+"/users/me", middleware.RequireAuthToken, api.GetUser)
	r.GET("/api/"+version+"/auth/logout", middleware.RequireAuthToken, api.Logout)
}
