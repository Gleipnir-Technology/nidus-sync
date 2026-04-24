import { RouteLocationRaw } from "vue-router";
import { ROUTE_NAMES } from "@/route/name";

export function useRoutes() {
	const CellDetail = (cell: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.CELL_DETAIL,
			params: {
				cell: cell,
			},
		};
	};
	const ReviewSite = (siteID: string): RouteLocationRaw => {
		return {
			name: ROUTE_NAMES.REVIEW_SITE,
			query: { site: siteID },
		};
	};
	return {
		CellDetail,
		ReviewSite,
	};
}
