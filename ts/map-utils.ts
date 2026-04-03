import type { Marker } from "@/types";
import { LngLat, LngLatBounds } from "maplibre-gl";

export function boundsMarkers(markers: Marker[]): LngLatBounds {
	let min_lat = 90;
	let min_lng = 180;
	let max_lat = -90;
	let max_lng = -180;
	markers.forEach((marker: Marker) => {
		min_lat = Math.min(marker.location.lat, min_lat);
		min_lng = Math.min(marker.location.lng, min_lng);
		max_lat = Math.min(marker.location.lat, max_lat);
		max_lng = Math.min(marker.location.lng, max_lng);
	});
	return new LngLatBounds(
		new LngLat(min_lng, min_lat),
		new LngLat(max_lng, max_lat),
	);
}
export function boundsDefault(): LngLatBounds {
	return new LngLatBounds(new LngLat(-70, 50), new LngLat(-125, 25));
}
