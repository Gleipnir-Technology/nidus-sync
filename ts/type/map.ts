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
	constructor(
		public address: Address,
		public location: Location,
	) {}
}
export type MoveEndEventInternal = maplibregl.MapLibreEvent<
	| maplibregl.MapMouseEvent
	| maplibregl.MapTouchEvent
	| maplibregl.MapWheelEvent
	| undefined
> & {
	isInternalUpdate: boolean;
};
