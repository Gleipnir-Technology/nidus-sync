import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import HomeBase from "@/rmo/view/Home.vue";
import HomeDistrict from "@/rmo/view/district/Home.vue";
import NuisanceBase from "@/rmo/view/Nuisance.vue";
import NuisanceDistrict from "@/rmo/view/district/Nuisance.vue";
import StatusBase from "@/rmo/view/Status.vue";
import StatusDistrict from "@/rmo/view/district/Status.vue";
import Water from "@/rmo/view/Water.vue";
import WaterDistrict from "@/rmo/view/district/Water.vue";
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
	{
		path: "/district/:slug/nuisance",
		name: "NuisanceDistrict",
		component: NuisanceDistrict,
		props: true,
	},
	{
		path: "/district/:slug/status",
		name: "StatusDistrict",
		component: StatusDistrict,
		props: true,
	},
	{
		path: "/district/:slug/water",
		name: "WaterDistrict",
		component: WaterDistrict,
		props: true,
	},
	{
		path: "/status",
		name: "StatusBase",
		component: StatusBase,
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
