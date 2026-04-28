export const ROUTE_NAMES = {
	COMPLIANCE_ADDRESS: "compliance-address",
	COMPLIANCE_COMPLETE: "compliance-complete",
	COMPLIANCE_CONCERN: "compliance-concern",
	COMPLIANCE_CONTACT: "compliance-contact",
	COMPLIANCE_EVIDENCE: "compliance-evidence",
	COMPLIANCE_INTRO: "compliance-intro",
	COMPLIANCE_PERMISSION: "compliance-permission",
	COMPLIANCE_PROCESS: "compliance-process",
	COMPLIANCE_SUBMIT: "compliance-submit",
	REGISTER_NOTIFICATIONS_COMPLETE: "register-notifications-complete",
	REVIEW_SITE: "review-site",
	STATUS_BY_ID: "status-by-id",
} as const;

export type RouteName = (typeof ROUTE_NAMES)[keyof typeof ROUTE_NAMES];
