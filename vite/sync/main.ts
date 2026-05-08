import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "@/AppSync.vue";
import * as config from "@/config";
import router from "@/route/config";
import * as sentry from "@/sentry";
import { useErrorHandler } from "@/composable/error-handler";

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

const { setError } = useErrorHandler();

const pinia = createPinia();
const app = createApp(App);
app.config.errorHandler = (err, instance, info) => {
	// err: the error object
	// instance: the component instance where error occurred
	// info: Vue-specific error info, e.g., lifecycle hook

	console.error("Global error:", err);
	console.error("Error info:", info);
	console.error("Error instance:", instance);

	// You could dispatch to a store, send to error tracking service, etc.
	// For example, trigger a global error state
	setError(err);
};

app.use(pinia);
app.use(router);
sentry.Init(app, pinia).then(() => {
	app.mount("#app");
});
