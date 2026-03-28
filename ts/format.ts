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

export function shortAddress(a: Address): string {
	if (!a) return "";
	return `${a.number} ${a.street}, ${a.locality}`;
}
