package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/routes"
	"github.com/sabt-dev/0-Project/internal/jobs"
	"log"
	"os"
)

var addr string

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase() // Uncomment this line to create the table
	initializers.CheckTLSFilesExistence()
}

func main() {
	router := gin.Default()
	
	// Apply the CORS configuration to the router
	router.Use(middleware.AllowCORS())
	
	// Apply the rate limiter middleware to the router
	router.Use(middleware.RateLimiter())

	// Start the cleanup job
    go jobs.CleanupUnverifiedUsers()
	
	// Apply the routes to the router
	routes.Routes(router)
	
	// Run the server with TLS
	addr = ":" + os.Getenv("PORT")
	err := router.RunTLS(addr, "../tls/cert.pem", "../tls/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
