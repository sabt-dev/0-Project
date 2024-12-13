package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/api"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"os"
)

func Routes(r *gin.Engine) {
	var version = os.Getenv("API_VERSION")

	r.POST("/api/"+version+"/signup", api.SignUp)
	r.POST("/api/"+version+"/login", api.Login)
	r.GET("/api/"+version+"/user", middleware.RequireAuthToken, api.GetUser)
	r.GET("/api/"+version+"/logout", middleware.RequireAuthToken, api.Logout)
}
