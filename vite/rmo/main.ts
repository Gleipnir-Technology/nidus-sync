import { createApp } from "vue";
import { createHead } from "@vueuse/head";
import { createPinia } from "pinia";
import "@/gen/custom-icons.scss";
import "@/style/rmo.scss";
import router from "@/rmo/router";
import App from "@/rmo/App.vue";

const app = createApp(App);
const head = createHead();
const pinia = createPinia();

app.use(head);
app.use(pinia);
app.use(router);
app.mount("#app");
