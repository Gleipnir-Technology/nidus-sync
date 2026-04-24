export const ROUTE_NAMES = {
	CELL_DETAIL: "cell-detail",
	COMPLIANCE_ADDRESS: "compliance-address",
	REVIEW_SITE: "review-site",
} as const;

export type RouteName = (typeof ROUTE_NAMES)[keyof typeof ROUTE_NAMES];
