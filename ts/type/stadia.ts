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
interface AddressProperties {
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

export interface Address {
	properties: AddressProperties;
	geometry?: {
		type: string;
		coordinates: [number, number];
	};
}
