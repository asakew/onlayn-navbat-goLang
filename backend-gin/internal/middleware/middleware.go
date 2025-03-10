package middleware

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateUUID generates a unique ID for session management
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// SessionMiddleware checks for a session cookie and sets it if not present
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			sessionID = GenerateUUID()
			c.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
		}
		c.Set("sessionID", sessionID)
		c.Next()
	}
}

// RateLimiter implements a simple rate limiting middleware
func RateLimiter() gin.HandlerFunc {
	// Using a map to store IP addresses and their request counts
	limitMap := make(map[string]int)
	var mu sync.Mutex

	// Reset the map every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for k := range limitMap {
				delete(limitMap, k)
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		count := limitMap[ip]
		if count >= 60 { // 60 requests per minute
			mu.Unlock()
			c.JSON(429, gin.H{
				"error": "Rate limit exceeded, please try again later",
			})
			c.Abort()
			return
		}
		limitMap[ip] = count + 1
		mu.Unlock()
		c.Next()
	}
}

// Logger middleware for request logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details after processing
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		fmt.Printf("[%s] %s %s | %d | %v | %s\n",
			endTime.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
		)
	}
}
