package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sabt-dev/0-Project/internal/config"
	"github.com/sabt-dev/0-Project/internal/initializers"
	// "github.com/sabt-dev/0-Project/internal/services"
	"github.com/sabt-dev/0-Project/internal/middleware"
	"github.com/sabt-dev/0-Project/internal/routes"
)

func init() {
	config.LoadConfig()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	initializers.CheckTLSFilesExistence()
}

func main() {
	router := gin.Default()
	
	// Apply the CORS configuration to the router
	router.Use(middleware.AllowCORS())
	
	// Apply the rate limiter middleware to the router
	router.Use(middleware.RateLimiter())

	// Start the cleanup job
    // go services.CleanupUnverifiedUsers()
	
	// Apply the routes to the router
	routes.Routes(router)
	
	// Run the server with TLS
	addr := ":" + config.AppConfig.Port
	err := router.RunTLS(addr, "../tls/cert.pem", "../tls/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
