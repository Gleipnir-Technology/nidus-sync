// Define types for the SSE data structure
interface SSEMessage {
  type: string;
  count?: number;
  [key: string]: any; // Allow additional properties
}

type SSEHandler = (data: SSEMessage) => void;

interface SSEManagerType {
  connect: (url: string) => Promise<EventSource>;
  disconnect: () => void;
  subscribe: (eventType: string, handler: SSEHandler) => void;
  unsubscribe: (eventType: string, handler: SSEHandler) => void;
  ready: (callback: (eventSource: EventSource) => void) => void;
}

declare global {
  interface Window {
    SSEManager: SSEManagerType;
  }
}

export const SSEManager: SSEManagerType = (function (): SSEManagerType {
  let eventSource: EventSource | null = null;
  let subscribers: Map<string, SSEHandler[]> = new Map();
  let isConnected: boolean = false;
  let connectionPromise: Promise<EventSource> | null = null;

  function subscribe(eventType: string, handler: SSEHandler): void {
    if (!subscribers.has(eventType)) {
      subscribers.set(eventType, []);
    }
    subscribers.get(eventType)!.push(handler);

    // If already connected, attach the listener immediately
    if (isConnected && eventSource) {
      eventSource.addEventListener(eventType, handler as EventListener);
    }
  }

  function unsubscribe(eventType: string, handler: SSEHandler): void {
    if (subscribers.has(eventType)) {
      const handlers = subscribers.get(eventType)!;
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
    if (eventSource) {
      eventSource.removeEventListener(eventType, handler as EventListener);
    }
  }

  function connect(url: string): Promise<EventSource> {
    if (connectionPromise) {
      return connectionPromise;
    }

    connectionPromise = new Promise((resolve, reject) => {
      eventSource = new EventSource(url);

      eventSource.onopen = function (): void {
        isConnected = true;

        // Attach all pre-registered handlers
        subscribers.forEach((handlers: SSEHandler[], eventType: string) => {
          handlers.forEach((handler: SSEHandler) => {
            eventSource!.addEventListener("message", (message: MessageEvent) => {
              const data: SSEMessage = JSON.parse(message.data);
              if (eventType === "*" || eventType === data.type) {
                handler(data);
              }
            });
          });
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

  return {
    connect,
    disconnect,
    subscribe,
    unsubscribe,
    ready,
  };
})();

function updateNotificationBadge(data: SSEMessage): void {
  const badge = document.querySelector<HTMLElement>(".notification-badge");
  if (badge) {
    badge.textContent = String(data.count || 0);
    badge.style.display = (data.count || 0) > 0 ? "block" : "none";
  }
}
