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
export interface LogEntryDTO {
	created: string;
	message: string;
	type: string;
	user_id: number;
}
export class LogEntry {
	constructor(
		public created: Date,
		public message: string,
		public type: string,
		public user_id: number,
	) {}
	static fromJSON(json: LogEntryDTO): LogEntry {
		return new LogEntry(
			new Date(json.created),
			json.message,
			json.type,
			json.user_id,
		);
	}
}
export interface PublicreportDTO {
	address: Address;
	created: string;
	district: string;
	id: string;
	image_count: number;
	location: Location;
	log: LogEntryDTO[];
	status: string;
	type: string;
	uri: string;
}
export class Publicreport {
	constructor(
		public address: Address,
		public created: Date,
		public district: string,
		public id: string,
		public image_count: number,
		public location: Location,
		public log: LogEntry[],
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
			json.log.map((l: LogEntryDTO) => LogEntry.fromJSON(l)),
			json.status,
			json.type,
			json.uri,
		);
	}
}
