import type { Map as MapLibreMap } from "maplibre-gl";
import { Location } from "@/type/api";

export interface Bounds {
	min: Location;
	max: Location;
}
export interface Changes {
	updated: string[];
	unchanged: string[];
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

export interface Point {
	x: number;
	y: number;
}
