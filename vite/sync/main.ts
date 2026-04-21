import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "@/AppSync.vue";
import router from "@/router";
import * as Sentry from "@sentry/vue";
import * as config from "@/config";

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

const pinia = createPinia();
const app = createApp(App);
app.use(pinia);
app.use(router);
app.mount("#app");

Sentry.init({
	dsn: config.DSN,
	integrations: [Sentry.browserTracingIntegration({ router })],
	environment: config.ENVIRONMENT,
	release: config.RELEASE,
	tracesSampleRate: 0.01,
});
