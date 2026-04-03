import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "@/rmo/view/Home.vue";
import Nuisance from "@/rmo/view/Nuisance.vue";
import Water from "@/rmo/view/Water.vue";
const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
	},
	{
		path: "/nuisance",
		name: "Nuisance",
		component: Nuisance,
	},
	{
		path: "/water",
		name: "Water",
		component: Water,
	},
];

export const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

export default router;
