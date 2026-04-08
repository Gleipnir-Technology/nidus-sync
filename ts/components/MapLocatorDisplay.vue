<style scoped>
@import url("https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
.map {
	width: 100%;
	height: 100%;
	transition: filter 0.2s ease;
}

/* Ensure map fills container on all devices */
:deep(.maplibregl-map) {
	width: 100%;
	height: 100%;
}

:deep(.maplibregl-canvas) {
	width: 100%;
	height: 100%;
}
</style>

<template>
	<!-- Map container -->
	<div ref="mapContainer" class="map-container"></div>
</template>

<script setup lang="ts">
import maplibregl from "maplibre-gl";
import {
	onMounted,
	onUnmounted,
	ref,
	type Ref,
	shallowRef,
	useTemplateRef,
	watch,
} from "vue";
import { boundsMarkers, boundsDefault } from "@/map-utils";
import type { Marker } from "@/types";
import type { Location } from "@/type/api";
import type { Camera, MoveEndEventInternal } from "@/type/map";

// Emits interface
interface Emits {}

// Props
interface Props {
	markers?: Marker[];
}

const props = withDefaults(defineProps<Props>(), {
	markers: () => [],
});

// Refs
const map: Ref<maplibregl.Map | null> = shallowRef(null);
const mapContainer = ref<HTMLDivElement | null>(null);
const mapMarkers: Ref<Map<string, maplibregl.Marker>> = shallowRef<
	Map<string, maplibregl.Marker>
>(new Map());

// Initialize map
function initializeMap() {
	if (!mapContainer.value) return;

	let bounds = boundsDefault();
	if (props.markers.length > 0) {
		bounds = boundsMarkers(props.markers);
	}
	const _map = new maplibregl.Map({
		bounds: bounds,
		container: mapContainer.value,
		style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
		// Disable interactions by default
		doubleClickZoom: false,
		dragPan: false,
		scrollZoom: false,
		touchZoomRotate: false,
	});
	_map.addControl(new maplibregl.NavigationControl(), "top-left");
	map.value = _map;
}

// Update markers on the map
const updateMarkers = () => {
	if (!map.value) return;

	// Remove markers that no longer exist
	const currentMarkerIds = new Set(props.markers.map((m) => m.id));
	for (const [id, marker] of mapMarkers.value) {
		if (!currentMarkerIds.has(id)) {
			marker.remove();
			mapMarkers.value.delete(id);
		}
	}

	// Add or update markers
	props.markers.forEach((markerData) => {
		let marker = mapMarkers.value.get(markerData.id);

		if (marker) {
			// Update existing marker
			marker.setLngLat([
				markerData.location.longitude,
				markerData.location.latitude,
			]);
			marker.setDraggable(false);
		} else {
			marker = new maplibregl.Marker({
				draggable: false,
			})
				.setLngLat([
					markerData.location.longitude,
					markerData.location.latitude,
				])
				.addTo(map.value!);

			mapMarkers.value.set(markerData.id, marker);
		}
	});
	frameMarkers();
};

// Frame all markers in view
const frameMarkers = () => {
	if (!map.value || props.markers.length === 0) return;

	if (props.markers.length === 1) {
		// Single marker: pan to it
		map.value.panTo(
			{
				lat: props.markers[0].location.latitude,
				lng: props.markers[0].location.longitude,
			},
			{ duration: 1000, zoom: 15 },
			{ isInternalUpdate: true },
		);
	} else {
		// Multiple markers: fit bounds
		const bounds = new maplibregl.LngLatBounds();
		props.markers.forEach((marker) => {
			bounds.extend([marker.location.longitude, marker.location.latitude]);
		});
		map.value.fitBounds(
			bounds,
			{ padding: 10, duration: 1000 },
			{ isInternalUpdate: true },
		);
	}
};

// Watch for markers changes
watch(
	() => props.markers,
	() => {
		updateMarkers();
	},
	{ deep: true },
);

// Lifecycle hooks
onMounted(() => {
	setTimeout(() => {
		initializeMap();
		updateMarkers();
	}, 0);
});

onUnmounted(() => {
	// Remove all markers
	mapMarkers.value.forEach((marker) => marker.remove());
	mapMarkers.value.clear();

	// Remove map
	if (map.value) {
		map.value.remove();
		map.value = null;
	}
});
</script>
