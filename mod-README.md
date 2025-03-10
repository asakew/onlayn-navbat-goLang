Create a real-time online queuing system using Go (Gin Framework) for the backend and HTML/CSS/JavaScript for the frontend with the following requirements:

1. Backend (GoLang/goFiber):
- Implement an endpoint to handle queue registration
- Generate and assign unique sequential numbers to users who join the queue
- Maintain the current queue state and position tracking
- Provide WebSocket support for real-time updates
- Include error handling for edge cases (queue full, invalid requests)

2. Frontend (HTML/CSS/JS):
- Create a clean, responsive interface showing:
    * Current queue position
    * Total people in queue
    * Estimated waiting time
    * Real-time updates when queue position changes
- Implement WebSocket connection to receive live updates
- Add visual feedback for queue status (joining, waiting, error states)
- Include a "Join Queue" button that becomes disabled once in queue

Technical Requirements:
- Use Gin's built-in middleware for CORS and request handling
- Implement proper session management to prevent duplicate queue entries
- Ensure the system can handle concurrent users
- Add basic rate limiting to prevent abuse
- Include logging for system monitoring

The solution should be scalable and maintain queue integrity during high traffic.