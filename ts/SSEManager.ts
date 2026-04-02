// Define types for the SSE data structure
export interface SSEMessage {
	resource: string;
	time: string;
	type: string;
	uri: string;
}

type SSEHandler = (data: SSEMessage) => void;

interface SSEManagerType {
	connect: (url: string) => Promise<EventSource>;
	disconnect: () => void;
	subscribe: (handler: SSEHandler) => string;
	unsubscribe: (uuid: string) => void;
	ready: (callback: (eventSource: EventSource) => void) => void;
}

/*
declare global {
	interface Window {
		SSEManager: SSEManagerType;
	}
}
*/

export const SSEManager: SSEManagerType = (function (): SSEManagerType {
	let eventSource: EventSource | null = null;
	let subscribers: Map<string, SSEHandler> = new Map();
	let isConnected: boolean = false;
	let connectionPromise: Promise<EventSource> | null = null;

	function connect(url: string): Promise<EventSource> {
		if (connectionPromise) {
			return connectionPromise;
		}

		connectionPromise = new Promise((resolve, reject) => {
			eventSource = new EventSource(url);

			eventSource.onopen = function (): void {
				isConnected = true;

				eventSource!.addEventListener("message", (message: MessageEvent) => {
					const data: SSEMessage = JSON.parse(message.data);
					handleMessage(data);
				});

				console.log("SSE connected");
				resolve(eventSource!);
			};

			eventSource.onerror = function (err: Event): void {
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

	function disconnect(): void {
		if (eventSource) {
			eventSource.close();
			eventSource = null;
			isConnected = false;
			connectionPromise = null;
		}
	}

	function handleMessage(msg: SSEMessage) {
		if (msg.type == "heartbeat") {
			return;
		}
		subscribers.forEach((handler: SSEHandler, _: string) => {
			handler(msg);
		});
	}

	function ready(callback: (eventSource: EventSource) => void): void {
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

	function subscribe(handler: SSEHandler): string {
		const uuid = crypto.randomUUID();
		subscribers.set(uuid.toString(), handler);
		return uuid;
	}

	function unsubscribe(uuid: string): void {
		if (subscribers.has(uuid)) {
			subscribers.delete(uuid);
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
