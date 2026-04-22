import { RouteLocationRaw } from "vue-router";
import { ROUTE_NAMES } from "@/route/name";

export function useRoutes() {
	/*
	const RMOComplianceAddress = (publicID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.COMPLIANCE_ADDRESS,
			...(publicID && { query: { publicID: publicID } })
		}
	}
   */
	const ReviewSite = (siteID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.REVIEW_SITE,
			query: { site: siteID },
		};
	};
	return {
		ReviewSite,
	};
}
