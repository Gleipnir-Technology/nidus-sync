import * as bootstrap from "bootstrap";

declare global {
	interface Window {
		Alpine: any;
		SSEManager: any;
		createAppPlanning: any;
		bootstrap: typeof bootstrap;
	}
}

export {};
