package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type Client struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var clients = make(map[string]*Client)
var mu sync.Mutex

func getClient(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    client, exists := clients[ip]
    if !exists {
        limiter := rate.NewLimiter(rate.Every(30*time.Second), 15) // 1 request per second with a burst of 5 requests
        clients[ip] = &Client{limiter, time.Now()}
        return limiter
    }

    client.lastSeen = time.Now()
    return client.limiter
}

func cleanupClients() {
    for {
        time.Sleep(time.Minute)
        mu.Lock()
        for ip, client := range clients {
			// Check if the client has not been seen for more than 3 minutes
            if time.Since(client.lastSeen) > 3*time.Minute {
                delete(clients, ip)
            }
        }
        mu.Unlock()
    }
}

func RateLimiter() gin.HandlerFunc {
    go cleanupClients()

    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := getClient(ip)

        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "status": "fail",
                "error": "Too many requests",
            })
            return
        }

        c.Next()
    }
}
