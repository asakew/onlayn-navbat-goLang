Design a comprehensive online queue management system that implements the following core features:

## Backend (GoLang/GIN Framework):

1. Queue Management System
- Create a robust API for queue creation, deletion, and modification
- Implement user authentication and role-based access control
- Design endpoints for adding, removing, and updating queue entries
- Include queue configuration options (max capacity, priority levels, service categories)

2. Real-time Queue Management
- Implement WebSocket connections for live updates
- Create bi-directional communication between server and clients
- Enable instant notifications for queue status changes
- Design a heartbeat mechanism to maintain connection stability
- Handle reconnection scenarios gracefully

3. Queue Position Tracking
- Develop a position calculation algorithm considering priority levels
- Create status update notifications (waiting, processing, completed)
- Implement estimated wait time calculations
- Design a dashboard for queue analytics and metrics
- Enable customer notification system (SMS/email) for status updates

Technical Requirements:
- Use WebSocket protocol for real-time communications
- Implement horizontal scaling capabilities
- Design for high availability and fault tolerance
- Include rate limiting and security measures
- Ensure data consistency across distributed systems
- Implement caching mechanisms for improved performance
- Create comprehensive logging and monitoring systems

Deliverables:
- System architecture diagram
- API documentation
- WebSocket event specifications
- Database schema
- Performance benchmarks
- Scaling strategy documentation