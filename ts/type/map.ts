import maplibregl from "maplibre-gl";
import type { Location } from "@/type/api";

export interface Camera {
	location: Location;
	zoom: number;
}

export type MoveEndEventInternal = maplibregl.MapLibreEvent<
	| maplibregl.MapMouseEvent
	| maplibregl.MapTouchEvent
	| maplibregl.MapWheelEvent
	| undefined
> & {
	isInternalUpdate: boolean;
};
