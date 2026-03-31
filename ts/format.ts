import { Address } from "./types";

export function formatAddress(address?: Address): string {
	if (!address) {
		return "undefined";
	}
	if (address.number === "" && address.street === "") {
		return "no address provided";
	}
	return `${address.number} ${address.street}, ${address.locality}`;
}
export function formatDistance(meters: number | undefined) {
	if (meters === undefined || meters === null) {
		return "unknown";
	}
	if (meters < 1) {
		const mm = Math.round(meters * 1000);
		return `${mm} mm`;
	} else if (meters >= 1000) {
		const km = Math.round(meters / 1000);
		return `${km} km`;
	} else {
		const m = Math.round(meters);
		return `${m} m`;
	}
}
export function formatRelativeTime(dateString: string): string {
	if (!dateString) return "";

	const date = new Date(dateString);
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffMins = Math.floor(diffMs / 60000);
	const diffHours = Math.floor(diffMins / 60);
	const diffDays = Math.floor(diffHours / 24);

	if (diffMins < 1) return "just now";
	if (diffMins < 60) return `${diffMins} min ago`;
	if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? "s" : ""} ago`;
	return `${diffDays} day${diffDays > 1 ? "s" : ""} ago`;
}

export function shortAddress(a: Address | undefined): string {
	if (!a) return "unknown";
	return `${a.number} ${a.street}, ${a.locality}`;
}
