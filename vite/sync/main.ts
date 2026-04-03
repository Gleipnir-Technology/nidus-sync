import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "@/AppSync.vue";
import router from "@/router";
import { SSEManager, type SSEMessage } from "@/SSEManager";
import "maplibre-gl/dist/maplibre-gl.css";

// Import Bootstrap Icons CSS
import "bootstrap-icons/font/bootstrap-icons.css";
// Import Bootstrap SCSS
import "@/style/style.scss";
// Import custom icons
import "@/gen/custom-icons.scss";

// Import Bootstrap JavaScript and make it available globally
import * as bootstrap from "bootstrap";
window.bootstrap = bootstrap;

// Make SSEManager available to all the JavaScript
window.SSEManager = SSEManager;

document.addEventListener("DOMContentLoaded", () => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
});
document.addEventListener("init", () => {
	const user = {
		display_name: "",
		initials: "",
		notifications: [],
		notification_counts: {
			communication: 0,
			home: 0,
			review: 0,
		},
		organization: {
			name: "",
		},
		role: "",
		username: "",
	};
});
interface GreetingComponent {
	message: string;
	name: string;
	updateMessage(): void;
}

const pinia = createPinia();
const app = createApp(App);
app.use(pinia);
app.use(router);
app.mount("#app");
