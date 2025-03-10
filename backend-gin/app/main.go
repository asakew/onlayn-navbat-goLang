package main

import (
	"log"
	"net/http"
	_ "net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"queue_system/internal/handlers"
	"queue_system/internal/middleware"
	"queue_system/internal/models"
)

func main() {
	// Set up the queue manager
	queueManager := models.NewQueueManager()

	// Start broadcasting queue updates
	go queueManager.BroadcastQueueState()

	// Set up Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimiter())
	router.Use(middleware.SessionMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static files for web
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/queue/status", handlers.GetQueueStatus(queueManager))
		api.POST("/queue/join", handlers.JoinQueue(queueManager))
		api.POST("/queue/advance", handlers.AdvanceQueue(queueManager)) // Admin endpoint
		api.GET("/ws", handlers.HandleWebsocket(queueManager))
	}

	// Start server
	log.Println("Starting server on http://localhost:8081/")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
