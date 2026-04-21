import { Address } from "@/type/api";

export function formatAddress(address?: Address): string {
	if (!address) {
		return "undefined";
	}
	if (address.number === "" && address.street === "") {
		return "no address provided";
	}
	return `${address.number.trim()} ${address.street.trim()}, ${address.locality}`;
}
export function formatBigNumber(n: number): string {
	// Convert the number to a string
	const numStr = n.toString();

	// Add commas every three digits from the right
	let result = "";
	for (let i = 0; i < numStr.length; i++) {
		if (i > 0 && (numStr.length - i) % 3 === 0) {
			result += ",";
		}
		result += numStr[i];
	}

	return result;
}
export function formatDate(date: Date): string {
	return new Intl.DateTimeFormat("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	}).format(date);
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
export function formatReportID(id: string): string {
	if (id.length === 12) {
		return `${id.substring(0, 4)}-${id.substring(4, 8)}-${id.substring(8)}`;
	}
	return id;
}
export function formatTimeRelative(t: Date): string {
	const now = new Date();
	const diffMs = now.getTime() - t.getTime();

	const hours = diffMs / (1000 * 60 * 60);

	if (hours > 0) {
		if (hours < 1) {
			const minutes = diffMs / (1000 * 60);
			return `${Math.floor(minutes)} minutes ago`;
		} else if (hours < 24) {
			return `${Math.floor(hours)} hours ago`;
		} else {
			const days = hours / 24;
			return `${Math.floor(days)} days ago`;
		}
	} else {
		if (hours < -24) {
			const days = hours / 24;
			return `in ${Math.floor(-1 * days)} days`;
		} else if (hours < -1) {
			return `in ${Math.floor(-1 * hours)} hours`;
		} else {
			const minutes = diffMs / (1000 * 60);
			if (minutes > -1) {
				const seconds = diffMs / 1000;
				return `in ${Math.floor(-1 * seconds)} seconds`;
			}
			return `in ${Math.floor(-1 * minutes)} minutes`;
		}
	}
}
export function shortAddress(a: Address | undefined): string {
	if (!a) return "unknown";
	return `${a.number} ${a.street}, ${a.locality}`;
}
