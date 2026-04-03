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

document.addEventListener("DOMContentLoaded", () => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
});

const pinia = createPinia();
const app = createApp(App);
app.use(pinia);
app.use(router);
app.mount("#app");
