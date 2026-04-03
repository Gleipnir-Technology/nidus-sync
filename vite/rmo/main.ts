import { createApp } from "vue";
import { createHead } from "@vueuse/head";
import "@/gen/custom-icons.scss";
import "@/style/rmo.scss";
import router from "@/rmo/router";
import App from "@/rmo/App.vue";

const app = createApp(App);
const head = createHead();

app.use(head);
app.use(router);
app.mount("#app");
