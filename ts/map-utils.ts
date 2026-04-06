import { LngLat, LngLatBounds } from "maplibre-gl";
import type { Marker } from "@/types";
import type { Location } from "@/type/api";

export function boundsMarkers(markers: Marker[]): LngLatBounds {
	let min_lat = 90;
	let min_lng = 180;
	let max_lat = -90;
	let max_lng = -180;
	markers.forEach((marker: Marker) => {
		min_lat = Math.min(marker.location.latitude, min_lat);
		min_lng = Math.min(marker.location.longitude, min_lng);
		max_lat = Math.min(marker.location.latitude, max_lat);
		max_lng = Math.min(marker.location.longitude, max_lng);
	});
	return new LngLatBounds(
		new LngLat(min_lng, min_lat),
		new LngLat(max_lng, max_lat),
	);
}
export function boundsDefault(): LngLatBounds {
	return new LngLatBounds(new LngLat(-125, 50), new LngLat(-70, 25));
}

// Helper functions (outside component)
const getBoundingBox = (points: Location[]) => {
	if (!points || points.length === 0) {
		return null;
	}

	let minLat = points[0].latitude;
	let maxLat = points[0].latitude;
	let minLng = points[0].longitude;
	let maxLng = points[0].longitude;

	for (const point of points) {
		if (point.latitude < minLat) minLat = point.latitude;
		if (point.latitude > maxLat) maxLat = point.latitude;
		if (point.longitude < minLng) minLng = point.longitude;
		if (point.longitude > maxLng) maxLng = point.longitude;
	}

	return new window.maplibregl.LngLatBounds(
		new window.maplibregl.LngLat(minLng, minLat),
		new window.maplibregl.LngLat(maxLng, maxLat),
	);
};
