export enum PermissionType {
	DENIED = "denied",
	GRANTED = "granted",
	UNSELECTED = "unselected",
	WITH_OWNER = "with-owner",
}

function isPermissionType(value: string): value is PermissionType {
	return Object.values(PermissionType).includes(value as PermissionType);
}
function toPermissionType(
	value: string,
	defaultValue: PermissionType = PermissionType.UNSELECTED,
): PermissionType {
	if (Object.values(PermissionType).includes(value as PermissionType)) {
		return value as PermissionType;
	}
	return defaultValue;
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
		public location?: Location,
	) {}
}
export interface TegolaURLs {
	nidus: string;
	rmo: string;
}
export interface APIProperties {
	environment: string;
	sentry_dsn: string;
	tegola: TegolaURLs;
	version: string;
}
export interface BoundsDTO {
	min: Location;
	max: Location;
}

export class Bounds {
	min: Location;
	max: Location;
	constructor(
		min: Location = new Location(90, 180),
		max: Location = new Location(-90, -180),
	) {
		this.min = min;
		this.max = max;
	}
	addLocation(l: Location) {
		this.min.latitude = Math.min(this.min.latitude, l.latitude);
		this.min.longitude = Math.min(this.min.longitude, l.longitude);
		this.max.latitude = Math.max(this.max.latitude, l.latitude);
		this.max.longitude = Math.max(this.max.longitude, l.longitude);
	}
}
export interface ContactOptions {
	can_sms: boolean;
	email?: string;
	has_email: boolean;
	has_phone: boolean;
	name?: string;
	phone?: string;
}
export class Contact {
	can_sms: boolean;
	email: string;
	has_email: boolean;
	has_phone: boolean;
	name: string;
	phone: string;
	constructor(options?: ContactOptions) {
		this.can_sms = options?.can_sms ?? false;
		this.email = options?.email ?? "";
		this.has_email = options?.has_email ?? false;
		this.has_phone = options?.has_phone ?? false;
		this.name = options?.name ?? "";
		this.phone = options?.phone ?? "";
	}
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
export interface ComplianceReportRequest {
	id: number;
	public_id: string;
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
export interface ComplianceUpdate {
	access_instructions?: string;
	address?: Address;
	availability_notes?: string;
	comments?: string;
	gate_code?: string;
	has_dog?: boolean;
	//images?: Image[];
	location?: Location;
	permission_type?: string;
	reporter?: Contact;
	//uri: string;
	wants_scheduled?: boolean;
}
export interface Concern {
	type: string;
	url: string;
}
export interface PublicReportDTO {
	address: Address;
	created: string;
	district: string;
	images: Image[];
	location: Location;
	log: LogEntryDTO[];
	public_id: string;
	reporter: Contact;
	status: string;
	type: string;
	uri: string;
}
export interface PublicReportUpdate {
	address?: Address;
	created?: string;
	district?: string;
	images?: Image[];
	location?: Location;
	public_id?: string;
	reporter?: Contact;
	status?: string;
	type?: string;
	uri?: string;
}
export interface PublicReportComplianceCreateRequest {
	client_id: string;
	district?: string;
	mailer_id?: string;
}
export interface PublicReportOptions {
	address: Address;
	created: Date;
	district: string;
	images: Image[];
	location: Location;
	log: LogEntry[];
	public_id: string;
	reporter: Contact;
	status: string;
	type: string;
	uri: string;
}
export class PublicReport {
	address: Address;
	created: Date;
	district: string;
	images: Image[];
	log: LogEntry[];
	public_id: string;
	reporter: Contact;
	status: string;
	type: string;
	uri: string;
	location?: Location;
	constructor(options?: PublicReportOptions) {
		this.address = options?.address ?? new Address();
		this.created = options?.created ?? new Date();
		this.district = options?.district ?? "";
		this.images = options?.images ?? [];
		this.log = options?.log ?? [];
		this.public_id = options?.public_id ?? "";
		this.reporter = options?.reporter ?? new Contact();
		this.status = options?.status ?? "";
		this.type = options?.type ?? "";
		this.uri = options?.uri ?? "";
		this.location = options?.location ?? new Location();
	}
	static fromJSON(json: PublicReportDTO): PublicReport {
		switch (json.type) {
			case "compliance":
				return PublicReportCompliance.fromJSON(
					json as PublicReportComplianceDTO,
				);
			case "nuisance":
				return PublicReportNuisance.fromJSON(json as PublicReportNuisanceDTO);
			case "water":
				return PublicReportWater.fromJSON(json as PublicReportWaterDTO);
			default:
				throw new Error(`Unrecognized public report type '${json.type}'`);
		}
	}
}
export interface PublicReportComplianceDTO extends PublicReportDTO {
	access_instructions: string;
	availability_notes: string;
	comments: string;
	concerns: Concern[];
	gate_code: string;
	has_dog: boolean;
	permission_type: PermissionType;
	wants_scheduled: boolean;
}
export interface PublicReportComplianceOptions extends PublicReportOptions {
	access_instructions: string;
	availability_notes: string;
	comments: string;
	concerns: Concern[];
	gate_code: string;
	has_dog: boolean;
	permission_type: PermissionType;
	wants_scheduled: boolean;
}
export interface PublicReportComplianceUpdate extends PublicReportUpdate {
	access_instructions?: string;
	availability_notes?: string;
	comments?: string;
	gate_code?: string;
	has_dog?: boolean;
	permission_type?: PermissionType;
	wants_scheduled?: boolean;
}
export class PublicReportCompliance extends PublicReport {
	access_instructions: string;
	availability_notes: string;
	comments: string;
	concerns: Concern[];
	gate_code: string;
	has_dog: boolean;
	permission_type: PermissionType;
	wants_scheduled: boolean;
	constructor(options?: PublicReportComplianceOptions) {
		super(options);
		this.access_instructions = options?.access_instructions ?? "";
		this.availability_notes = options?.availability_notes ?? "";
		this.comments = options?.comments ?? "";
		this.concerns = options?.concerns ?? [];
		this.gate_code = options?.gate_code ?? "";
		this.has_dog = options?.has_dog ?? false;
		this.permission_type = toPermissionType(
			options?.permission_type ?? PermissionType.UNSELECTED,
		);
		this.wants_scheduled = options?.wants_scheduled ?? false;
	}
	static fromJSON(json: PublicReportComplianceDTO): PublicReportCompliance {
		return new PublicReportCompliance({
			...json,
			created: new Date(json.created),
			log: json.log.map((l: LogEntryDTO) => LogEntry.fromJSON(l)),
		});
	}
}
export interface PublicReportNuisanceDTO extends PublicReportDTO {
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
export interface PublicReportNuisanceOptions extends PublicReportOptions {
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
export class PublicReportNuisance extends PublicReport {
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
	constructor(options: PublicReportNuisanceOptions) {
		super(options);
		this.additional_info = options.additional_info;
		this.duration = options.duration;
		this.is_location_backyard = options.is_location_backyard;
		this.is_location_frontyard = options.is_location_frontyard;
		this.is_location_garden = options.is_location_garden;
		this.is_location_other = options.is_location_other;
		this.is_location_pool = options.is_location_pool;
		this.source_container = options.source_container;
		this.source_description = options.source_description;
		this.source_gutter = options.source_gutter;
		this.source_stagnant = options.source_stagnant;
		this.time_of_day_day = options.time_of_day_day;
		this.time_of_day_early = options.time_of_day_early;
		this.time_of_day_evening = options.time_of_day_evening;
		this.time_of_day_night = options.time_of_day_night;
	}
	static fromJSON(json: PublicReportNuisanceDTO): PublicReportNuisance {
		return new PublicReportNuisance({
			...json,
			created: new Date(json.created),
			log: json.log.map((l: LogEntryDTO) => LogEntry.fromJSON(l)),
		});
	}
}

export interface PublicReportWaterDTO extends PublicReportDTO {
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
export interface PublicReportWaterOptions extends PublicReportOptions {
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
export class PublicReportWater extends PublicReport {
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
	constructor(options: PublicReportWaterOptions) {
		super(options);
		this.access_comments = options.access_comments;
		this.access_gate = options.access_gate;
		this.access_fence = options.access_fence;
		this.access_locked = options.access_locked;
		this.access_dog = options.access_dog;
		this.access_other = options.access_other;
		this.comments = options.comments;
		this.has_adult = options.has_adult;
		this.has_backyard_permission = options.has_backyard_permission;
		this.has_larvae = options.has_larvae;
		this.has_pupae = options.has_pupae;
		this.is_reporter_confidential = options.is_reporter_confidential;
		this.is_reporter_owner = options.is_reporter_owner;
		this.owner = options.owner;
	}
	static fromJSON(json: PublicReportWaterDTO): PublicReportWater {
		return new PublicReportWater({
			...json,
			created: new Date(json.created),
			log: json.log.map((l: LogEntryDTO) => LogEntry.fromJSON(l)),
		});
	}
}
/*
	address: new Address(),
	comments: "",
	contact: {
		name: "",
		phone: "",
		can_sms: true,
		email: "",
	},
	id: "",
	images: [],
	location: {
		latitude: 0,
		longitude: 0,
	},
	permission: {
		access: PermissionType.UNSELECTED,
		access_instructions: "",
		availability_notes: "",
		gate_code: "",
		has_dog: false,
		wants_scheduled: false,
	},
	uri: "",
});
*/
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
export interface Parcel {
	apn: string;
	description: string;
	id: number;
}
export interface Feature {
	id: number;
	location: Location;
	type: string;
}
export interface Lead {
	compliance_report_requests: ComplianceReportRequest[];
	id: number;
	site_id: number;
	type: string;
}
export interface Site {
	address: Address;
	created: string;
	creator_id: number;
	features: Feature[];
	file_id: number;
	id: number;
	leads: Lead[];
	notes: string;
	organization_id: number;
	owner?: Contact;
	parcel: Parcel;
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
export interface ReviewTaskListResponse {
	tasks: ReviewTask[];
	total: number;
}
export interface SiteListResponse {
	sites: Site[];
	total: number;
}
export interface UploadDTO {
	created: string;
	error: string;
	filename: string;
	id: number;
	recordcount: number;
	status: string;
	type: string;
	csv_pool?: CSVPoolDetail;
}
export interface UploadOptions {
	created: Date;
	error: string;
	filename: string;
	id: number;
	recordcount: number;
	status: string;
	type: string;
	csv_pool?: CSVPoolDetail;
}
export class Upload {
	created: Date;
	error: string;
	filename: string;
	id: number;
	recordcount: number;
	status: string;
	type: string;
	csv_pool?: CSVPoolDetail;
	constructor(options: UploadOptions) {
		this.created = options.created;
		this.error = options.error;
		this.filename = options.filename;
		this.id = options.id;
		this.recordcount = options.recordcount;
		this.status = options.status;
		this.type = options.type;
		this.csv_pool = options.csv_pool;
	}
	static fromJSON(json: UploadDTO): Upload {
		return new Upload({
			...json,
			created: new Date(json.created),
		});
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
export type MailerStatus = "created" | "printed" | "mailed" | "completed";
export interface MailerDTO {
	address: Address;
	compliance_report_request_id?: string;
	created: string;
	id: string;
	recipient: string;
	status: MailerStatus;
	site_id: string;
	uri: string;
}
export interface MailerOptions {
	address: Address;
	compliance_report_request_id?: string;
	created: Date;
	id: string;
	recipient: string;
	site_id: string;
	status: MailerStatus;
	uri: string;
}
export class Mailer {
	address: Address;
	compliance_report_request_id?: string;
	created: Date;
	id: string;
	recipient: string;
	site_id: string;
	status: MailerStatus;
	uri: string;
	constructor(options: MailerOptions) {
		this.address = options.address;
		this.compliance_report_request_id = options.compliance_report_request_id;
		this.created = options.created;
		this.id = options.id;
		this.recipient = options.recipient;
		this.site_id = options.site_id;
		this.status = options.status;
		this.uri = options.uri;
	}
	pdfUrl(): string {
		return `/mailer/mode-3/${this.compliance_report_request_id}/preview`;
	}
	static fromJSON(json: MailerDTO): Mailer {
		return new Mailer({
			...json,
			created: new Date(json.created),
		});
	}
}

export interface Organization {
	id: number;
	name: string;
	lob_address_id: string;
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
export interface ServiceRequestDTO {
	address: string;
	assigned_technician: string;
	city: string;
	created: string;
	h3cell: number;
	has_dog: boolean;
	has_spanish_speaker: boolean;
	id: string;
	priority: string;
	recorded_date: string;
	source: string;
	status: string;
	target: string;
	zip: string;
}
export interface ServiceRequestOptions {
	address: string;
	assigned_technician: string;
	city: string;
	created: Date;
	h3cell: number;
	has_dog: boolean;
	has_spanish_speaker: boolean;
	id: string;
	priority: string;
	recorded_date: Date;
	source: string;
	status: string;
	target: string;
	zip: string;
}
export class ServiceRequest {
	address: string;
	assigned_technician: string;
	city: string;
	created: Date;
	h3cell: number;
	has_dog: boolean;
	has_spanish_speaker: boolean;
	id: string;
	priority: string;
	recorded_date: Date;
	source: string;
	status: string;
	target: string;
	zip: string;
	constructor(options: ServiceRequestOptions) {
		this.address = options.address;
		this.assigned_technician = options.assigned_technician;
		this.city = options.city;
		this.created = options.created;
		this.h3cell = options.h3cell;
		this.has_dog = options.has_dog;
		this.has_spanish_speaker = options.has_spanish_speaker;
		this.id = options.id;
		this.priority = options.priority;
		this.recorded_date = options.recorded_date;
		this.source = options.source;
		this.status = options.status;
		this.target = options.target;
		this.zip = options.zip;
	}
	static fromJSON(json: ServiceRequestDTO): ServiceRequest {
		return new ServiceRequest({
			...json,
			created: new Date(json.created),
			recorded_date: new Date(json.recorded_date),
		});
	}
}
export interface SyncDTO {
	created: Date;
	id: string;
	organization: string;
	records_created: number;
	records_unchanged: number;
	records_updated: number;
}
export class Sync {
	constructor(
		public created: Date,
		public id: string,
		public organization: string,
		public records_created: number,
		public records_unchanged: number,
		public records_updated: number,
	) {}
	static fromJSON(json: SyncDTO): Sync {
		return new Sync(
			new Date(json.created),
			json.id,
			json.organization,
			json.records_created,
			json.records_unchanged,
			json.records_updated,
		);
	}
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
	mailer: string;
	publicreport_message: string;
	review_task: string;
	service_request: string;
	signal: string;
	site: string;
	sync: string;
	upload: string;
	user: string;
}
