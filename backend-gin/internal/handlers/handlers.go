package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"queue_system/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now, can be restricted in production
	},
}

// HandleWebsocket upgrades the HTTP connection to a WebSocket connection
func HandleWebsocket(qm *models.QueueManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %+v", err)
			return
		}

		// Register the new client
		qm.Mu.Lock()
		qm.Clients[conn] = true
		qm.Mu.Unlock()

		// Send current state to the new client
		qm.Broadcast <- models.QueueUpdate{
			TotalInQueue:    len(qm.Users),
			CurrentPosition: qm.CurrentPosition,
			CurrentNumber:   qm.CurrentNumber,
		}

		// Handle disconnects
		defer func() {
			qm.Mu.Lock()
			delete(qm.Clients, conn)
			qm.Mu.Unlock()
			conn.Close()
		}()

		// Simple ping-pong to keep connection alive
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Echo the message back as a pong
			if messageType == websocket.PingMessage {
				if err := conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
					break
				}
			}
		}
	}
}

// JoinQueue handles requests to join the queue
func JoinQueue(qm *models.QueueManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("sessionID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session not found"})
			return
		}

		user, isNew := qm.AddUser(sessionID.(string))
		if !isNew {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "Already in queue",
				"queue_number": user.Number,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"queue_number": user.Number,
			"message":      "You have joined the queue",
		})
	}
}

// GetQueueStatus retrieves the current queue status for a user
func GetQueueStatus(qm *models.QueueManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("sessionID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session not found"})
			return
		}

		status := qm.GetQueueStatus(sessionID.(string))
		c.JSON(http.StatusOK, status)
	}
}

// AdvanceQueue moves the queue forward (admin endpoint)
func AdvanceQueue(qm *models.QueueManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real system, you would add authentication here
		// For simplicity, we're not implementing authentication in this example

		newPosition := qm.AdvanceQueue()
		c.JSON(http.StatusOK, gin.H{
			"success":          true,
			"current_position": newPosition,
		})
	}
}
