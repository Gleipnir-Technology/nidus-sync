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
export interface District {
	name: string;
	phone_office: string;
	slug: string;
	uri: string;
	url_logo: string;
	url_website: string;
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
export interface PublicreportDTO {
	address: string;
	created: string;
	district: string;
	id: string;
	image_count: number;
	location: Location;
	status: string;
	type: string;
	uri: string;
}
export class Publicreport {
	constructor(
		public address: string,
		public created: Date,
		public district: string,
		public id: string,
		public image_count: number,
		public location: Location,
		public status: string,
		public type: string,
		public uri: string,
	) {}
	static fromJSON(json: PublicreportDTO): Publicreport {
		return new Publicreport(
			json.address,
			new Date(json.created),
			json.district,
			json.id,
			json.image_count,
			json.location,
			json.status,
			json.type,
			json.uri,
		);
	}
}
