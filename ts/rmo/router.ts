import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "@/rmo/view/Home.vue";
const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
	},
];

export const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

export default router;
