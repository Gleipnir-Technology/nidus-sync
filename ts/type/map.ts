import maplibregl from "maplibre-gl";
import type { Address, Location } from "@/type/api";

export interface Camera {
	location: Location;
	zoom: number;
}
export interface Locator {
	address: Address;
	location: Location;
}
export type MoveEndEventInternal = maplibregl.MapLibreEvent<
	| maplibregl.MapMouseEvent
	| maplibregl.MapTouchEvent
	| maplibregl.MapWheelEvent
	| undefined
> & {
	isInternalUpdate: boolean;
};
