// Interface definitions
interface AddressComponents {
	number?: string;
	street?: string;
	city?: string;
	state?: string;
	zip?: string;
}

interface AddressProperties {
	gid: string;
	name: string;
	coarse_location?: string;
	formatted_address_line?: string;
	address_components?: AddressComponents;
	coordinates?: {
		lat: number;
		lon: number;
	};
}

export interface Address {
	properties: AddressProperties;
	geometry?: {
		type: string;
		coordinates: [number, number];
	};
}
