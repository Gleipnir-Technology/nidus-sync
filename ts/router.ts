import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "./view/Home.vue";
import About from "./view/About.vue";
import Intelligence from "./view/Intelligence.vue";

const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
	},
	{
		path: "/intelligence",
		name: "Intelligence",
		component: Intelligence,
	},
];

const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

export default router;
