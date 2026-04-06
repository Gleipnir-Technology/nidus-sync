export interface Address {
	country: string;
	gid: string;
	locality: string;
	number: string;
	postal_code: string;
	raw: string;
	region: string;
	street: string;
	unit: string;
}
export interface Location {
	latitude: number;
	longitude: number;
}
export interface GeocodeSuggestion {
	detail: string;
	gid: string;
	locality: string;
	type: string;
}
export interface Geocode {
	address: Address;
	cell: number;
	location: Location;
}
