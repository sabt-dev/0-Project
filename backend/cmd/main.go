package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	//initializers.SyncDatabase() // Uncomment this line to create the table
}

func main() {
	router := gin.Default()
	var addr string = ":"+os.Getenv("PORT")

	// Apply the CORS configuration to the router
	router.Use(middleware.AllowCORS())

	// Apply the routes to the router
	routes.Routes(router)

	err := router.RunTLS(addr, "../tls/cert.pem", "../tls/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
