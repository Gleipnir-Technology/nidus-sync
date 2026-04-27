import maplibregl from "maplibre-gl";
import { Address, Location } from "@/type/api";

export class Camera {
	location: Location;
	zoom: number;
	constructor(location: Location = new Location(), zoom: number = 0) {
		this.location = location;
		this.zoom = zoom;
	}
}
export class Locator {
	address: Address;
	location: Location;
	constructor(
		address: Address = new Address(),
		location: Location = new Location(),
	) {
		this.address = address;
		this.location = location;
	}
}
export type MoveEndEventInternal = maplibregl.MapLibreEvent<
	| maplibregl.MapMouseEvent
	| maplibregl.MapTouchEvent
	| maplibregl.MapWheelEvent
	| undefined
> & {
	isInternalUpdate: boolean;
};
