import * as bootstrap from "bootstrap";

declare global {
	interface Window {
		Alpine: any;
		SSEManager: any;
		bootstrap: typeof bootstrap;
	}
}

export {};
