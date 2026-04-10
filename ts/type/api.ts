export enum PermissionAccess {
	DENIED = "denied",
	GRANTED = "granted",
	UNSELECTED = "unselected",
	WITH_OWNER = "with-owner",
}
export class Address {
	constructor(
		public country: string = "",
		public gid: string = "",
		public locality: string = "",
		public number: string = "",
		public postal_code: string = "",
		public raw: string = "",
		public region: string = "",
		public street: string = "",
		public unit: string = "",
	) {}
}
export interface Bounds {
	min: Location;
	max: Location;
}
export interface Contact {
	has_email: boolean;
	has_phone: boolean;
	name?: string;
}
export interface District {
	name: string;
	phone_office: string;
	slug: string;
	uri: string;
	url_logo: string;
	url_website: string;
}
export class Location {
	accuracy?: number;
	latitude: number;
	longitude: number;
	constructor(latitude: number = 0, longitude: number = 0, accuracy?: number) {
		this.accuracy = accuracy;
		this.latitude = latitude;
		this.longitude = longitude;
	}
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
export interface CSVPoolDetailCount {
	existing: number;
	new: number;
	outside: number;
}
export interface CSVPoolError {
	column: number;
	line: number;
	message: string;
}
export interface Followup {
	description: string;
	id: number;
	title: string;
}
export interface Lead {
	description: string;
	id: number;
	title: string;
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
export interface Exif {
	created: string;
	make: string;
	model: string;
}
export interface Image {
	distance_from_reporter_meters?: number;
	exif: Exif;
	exif_make: string;
	exif_model: string;
	exif_datetime: string;
	location?: Location;
	report_id: number;
	url_content: string;
	uuid: string;
}
export interface Nuisance {
	additional_info: string;
	duration: string;
	is_location_backyard: boolean;
	is_location_frontyard: boolean;
	is_location_garden: boolean;
	is_location_other: boolean;
	is_location_pool: boolean;
	source_container: boolean;
	source_description: string;
	source_gutter: boolean;
	source_stagnant: boolean;
	time_of_day_day: boolean;
	time_of_day_early: boolean;
	time_of_day_evening: boolean;
	time_of_day_night: boolean;
}
export interface Water {
	access_comments: string;
	access_gate: boolean;
	access_fence: boolean;
	access_locked: boolean;
	access_dog: boolean;
	access_other: boolean;
	comments: string;
	has_adult: boolean;
	has_backyard_permission: boolean;
	has_larvae: boolean;
	has_pupae: boolean;
	is_reporter_confidential: boolean;
	is_reporter_owner: boolean;
	owner: Contact;
}
export interface PublicReportDTO {
	address: Address;
	created: string;
	district: string;
	id: string;
	images: Image[];
	location: Location;
	log: LogEntryDTO[];
	nuisance: Nuisance;
	reporter: Contact;
	status: string;
	type: string;
	water: Water;
	uri: string;
}
export class PublicReport {
	constructor(
		public address: Address,
		public created: Date,
		public district: string,
		public id: string,
		public images: Image[],
		public log: LogEntry[],
		public reporter: Contact,
		public status: string,
		public type: string,
		public uri: string,
		public location?: Location,
		public nuisance?: Nuisance,
		public water?: Water,
	) {}
	static fromJSON(json: PublicReportDTO): PublicReport {
		return new PublicReport(
			json.address,
			new Date(json.created),
			json.district,
			json.id,
			json.images,
			json.log.map((l: LogEntryDTO) => LogEntry.fromJSON(l)),
			json.reporter,
			json.status,
			json.type,
			json.uri,
			json.location,
			json.nuisance,
			json.water,
		);
	}
}
export interface CommunicationDTO {
	created: string;
	id: string;
	public_report?: PublicReportDTO;
	type: string;
}
export class Communication {
	constructor(
		public created: Date,
		public id: string,
		public type: string,
		public public_report?: PublicReport,
	) {}
	static fromJSON(json: CommunicationDTO): Communication {
		return new Communication(
			new Date(json.created),
			json.id,
			json.type,
			json.public_report == undefined
				? undefined
				: PublicReport.fromJSON(json.public_report),
		);
	}
}

export interface Pool {
	condition: string;
	id: number;
	location: Location;
	site: Site;
}
export interface SignalDTO {
	address?: Address;
	addressed?: string;
	addressor?: number;
	created: string;
	creator: number;
	id: number;
	location: Location;
	pool?: Pool;
	report?: PublicReport;
	species?: string;
	type: string;
}
export class Signal {
	constructor(
		public created: Date,
		public creator: number,
		public id: number,
		public location: Location,
		public type: string,
		public address?: Address,
		public addressed?: string,
		public addressor?: number,
		public pool?: Pool,
		public report?: PublicReport,
		public species?: string,
	) {}
	static fromJSON(json: SignalDTO): Signal {
		return new Signal(
			new Date(json.created),
			json.creator,
			json.id,
			json.location,
			json.type,
			json.address,
			json.addressed,
			json.addressor,
			json.pool,
			json.report,
			json.species,
		);
	}
}
export interface Site {
	address: Address;
	created: string;
	creator_id: number;
	file_id: number;
	id: number;
	location: Location;
	notes: string;
	organization_id: number;
	owner?: Contact;
	parcel_id?: number;
	resident?: Contact;
	resident_owned: boolean;
	tags: Map<string, string>;
	version: number;
}
export interface ReviewTaskPool {
	condition: string;
	location: Location;
	owner: Contact;
	site: Site;
}
export interface ReviewTask {
	address: Address;
	addressed?: string;
	addressor?: User;
	created: string;
	creator: User;
	pool?: ReviewTaskPool;
	id: number;
}
export interface UploadDTO {
	created: string;
	filename: string;
	id: number;
	recordcount: number;
	status: string;
	type: string;
	csv_pool?: CSVPoolDetail;
}
export class Upload {
	constructor(
		public created: Date,
		public filename: string,
		public id: number,
		public recordcount: number,
		public status: string,
		public type: string,
		public csv_pool?: CSVPoolDetail,
	) {}
	static fromJSON(json: UploadDTO): Upload {
		return new Upload(
			new Date(json.created),
			json.filename,
			json.id,
			json.recordcount,
			json.status,
			json.type,
			json.csv_pool,
		);
	}
}

export interface UploadPoolRow {
	address: Address;
	condition: string;
	errors: UploadPoolError[];
	status: string;
	tags: Map<string, string>;
}
export interface UploadPoolError {
	column: number;
	line: number;
	message: string;
}

export interface CSVPoolDetail {
	count: CSVPoolDetailCount;
	errors: CSVPoolError[];
	pools: UploadPoolRow[];
}
export interface User {
	avatar: string;
	display_name: string;
	id: number;
	initials: string;
	is_active: boolean;
	role: string;
	tags: string[];
	uri: string;
	username: string;
}
export interface Organization {
	id: number;
	service_area?: Bounds;
}
export interface UserNotificationCounts {
	communication: number;
	home: number;
	review: number;
}
export interface SessionNotificationCounts {
	communication: number;
	home: number;
	review: number;
}
export interface Session {
	impersonating?: string;
	notifications: Notification[];
	notification_counts: SessionNotificationCounts;
	organization: Organization;
	self: User;
	urls: URLs;
}
export interface URLs {
	api: URLsAPI;
	tegola: string;
	tile: string;
}
// Define interfaces matching your Go structs
interface URLsAPI {
	avatar: string;
	communication: string;
	impersonation: string;
	publicreport_message: string;
	review_task: string;
	signal: string;
	upload: string;
	user: string;
}
