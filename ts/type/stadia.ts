// Interface definitions
interface AddressComponents {
	number?: string;
	postal_code?: string;
	street?: string;
}

interface AddressContext {
	iso_3166_a2: string; // "US"
	iso_3166_a3: string; // "USA"
	whosonfirst?: AddressContextWhosOnFirst;
}
interface AddressContextWhosOnFirst {
	country: WhosOnFirstEntry;
	county: WhosOnFirstEntry;
	locality: WhosOnFirstEntry;
	region: WhosOnFirstEntry;
}
interface WhosOnFirstEntry {
	abbreviation?: string; // "SL" or "UT"
	gid: string; // "whosonfirst:county:102082877"
	name: string; // "Salt Lake County"
}
interface Properties {
	address_components?: AddressComponents;
	coarse_location?: string;
	context: AddressContext;
	coordinates?: {
		lat: number;
		lon: number;
	};
	distance?: number;
	formatted_address_line?: string;
	gid: string;
	precision?: string; // "centroid"
	name: string;
}
export interface Geometry {
	type: string;
	coordinates: [number, number];
}
export interface Address {
	geometry?: Geometry;
	properties: Properties;
}
export interface GeocodeFeature {
	geometry: Geometry;
	properties: Properties;
	type: string; // "Feature"
}
export interface Query {
	"point.lat": number;
	"point.lng": number;
}
export interface Geocoding {
	attribution: string; // https://stadiamaps.com/attribution
	query: Query;
}
export interface Geocode {
	bbox: [number, number, number, number];
	features: GeocodeFeature[];
	geocoding: Geocoding;
	type: string; // "FeatureCollection"
}
