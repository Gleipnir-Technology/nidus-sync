declare global {
	interface Window {
		bootstrap: typeof import("bootstrap");
		SSEManager: typeof import("./sse-manager").SSEManager;
	}
}

export {};
