import { RouteLocationRaw } from "vue-router";
import { ROUTE_NAMES } from "@/rmo/route/name";

function complianceRoute(name: string) {
	return (publicID: string): RouteLocationRaw => {
		return {
			name: name,
			params: {
				public_id: publicID,
			},
		};
	};
}
export function useRoutes() {
	/*
	const RMOComplianceAddress = (publicID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.COMPLIANCE_ADDRESS,
			...(publicID && { query: { publicID: publicID } })
		}
	}
   */
	const ComplianceAddress = complianceRoute(ROUTE_NAMES.COMPLIANCE_ADDRESS);
	const ComplianceComplete = complianceRoute(ROUTE_NAMES.COMPLIANCE_COMPLETE);
	const ComplianceConcern = complianceRoute(ROUTE_NAMES.COMPLIANCE_CONCERN);
	const ComplianceContact = complianceRoute(ROUTE_NAMES.COMPLIANCE_CONTACT);
	const ComplianceEvidence = complianceRoute(ROUTE_NAMES.COMPLIANCE_EVIDENCE);
	const ComplianceIntro = complianceRoute(ROUTE_NAMES.COMPLIANCE_INTRO);
	const CompliancePermission = complianceRoute(
		ROUTE_NAMES.COMPLIANCE_PERMISSION,
	);
	const ComplianceProcess = complianceRoute(ROUTE_NAMES.COMPLIANCE_PROCESS);
	const ComplianceSubmit = complianceRoute(ROUTE_NAMES.COMPLIANCE_SUBMIT);
	const RegisterNotificationsComplete = (
		publicID: string,
	): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.REGISTER_NOTIFICATIONS_COMPLETE,
			params: {
				public_id: publicID,
			},
		};
	};
	return {
		ComplianceAddress,
		ComplianceComplete,
		ComplianceConcern,
		ComplianceContact,
		ComplianceEvidence,
		ComplianceIntro,
		CompliancePermission,
		ComplianceProcess,
		ComplianceSubmit,
		RegisterNotificationsComplete,
	};
}
