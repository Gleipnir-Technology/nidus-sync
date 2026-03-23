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
import ConfigurationUser from "./view/configuration/User.vue";
import ConfigurationUserAdd from "./view/configuration/UserAdd.vue";
import Intelligence from "./view/Intelligence.vue";
import NotFound from "./view/NotFound.vue";
import OAuthRefreshArcgis from "./view/OAuthRefreshArcgis.vue";
import Operations from "./view/Operations.vue";
import Planning from "./view/Planning.vue";
import Review from "./view/Review.vue";
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
		path: "/communication",
		name: "Communication",
		component: Communication,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration",
		name: "Configuration",
		component: ConfigurationRoot,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/integration",
		name: "Integration Configuration",
		component: ConfigurationIntegration,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/integration/arcgis",
		name: "Arcgis Integration Configuration",
		component: ConfigurationIntegrationArcgis,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/organization",
		name: "Organization Configuration",
		component: ConfigurationOrganization,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/pesticide",
		name: "Pesticide Configuration",
		component: ConfigurationPesticide,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/pesticide/add",
		name: "Pesticide Add",
		component: ConfigurationPesticideAdd,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/user",
		name: "User Configuration",
		component: ConfigurationUser,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/configuration/user/add",
		name: "User Add Configuration",
		component: ConfigurationUserAdd,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/intelligence",
		name: "Intelligence",
		component: Intelligence,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/oauth/refresh/arcgis",
		name: "Arcgis OAuth Refresh",
		component: OAuthRefreshArcgis,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/operations",
		name: "Operations",
		component: Operations,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/planning",
		name: "Planning",
		component: Planning,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/review",
		name: "Review",
		component: Review,
		meta: { requiresAuth: true, showSidebar: true },
	},
	{
		path: "/signin",
		name: "Signin",
		component: Signin,
		meta: { requiresAuth: false, showSidebar: false },
	},
	{
		path: "/sudo",
		name: "Sudo",
		component: Sudo,
		meta: { requiresAuth: true, showSidebar: true },
	},
	// Catch-all route - must be last
  { 
    path: '/:pathMatch(.*)*', 
    name: 'NotFound',
    component: NotFound 
  }
];

const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

// Global navigation guard
router.beforeEach(async (to, from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
    
    if (requiresAuth) {
        try {
            // Check if user is authenticated (could be an API call)
            const isAuthenticated = await apiClient.isAuthenticated();
            if (!isAuthenticated) {
                next('/signin');
            } else {
                next();
            }
        } catch (error) {
						console.log("check auth failed");
            next('/signin');
        }
    } else {
        next();
    }
});


export default router;
