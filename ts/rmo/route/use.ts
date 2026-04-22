import { RouteLocationRaw } from "vue-router";
import { ROUTE_NAMES } from "@/rmo/route/name";

export function useRoutes() {
	/*
	const RMOComplianceAddress = (publicID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.COMPLIANCE_ADDRESS,
			...(publicID && { query: { publicID: publicID } })
		}
	}
   */
	const ComplianceAddress = (publicID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.COMPLIANCE_ADDRESS,
			params: {
				public_id: publicID,
			},
		};
	};
	return {
		ComplianceAddress,
	};
}
