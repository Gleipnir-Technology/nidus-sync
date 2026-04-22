import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Compliance from "@/rmo/view/Compliance.vue";
import ComplianceAddress from "@/rmo/content/compliance/Address.vue";
import ComplianceComplete from "@/rmo/content/compliance/Complete.vue";
import ComplianceConcern from "@/rmo/content/compliance/Concern.vue";
import ComplianceContact from "@/rmo/content/compliance/Contact.vue";
import ComplianceDistrict from "@/rmo/view/ComplianceDistrict.vue";
import ComplianceEvidence from "@/rmo/content/compliance/Evidence.vue";
import ComplianceIntro from "@/rmo/content/compliance/Intro.vue";
import ComplianceMailer from "@/rmo/view/ComplianceMailer.vue";
import CompliancePermission from "@/rmo/content/compliance/Permission.vue";
import ComplianceProcess from "@/rmo/content/compliance/Process.vue";
import ComplianceSubmit from "@/rmo/content/compliance/Submit.vue";
import HomeBase from "@/rmo/view/Home.vue";
import HomeDistrict from "@/rmo/view/district/Home.vue";
import NuisanceBase from "@/rmo/view/Nuisance.vue";
import NuisanceDistrict from "@/rmo/view/district/Nuisance.vue";
import ReportSubmitted from "@/rmo/view/ReportSubmitted.vue";
import StatusBase from "@/rmo/view/Status.vue";
import StatusByID from "@/rmo/view/StatusByID.vue";
import StatusDistrict from "@/rmo/view/district/Status.vue";
import Water from "@/rmo/view/Water.vue";
import WaterDistrict from "@/rmo/view/district/Water.vue";

import { ROUTE_NAMES } from "@/rmo/route/name";

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
		children: [
			{
				component: ComplianceIntro,
				name: ROUTE_NAMES.COMPLIANCE_INTRO,
				path: "",
			},
			{
				component: ComplianceAddress,
				name: ROUTE_NAMES.COMPLIANCE_ADDRESS,
				path: "address",
			},
			{
				component: ComplianceComplete,
				name: ROUTE_NAMES.COMPLIANCE_COMPLETE,
				path: "complete",
			},
			{
				component: ComplianceConcern,
				name: ROUTE_NAMES.COMPLIANCE_CONCERN,
				path: "concern",
			},
			{
				component: ComplianceContact,
				name: ROUTE_NAMES.COMPLIANCE_CONTACT,
				path: "contact",
			},
			{
				component: ComplianceEvidence,
				name: ROUTE_NAMES.COMPLIANCE_EVIDENCE,
				path: "evidence",
			},
			{
				component: CompliancePermission,
				name: ROUTE_NAMES.COMPLIANCE_PERMISSION,
				path: "permission",
			},
			{
				component: ComplianceProcess,
				name: ROUTE_NAMES.COMPLIANCE_PROCESS,
				path: "process",
			},
			{
				component: ComplianceSubmit,
				name: ROUTE_NAMES.COMPLIANCE_SUBMIT,
				path: "submit",
			},
		],
		component: Compliance,
		path: "/compliance/:public_id",
		name: "Compliance",
		props: true,
	},
	{
		component: ComplianceDistrict,
		path: "/district/:slug/compliance",
		name: "ComplianceDistrict",
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
		path: "/mailer/:public_id",
		name: "ComplianceMailer",
		component: ComplianceMailer,
		props: true,
	},
	{
		path: "/submitted/:id",
		name: "ReportSubmitted",
		component: ReportSubmitted,
		props: true,
	},
	{
		path: "/status",
		name: "StatusBase",
		component: StatusBase,
	},
	{
		component: StatusByID,
		name: "StatusbyID",
		path: "/status/:id",
		props: true,
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
