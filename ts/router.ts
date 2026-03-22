import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "./view/Home.vue";
import About from "./view/About.vue";
import Communication from "./view/Communication.vue";
import Configuration from "./view/Configuration.vue";
import Intelligence from "./view/Intelligence.vue";
import Operations from "./view/Operations.vue";
import Planning from "./view/Planning.vue";
import Review from "./view/Review.vue";
import Sudo from "./view/Sudo.vue";

const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
	},
	{
		path: "/communication",
		name: "Communication",
		component: Communication,
	},
	{
		path: "/configuration",
		name: "Configuration",
		component: Configuration,
	},
	{
		path: "/intelligence",
		name: "Intelligence",
		component: Intelligence,
	},
	{
		path: "/operations",
		name: "Operations",
		component: Operations,
	},
	{
		path: "/planning",
		name: "Planning",
		component: Planning,
	},
	{
		path: "/review",
		name: "Review",
		component: Review,
	},
	{
		path: "/sudo",
		name: "Sudo",
		component: Sudo,
	},
];

const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

export default router;
