import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "./view/Home.vue";
import About from "./view/About.vue";
import Communication from "./view/Communication.vue";
import ConfigurationOrganization from "./view/configuration/Organization.vue";
import ConfigurationPesticide from "./view/configuration/Pesticide.vue";
import ConfigurationPesticideAdd from "./view/configuration/PesticideAdd.vue";
import ConfigurationRoot from "./view/configuration/Root.vue";
import ConfigurationUser from "./view/configuration/User.vue";
import ConfigurationUserAdd from "./view/configuration/UserAdd.vue";
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
		component: ConfigurationRoot,
	},
	{
		path: "/configuration/organization",
		name: "Organization Configuration",
		component: ConfigurationOrganization,
	},
	{
		path: "/configuration/pesticide",
		name: "Pesticide Configuration",
		component: ConfigurationPesticide,
	},
	{
		path: "/configuration/pesticide/add",
		name: "Pesticide Add",
		component: ConfigurationPesticideAdd,
	},
	{
		path: "/configuration/user",
		name: "User Configuration",
		component: ConfigurationUser,
	},
	{
		path: "/configuration/user/add",
		name: "User Add Configuration",
		component: ConfigurationUserAdd,
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
