document.addEventListener('DOMContentLoaded', function() {
    // DOM elements
    const joinQueueBtn = document.getElementById('join-queue');
    const tryAgainBtn = document.getElementById('try-again');
    const initialCard = document.getElementById('initial-card');
    const waitingCard = document.getElementById('waiting-card');
    const errorCard = document.getElementById('error-card');
    const totalCount = document.getElementById('total-count');
    const queueNumber = document.getElementById('queue-number');
    const position = document.getElementById('position');
    const totalInQueue = document.getElementById('total-in-queue');
    const estimatedTime = document.getElementById('estimated-time');
    const progressBar = document.getElementById('progress-bar');
    const errorMessage = document.getElementById('error-message');
    const currentTimeEl = document.getElementById('current-time');

    // Update current time
    function updateCurrentTime() {
        const now = new Date();
        const formatted = now.toISOString().replace('T', ' ').substring(0, 19);
        currentTimeEl.textContent = formatted;
    }

    // Initial time update
    updateCurrentTime();

    // Update time every minute
    setInterval(updateCurrentTime, 60000);

    // WebSocket connection
    let socket;

    // Connect to WebSocket
    function connectWebSocket() {
        // Close existing connection if any
        if (socket) {
            socket.close();
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws`;

        socket = new WebSocket(wsUrl);

        socket.onopen = () => {
            console.log('WebSocket connection established');
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);

            // Update total count on initial screen
            totalCount.textContent = data.total_in_queue;

            // Update queue status if we're in the waiting card
            if (waitingCard.classList.contains('active')) {
                fetchQueueStatus();
            }
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        socket.onclose = (event) => {
            console.log('WebSocket connection closed', event);

            // Try to reconnect after 5 seconds
            setTimeout(connectWebSocket, 5000);
        };

        // Send a ping every 30 seconds to keep connection alive
        setInterval(() => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000);
    }

    // Show a specific card
    function showCard(card) {
        initialCard.classList.remove('active');
        waitingCard.classList.remove('active');
        errorCard.classList.remove('active');

        card.classList.add('active');
    }

    // Fetch current queue status
    function fetchQueueStatus() {
        fetch('/api/queue/status')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to get queue status');
                }
                return response.json();
            })
            .then(data => {
                if (data.in_queue) {
                    // Show waiting card with user's queue info
                    showCard(waitingCard);
                    queueNumber.textContent = data.queue_number;
                    position.textContent = data.position;
                    totalInQueue.textContent = data.total_in_queue;
                    estimatedTime.textContent = data.estimated_time;

                    // Update progress bar
                    const maxPosition = data.total_in_queue;
                    const userPosition = data.position;

                    if (maxPosition > 0) {
                        const progressPercentage = Math.max(
                            0,
                            Math.min(
                                100,
                                ((maxPosition - userPosition) / maxPosition) * 100
                            )
                        );
                        progressBar.style.width = `${progressPercentage}%`;
                    }

                    // If position is 0, the user is at the front
                    if (data.position === 0) {
                        estimatedTime.textContent = "It's your turn!";
                    }
                } else {
                    // Show initial join screen
                    showCard(initialCard);
                    totalCount.textContent = data.total_in_queue;
                }
            })
            .catch(error => {
                console.error('Error fetching queue status:', error);
                errorMessage.textContent = 'Failed to get queue status. Please try again.';
                showCard(errorCard);
            });
    }

    // Join the queue
    function joinQueue() {
        // Disable the button to prevent multiple clicks
        joinQueueBtn.disabled = true;

        fetch('/api/queue/join', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(data => {
                        throw new Error(data.error || 'Failed to join queue');
                    });
                }
                return response.json();
            })
            .then(data => {
                fetchQueueStatus();
            })
            .catch(error => {
                console.error('Error joining queue:', error);
                errorMessage.textContent = error.message;
                showCard(errorCard);
                joinQueueBtn.disabled = false;
            });
    }

    // Event listeners
    joinQueueBtn.addEventListener('click', joinQueue);

    tryAgainBtn.addEventListener('click', function() {
        showCard(initialCard);
        joinQueueBtn.disabled = false;
    });

    // Initial setup
    connectWebSocket();
    fetchQueueStatus();

    // Periodically check queue status
    setInterval(fetchQueueStatus, 10000);
});