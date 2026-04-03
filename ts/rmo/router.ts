import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import HomeBase from "@/rmo/view/Home.vue";
import HomeDistrict from "@/rmo/view/district/Home.vue";
import NuisanceBase from "@/rmo/view/Nuisance.vue";
//import * as NuisanceDistrict from "@/rmo/view/district/Nuisance.vue";
import Status from "@/rmo/view/Status.vue";
import Water from "@/rmo/view/Water.vue";
const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "HomeBase",
		component: HomeBase,
	},
	{
		path: "/nuisance",
		name: "NuisanceBase",
		component: NuisanceBase,
	},
	{
		path: "/district/:slug",
		name: "HomeDistrict",
		component: HomeDistrict,
		props: true,
	},
	/*{
		path: "/district/{slug}/nuisance",
		name: "NuisanceDistrict",
		component: NuisanceDistrict,
		props: true,
	},*/
	{
		path: "/status",
		name: "Status",
		component: Status,
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
