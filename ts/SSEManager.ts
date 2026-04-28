// Define types for the SSE data structure
export interface SSEMessageBase {
	type: string;
}
export interface SSEMessageResource extends SSEMessageBase {
	resource: string;
	time: string;
	uri: string;
}

export interface SSEMessageStatus extends SSEMessageBase {
	build_time: Date;
	is_modified: boolean;
	revision: string;
	status: string;
}
type SSEHandlerResource = (data: SSEMessageResource) => void;
type SSEHandlerStatus = (data: SSEMessageStatus) => void;

interface SSEManagerType {
	connect: (url: string) => Promise<EventSource>;
	disconnect: () => void;
	subscribe: (handler: SSEHandlerResource) => string;
	subscribeStatus: (handler: SSEHandlerStatus) => string;
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
	let subscribersResource: Map<string, SSEHandlerResource> = new Map();
	let subscribersStatus: Map<string, SSEHandlerStatus> = new Map();
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
					const data: SSEMessageBase = JSON.parse(message.data);
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
					console.log("SSE reconnecting");
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
			console.log("SSE disconnected");
		}
	}

	function handleMessage(msg: SSEMessageBase) {
		if (msg.type == "heartbeat") {
			return;
		} else if (msg.type == "status") {
			subscribersStatus.forEach((handler: SSEHandlerStatus, _: string) => {
				handler(msg as SSEMessageStatus);
			});
		} else {
			subscribersResource.forEach((handler: SSEHandlerResource, _: string) => {
				handler(msg as SSEMessageResource);
			});
		}
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

	function subscribe(handler: SSEHandlerResource): string {
		const uuid = crypto.randomUUID();
		subscribersResource.set(uuid.toString(), handler);
		return uuid;
	}

	function subscribeStatus(handler: SSEHandlerStatus): string {
		const uuid = crypto.randomUUID();
		subscribersStatus.set(uuid.toString(), handler);
		return uuid;
	}

	function unsubscribe(uuid: string): void {
		if (subscribersResource.has(uuid)) {
			subscribersResource.delete(uuid);
		}
		if (subscribersStatus.has(uuid)) {
			subscribersStatus.delete(uuid);
		}
	}

	return {
		connect,
		disconnect,
		subscribe,
		subscribeStatus,
		unsubscribe,
		ready,
	};
})();
