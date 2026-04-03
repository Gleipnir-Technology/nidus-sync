import type { Map as MapLibreMap } from "maplibre-gl";

export interface Address {
	country: string;
	locality: string;
	number: string;
	postal_code: string;
	raw: string;
	region: string;
	street: string;
	unit: string;
}
export interface Bounds {
	min: Location;
	max: Location;
}
export interface Changes {
	updated: string[];
	unchanged: string[];
}

export interface Communication {
	created: string;
	id: string;
	public_report: PublicReport | null;
	type: string;
}
export interface Contact {
	has_email: boolean;
	has_phone: boolean;
	name?: string;
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
export interface CSVPoolDetail {
	count: CSVPoolDetailCount;
	errors: CSVPoolError[];
	pools: UploadPoolRow[];
}
export interface Exif {
	created: string;
	make: string;
	model: string;
}
export interface Followup {
	description: string;
	id: number;
	title: string;
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
export interface Lead {
	description: string;
	id: number;
	title: string;
}
export interface Location {
	lat: number;
	lng: number;
}
export interface LogEntry {
	created: string;
	id: number;
	message: string;
	report_id: number;
	type: string;
	user_id: number;
}
export interface MapClickEvent {
	location: Location;
	map: MapLibreMap;
	point: Point;
}
export interface Marker {
	color?: string;
	draggable?: boolean;
	id: string;
	location: Location;
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
export interface Organization {
	id: number;
	service_area?: Bounds;
}
export interface Point {
	x: number;
	y: number;
}
export interface Pool {
	condition: string;
	id: number;
	location: Location;
	site: Site;
}
export interface PublicReport {
	address: Address;
	address_raw: string;
	created: string;
	images: Image[];
	location?: Location;
	log: LogEntry[];
	nuisance?: Nuisance;
	public_id: string;
	reporter: Contact;
	status: string;
	type: string;
	water?: Water;
}
export interface Signal {
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
export interface ReviewTask {
	address: Address;
	addressed?: string;
	addressor?: User;
	created: string;
	creator: User;
	pool?: ReviewTaskPool;
	id: number;
}
export interface ReviewTaskPool {
	condition: string;
	location: Location;
	owner: Contact;
	site: Site;
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
export interface Upload {
	created: string;
	filename: string;
	id: number;
	recordcount: number;
	status: string;
	type: string;
	csv_pool?: CSVPoolDetail;
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
export interface UserNotificationCounts {
	communication: number;
	home: number;
	review: number;
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
