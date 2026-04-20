import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Compliance from "@/rmo/view/Compliance.vue";
import ComplianceAddress from "@/rmo/content/compliance/Address.vue";
import ComplianceComplete from "@/rmo/content/compliance/Complete.vue";
import ComplianceConcern from "@/rmo/content/compliance/Concern.vue";
import ComplianceContact from "@/rmo/content/compliance/Contact.vue";
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
				name: "ComplianceIntro",
				path: "",
			},
			{
				component: ComplianceAddress,
				name: "ComplianceAddress",
				path: "address",
			},
			{
				component: ComplianceComplete,
				name: "ComplianceComplete",
				path: "complete",
			},
			{
				component: ComplianceConcern,
				name: "ComplianceConcern",
				path: "concern",
			},
			{
				component: ComplianceContact,
				name: "ComplianceContact",
				path: "contact",
			},
			{
				component: ComplianceEvidence,
				name: "ComplianceEvidence",
				path: "evidence",
			},
			{
				component: CompliancePermission,
				name: "CompliancePermission",
				path: "permission",
			},
			{
				component: ComplianceProcess,
				name: "ComplianceProcess",
				path: "process",
			},
			{
				component: ComplianceSubmit,
				name: "ComplianceSubmit",
				path: "submit",
			},
		],
		component: Compliance,
		path: "/district/:slug/compliance",
		name: "Compliance",
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
		children: [
			{
				component: ComplianceIntro,
				name: "ComplianceIntro",
				path: "",
			},
			{
				component: ComplianceAddress,
				name: "ComplianceAddress",
				path: "address",
			},
			{
				component: ComplianceComplete,
				name: "ComplianceComplete",
				path: "complete",
			},
			{
				component: ComplianceConcern,
				name: "ComplianceConcern",
				path: "concern",
			},
			{
				component: ComplianceContact,
				name: "ComplianceContact",
				path: "contact",
			},
			{
				component: ComplianceEvidence,
				name: "ComplianceEvidence",
				path: "evidence",
			},
			{
				component: CompliancePermission,
				name: "CompliancePermission",
				path: "permission",
			},
			{
				component: ComplianceProcess,
				name: "ComplianceProcess",
				path: "process",
			},
			{
				component: ComplianceSubmit,
				name: "ComplianceSubmit",
				path: "submit",
			},
		],
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
