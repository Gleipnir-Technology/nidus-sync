import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "./view/Home.vue";
import About from "./view/About.vue";
import Communication from "./view/Communication.vue";
import ConfigurationIntegration from "./view/configuration/Integration.vue";
import ConfigurationIntegrationArcgis from "./view/configuration/IntegrationArcgis.vue";
import ConfigurationOrganization from "./view/configuration/Organization.vue";
import ConfigurationPesticide from "./view/configuration/Pesticide.vue";
import ConfigurationPesticideAdd from "./view/configuration/PesticideAdd.vue";
import ConfigurationRoot from "./view/configuration/Root.vue";
import ConfigurationUpload from "./view/configuration/Upload.vue";
import ConfigurationUploadDetail from "./view/configuration/UploadDetail.vue";
import ConfigurationUploadPool from "./view/configuration/UploadPool.vue";
import ConfigurationUploadPoolFlyover from "./view/configuration/UploadPoolFlyover.vue";
import ConfigurationUser from "./view/configuration/User.vue";
import ConfigurationUserAdd from "./view/configuration/UserAdd.vue";
import Intelligence from "./view/Intelligence.vue";
import NotFound from "./view/NotFound.vue";
import OAuthRefreshArcgis from "./view/OAuthRefreshArcgis.vue";
import Operations from "./view/Operations.vue";
import Planning from "./view/Planning.vue";
import ReviewPool from "./view/review/Pool.vue";
import ReviewRoot from "./view/review/Root.vue";
import Signin from "./view/Signin.vue";
import Sudo from "./view/Sudo.vue";
import apiClient from "@/client";

const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/communication",
		name: "Communication",
		component: Communication,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration",
		name: "Configuration",
		component: ConfigurationRoot,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/integration",
		name: "Integration Configuration",
		component: ConfigurationIntegration,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/integration/arcgis",
		name: "Arcgis Integration Configuration",
		component: ConfigurationIntegrationArcgis,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/organization",
		name: "Organization Configuration",
		component: ConfigurationOrganization,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/pesticide",
		name: "Pesticide Configuration",
		component: ConfigurationPesticide,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/pesticide/add",
		name: "Pesticide Add",
		component: ConfigurationPesticideAdd,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/upload",
		name: "Upload Configuration",
		component: ConfigurationUpload,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		component: ConfigurationUploadDetail,
		meta: { requiresAuth: true, showSidebar: true },
		name: "Upload Detail",
		path: "/_/configuration/upload/:id",
		props: true,
	},
	{
		path: "/_/configuration/upload/pool",
		name: "Pool Upload",
		component: ConfigurationUploadPool,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/upload/pool/flyover",
		name: "Flyover Upload",
		component: ConfigurationUploadPoolFlyover,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/user",
		name: "User Configuration",
		component: ConfigurationUser,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/configuration/user/add",
		name: "User Add Configuration",
		component: ConfigurationUserAdd,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/intelligence",
		name: "Intelligence",
		component: Intelligence,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/oauth/refresh/arcgis",
		name: "Arcgis OAuth Refresh",
		component: OAuthRefreshArcgis,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/operations",
		name: "Operations",
		component: Operations,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/planning",
		name: "Planning",
		component: Planning,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/review",
		name: "Review",
		component: ReviewRoot,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/review/pool",
		name: "Pool Review",
		component: ReviewPool,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/_/signin",
		name: "Signin",
		component: Signin,
		meta: { requiresAuth: false, showSidebar: false },
	},
	{
		path: "/_/sudo",
		name: "Sudo",
		component: Sudo,
		meta: { requiresAuth: true, showSidebar: true },
	},
	// Catch-all route - must be last
	{
		path: "/:pathMatch(.*)*",
		name: "NotFound",
		component: NotFound,
	},
];

export const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

// Global navigation guard
router.beforeEach(async (to, from) => {
	const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);

	if (requiresAuth) {
		try {
			// Check if user is authenticated (could be an API call)
			const isAuthenticated = await apiClient.isAuthenticated();
			if (!isAuthenticated) {
				return "/signin";
			} else {
				return;
			}
		} catch (error) {
			console.log("check auth failed");
			return "/signin";
		}
	}
});

export default router;
