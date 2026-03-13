// sse-manager.js - Include this in your common template
window.SSEManager = (function () {
	let eventSource = null;
	let subscribers = new Map();
	let isConnected = false;
	let connectionPromise = null;

	function subscribe(eventType, handler) {
		if (!subscribers.has(eventType)) {
			subscribers.set(eventType, []);
		}
		subscribers.get(eventType).push(handler);

		// If already connected, attach the listener immediately
		if (isConnected && eventSource) {
			eventSource.addEventListener(eventType, handler);
		}
	}

	function unsubscribe(eventType, handler) {
		if (subscribers.has(eventType)) {
			const handlers = subscribers.get(eventType);
			const index = handlers.indexOf(handler);
			if (index > -1) {
				handlers.splice(index, 1);
			}
		}
		if (eventSource) {
			eventSource.removeEventListener(eventType, handler);
		}
	}

	function connect(url) {
		if (connectionPromise) {
			return connectionPromise;
		}

		connectionPromise = new Promise((resolve, reject) => {
			eventSource = new EventSource(url);

			eventSource.onopen = function () {
				isConnected = true;

				// Attach all pre-registered handlers
				subscribers.forEach((handlers, eventType) => {
					handlers.forEach((handler) => {
						eventSource.addEventListener("message", (message) => {
							const data = JSON.parse(message.data);
							if (eventType == "*" || eventType == data.type) {
								handler(data);
							}
						});
					});
				});

				console.log("SSE connected");
				resolve(eventSource);
			};

			eventSource.onerror = function (err) {
				console.error("SSE error:", err);
				isConnected = false;

				// Close old connection
				if (eventSource) {
					eventSource.close();
				}

				// Reconnect after delay
				setTimeout(() => {
					connectionPromise = null;
					connect(url);
				}, 5000);

				if (!isConnected) {
					reject(err);
				}
			};
		});

		return connectionPromise;
	}

	function disconnect() {
		if (eventSource) {
			eventSource.close();
			eventSource = null;
			isConnected = false;
			connectionPromise = null;
		}
	}

	function ready(callback) {
		if (connectionPromise) {
			connectionPromise.then(callback);
		} else {
			// If connect hasn't been called yet, queue it
			const checkInterval = setInterval(() => {
				if (connectionPromise) {
					clearInterval(checkInterval);
					connectionPromise.then(callback);
				}
			}, 50);
		}
	}

	return {
		connect,
		disconnect,
		subscribe,
		unsubscribe,
		ready,
	};
})();

// Initialize SSE for navigation notifications
document.addEventListener("DOMContentLoaded", function () {
	SSEManager.connect("/api/events");
});

function updateNotificationBadge(data) {
	const badge = document.querySelector(".notification-badge");
	if (badge) {
		badge.textContent = data.count;
		badge.style.display = data.count > 0 ? "block" : "none";
	}
}
