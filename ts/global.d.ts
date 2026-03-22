import * as bootstrap from "bootstrap";

declare global {
	interface Window {
		SSEManager: SSEManagerType;
		bootstrap: typeof bootstrap;
	}
}

export {};
