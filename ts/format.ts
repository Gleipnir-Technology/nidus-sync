import { Address } from "./types";

export function shortAddress(a: Address): string {
	if (!a) return "";
	return `${a.number} ${a.street}, ${a.locality}`;
}
