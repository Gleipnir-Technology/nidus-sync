import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "@/AppSync.vue";
import * as config from "@/config";
import router from "@/route/config";
import * as sentry from "@/sentry";

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
sentry.Init(app, pinia).then(() => {
	app.mount("#app");
});
