import maplibregl from "maplibre-gl";
import type { Marker } from "@/types";
import type { Bounds, Location } from "@/type/api";

export function boundsDefault(): maplibregl.LngLatBounds {
	return new maplibregl.LngLatBounds(
		new maplibregl.LngLat(-125, 50),
		new maplibregl.LngLat(-70, 25),
	);
}

export function boundsFromAPI(b: Bounds): maplibregl.LngLatBounds {
	return new maplibregl.LngLatBounds(
		new maplibregl.LngLat(b.min.longitude, b.max.latitude),
		new maplibregl.LngLat(b.max.longitude, b.min.latitude),
	);
}

export function boundsMarkers(markers: Marker[]): maplibregl.LngLatBounds {
	let min_lat = 90;
	let min_lng = 180;
	let max_lat = -90;
	let max_lng = -180;
	markers.forEach((marker: Marker) => {
		min_lat = Math.min(marker.location.latitude, min_lat);
		min_lng = Math.min(marker.location.longitude, min_lng);
		max_lat = Math.max(marker.location.latitude, max_lat);
		max_lng = Math.max(marker.location.longitude, max_lng);
	});
	return new maplibregl.LngLatBounds(
		new maplibregl.LngLat(min_lng, min_lat),
		new maplibregl.LngLat(max_lng, max_lat),
	);
}
export function boundsWithPadding(
	min: Location,
	max: Location,
	padding: number,
) {
	return new maplibregl.LngLatBounds(
		new maplibregl.LngLat(min.longitude - padding, min.latitude - padding),
		new maplibregl.LngLat(max.longitude + padding, max.latitude + padding),
	);
}
// Helper functions (outside component)
const getBoundingBox = (points: Location[]) => {
	if (!points || points.length === 0) {
		return null;
	}

	let maxLat = points[0].latitude;
	let maxLng = points[0].longitude;
	let minLat = points[0].latitude;
	let minLng = points[0].longitude;

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
