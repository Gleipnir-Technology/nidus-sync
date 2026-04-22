import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";

import { useSessionStore } from "@/store/session";
import Home from "@/view/Home.vue";
import Authenticated from "@/view/Authenticated.vue";
import Cell from "@/view/Cell.vue";
import Communication from "@/view/Communication.vue";
import ConfigurationIntegration from "@/view/configuration/Integration.vue";
import ConfigurationIntegrationArcgis from "@/view/configuration/IntegrationArcgis.vue";
import ConfigurationOrganization from "@/view/configuration/Organization.vue";
import ConfigurationPesticide from "@/view/configuration/Pesticide.vue";
import ConfigurationPesticideAdd from "@/view/configuration/PesticideAdd.vue";
import ConfigurationRoot from "@/view/configuration/Root.vue";
import ConfigurationUpload from "@/view/configuration/Upload.vue";
import ConfigurationUploadDetail from "@/view/configuration/UploadDetail.vue";
import ConfigurationUploadPool from "@/view/configuration/UploadPool.vue";
import ConfigurationUploadPoolCustom from "@/view/configuration/UploadPoolCustom.vue";
import ConfigurationUploadPoolFlyover from "@/view/configuration/UploadPoolFlyover.vue";
import ConfigurationUser from "@/view/configuration/User.vue";
import ConfigurationUserAdd from "@/view/configuration/UserAdd.vue";
import ConfigurationUserEdit from "@/view/configuration/UserEdit.vue";
import Dash from "@/view/Dash.vue";
import Intelligence from "@/view/Intelligence.vue";
import NotFound from "@/view/NotFound.vue";
import OAuthRefreshArcgis from "@/view/OAuthRefreshArcgis.vue";
import Operations from "@/view/Operations.vue";
import Planning from "@/view/Planning.vue";
import ReviewMailer from "@/view/review/Mailer.vue";
import ReviewPool from "@/view/review/Pool.vue";
import ReviewRoot from "@/view/review/Root.vue";
import ReviewSite from "@/view/review/Site.vue";
import Signin from "@/view/Signin.vue";
import Signout from "@/view/Signout.vue";
import Signup from "@/view/Signup.vue";
import Sudo from "@/view/Sudo.vue";
import { apiClient } from "@/client";

import { ROUTE_NAMES } from "@/route/name";

const routes: RouteRecordRaw[] = [
	{
		path: "/",
		name: "Home",
		component: Home,
	},
	{
		children: [
			{
				component: Cell,
				name: "Cell",
				path: "/_/cell/:cell",
				props: true,
			},
			{
				path: "/_/communication",
				name: "Communication",
				component: Communication,
			},
			{
				path: "/_/configuration",
				name: "Configuration",
				component: ConfigurationRoot,
			},
			{
				path: "/_/configuration/integration",
				name: "Integration Configuration",
				component: ConfigurationIntegration,
			},
			{
				path: "/_/configuration/integration/arcgis",
				name: "Arcgis Integration Configuration",
				component: ConfigurationIntegrationArcgis,
			},
			{
				path: "/_/configuration/organization",
				name: "Organization Configuration",
				component: ConfigurationOrganization,
			},
			{
				path: "/_/configuration/pesticide",
				name: "Pesticide Configuration",
				component: ConfigurationPesticide,
			},
			{
				path: "/_/configuration/pesticide/add",
				name: "Pesticide Add",
				component: ConfigurationPesticideAdd,
			},
			{
				path: "/_/configuration/upload",
				name: "Upload Configuration",
				component: ConfigurationUpload,
			},
			{
				component: ConfigurationUploadDetail,
				name: "Upload Detail",
				path: "/_/configuration/upload/:id",
				props: (route) => ({
					id: parseInt(route.params.id as string, 10),
				}),
			},
			{
				path: "/_/configuration/upload/pool",
				name: "Pool Upload",
				component: ConfigurationUploadPool,
			},
			{
				path: "/_/configuration/upload/pool/custom",
				name: "Custom Pool Upload",
				component: ConfigurationUploadPoolCustom,
			},
			{
				path: "/_/configuration/upload/pool/flyover",
				name: "Flyover Upload",
				component: ConfigurationUploadPoolFlyover,
			},
			{
				path: "/_/configuration/user",
				name: "User Configuration",
				component: ConfigurationUser,
			},
			{
				path: "/_/configuration/user/add",
				name: "User Add Configuration",
				component: ConfigurationUserAdd,
			},
			{
				component: ConfigurationUserEdit,
				name: "User Edit",
				path: "/_/configuration/user/:id",
				props: (route) => ({
					id: parseInt(route.params.id as string, 10),
				}),
			},
			{
				path: "/_/dash",
				name: "Dash",
				component: Dash,
			},
			{
				path: "/_/intelligence",
				name: "Intelligence",
				component: Intelligence,
			},
			{
				path: "/_/oauth/refresh/arcgis",
				name: "Arcgis OAuth Refresh",
				component: OAuthRefreshArcgis,
			},
			{
				path: "/_/operations",
				name: "Operations",
				component: Operations,
			},
			{
				path: "/_/planning",
				name: "Planning",
				component: Planning,
			},
			{
				path: "/_/review",
				name: "Review",
				component: ReviewRoot,
			},
			{
				path: "/_/review/mailer",
				name: "Mailer Review",
				component: ReviewMailer,
			},
			{
				path: "/_/review/pool",
				name: "Pool Review",
				component: ReviewPool,
			},
			{
				path: "/_/review/site",
				name: ROUTE_NAMES.REVIEW_SITE,
				component: ReviewSite,
			},
			{
				path: "/_/sudo",
				name: "Sudo",
				component: Sudo,
			},
		],
		component: Authenticated,
		path: "/_",
		name: "Authenticated",
	},
	{
		component: Signin,
		name: "Signin",
		path: "/signin",
	},
	{
		component: Signout,
		name: "Signout",
		path: "/signout",
	},
	{
		component: Signup,
		name: "Signup",
		path: "/signup",
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
	if (to.fullPath.startsWith("/signin") || to.fullPath == "/signup") {
		return;
	}
	const storeSession = useSessionStore();
	try {
		if (!storeSession.isLoading && !storeSession.isAuthenticated) {
			console.log("sending to signin because we're not authenticated");
			return `/signin?next=${from.fullPath}`;
		}
	} catch (error) {
		console.log("check auth failed");
	}
	return;
});

export default router;
