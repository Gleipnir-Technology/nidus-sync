import { createApp } from "vue";
import { createHead } from "@vueuse/head";
import { createPinia } from "pinia";
import "bootstrap-icons/font/bootstrap-icons.css";
import "@/gen/custom-icons.scss";
import "@/style/rmo.scss";
import router from "@/rmo/route/config";
import App from "@/rmo/App.vue";

const app = createApp(App);
const head = createHead();
const pinia = createPinia();

app.use(head);
app.use(pinia);
app.use(router);
app.mount("#app");
