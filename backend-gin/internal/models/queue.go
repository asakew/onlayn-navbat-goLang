package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// User represents a person in the queue
type User struct {
	ID       string    `json:"id"`
	Number   int       `json:"number"`
	JoinTime time.Time `json:"join_time"`
}

// QueueUpdate represents a message sent to clients
type QueueUpdate struct {
	TotalInQueue    int `json:"total_in_queue"`
	CurrentPosition int `json:"current_position"`
	CurrentNumber   int `json:"current_number"`
}

// QueueManager handles the queue state
type QueueManager struct {
	Users           map[string]*User
	CurrentNumber   int
	CurrentPosition int
	Mu              sync.Mutex
	Clients         map[*websocket.Conn]bool
	Broadcast       chan QueueUpdate
}

// NewQueueManager creates a new queue manager instance
func NewQueueManager() *QueueManager {
	return &QueueManager{
		Users:           make(map[string]*User),
		CurrentNumber:   0,
		CurrentPosition: 0,
		Clients:         make(map[*websocket.Conn]bool),
		Broadcast:       make(chan QueueUpdate),
	}
}

// BroadcastQueueState continuously broadcasts queue updates to all connected clients
func (qm *QueueManager) BroadcastQueueState() {
	for {
		update := <-qm.Broadcast
		qm.Mu.Lock()
		for client := range qm.Clients {
			err := client.WriteJSON(update)
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(qm.Clients, client)
			}
		}
		qm.Mu.Unlock()
	}
}

// AddUser adds a new user to the queue
func (qm *QueueManager) AddUser(sessionID string) (*User, bool) {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()

	// Check if user is already in queue
	if user, exists := qm.Users[sessionID]; exists {
		return user, false
	}

	// Add user to queue
	qm.CurrentNumber++
	user := &User{
		ID:       sessionID,
		Number:   qm.CurrentNumber,
		JoinTime: time.Now(),
	}
	qm.Users[sessionID] = user

	// Broadcast update
	qm.Broadcast <- QueueUpdate{
		TotalInQueue:    len(qm.Users),
		CurrentPosition: qm.CurrentPosition,
		CurrentNumber:   qm.CurrentNumber,
	}

	return user, true
}

// AdvanceQueue moves the queue forward by one position
func (qm *QueueManager) AdvanceQueue() int {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()

	qm.CurrentPosition++

	// Broadcast update
	qm.Broadcast <- QueueUpdate{
		TotalInQueue:    len(qm.Users),
		CurrentPosition: qm.CurrentPosition,
		CurrentNumber:   qm.CurrentNumber,
	}

	return qm.CurrentPosition
}

// GetQueueStatus returns the current status for a specific user
func (qm *QueueManager) GetQueueStatus(sessionID string) map[string]interface{} {
	qm.Mu.Lock()
	defer qm.Mu.Unlock()

	user, inQueue := qm.Users[sessionID]

	if !inQueue {
		return map[string]interface{}{
			"in_queue":       false,
			"total_in_queue": len(qm.Users),
			"current_number": qm.CurrentNumber,
		}
	}

	// Calculate position in queue
	position := user.Number - qm.CurrentPosition
	if position < 1 {
		position = 0
	}

	// Estimate waiting time (2 minutes per person)
	estimatedTime := position * 2

	return map[string]interface{}{
		"in_queue":       true,
		"queue_number":   user.Number,
		"position":       position,
		"total_in_queue": len(qm.Users),
		"estimated_time": estimatedTime,
	}
}
