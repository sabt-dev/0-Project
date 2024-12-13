package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/routes"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	//initializers.SyncDatabase() // Uncomment this line to create the table
}

func main() {
	router := gin.Default()

	// Apply the CORS configuration to the router
	router.Use(middleware.CORSMiddleware())

	// Apply the routes to the router
	routes.Routes(router)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
