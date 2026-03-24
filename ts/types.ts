export interface Address {
	locality: string;
	number: string;
	street: string;
}

export interface Communication {
	created: string;
	id: string;
	public_report: PublicReport | null;
	type: string;
}
export interface Point {
	lat: Number;
	lng: Number;
}
export interface Bounds {
	min: Point;
	max: Point;
}
export interface Marker {
	color: string;
	draggable: boolean;
	id: string;
	location: Point;
}

export interface PublicReport {
	created: string;
	type: string;
}

export interface Signal {}
