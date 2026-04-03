<style scoped>
@import url("https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");

.map-container {
	height: 100%;
	width: 100%;
}

.map-container :deep(img) {
	max-width: none;
	min-width: 0px;
	height: auto;
}
</style>

<template>
	<div ref="mapContainer" class="map-container"></div>
</template>

<script setup lang="ts">
import maplibregl from "maplibre-gl";
import type { LngLatBoundsLike, Map as MapLibreMap } from "maplibre-gl";
import { onMounted, onUnmounted, ref, type Ref, shallowRef, watch } from "vue";
import { boundsMarkers, boundsDefault } from "@/map-utils";
import type { Location, Marker } from "@/types";
import type { Camera, MoveEndEventInternal } from "@/type/map";

// Emits interface
interface Emits {
	(e: "update:modelValue", value: Camera): void;
	(e: "click", location: Location): void;
	(e: "load"): void;
	(e: "markerDragEnd", location: Location): void;
}

// Props
interface Props {
	modelValue: Camera | null;
	apiKey?: string;
	markers?: Marker[];
}

const props = withDefaults(defineProps<Props>(), {
	markers: () => [],
});

const emit = defineEmits<Emits>();

// Refs
const mapContainer = ref<HTMLDivElement | null>(null);
const map: Ref<MapLibreMap | null> = shallowRef(null);
const markerInstances: Ref<maplibregl.Marker[]> = shallowRef<
	maplibregl.Marker[]
>([]);

// Initialize map
const initializeMap = () => {
	if (!mapContainer.value) return;

	let bounds = boundsDefault();
	if (props.markers.length > 0) {
		bounds = boundsMarkers(props.markers);
	}
	const _map = new maplibregl.Map({
		bounds: bounds,
		container: mapContainer.value,
		style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
	});
	map.value = _map;
	_map.on("click", (e: maplibregl.MapLayerMouseEvent) => {
		e.preventDefault();
		console.log("internal click", e);
		emit("click", {
			lat: e.lngLat.lat,
			lng: e.lngLat.lng,
		});
	});

	_map.on("load", () => {
		emit("load");
		updateModel(_map);
	});

	_map.on("zoomend", (evt: MoveEndEventInternal) => {
		console.log("zoomend", evt);
		if (_map && !evt.isInternalUpdate) {
			updateModel(_map);
		}
	});

	_map.on("moveend", (evt: MoveEndEventInternal) => {
		console.log("moveend", evt);
		if (_map && !evt.isInternalUpdate) {
			updateModel(_map);
		}
	});
};

function updateModel(_map: maplibregl.Map) {
	const center = _map.getCenter();
	const newCamera: Camera = {
		location: center,
		zoom: _map.getZoom(),
	};
	emit("update:modelValue", newCamera);
}
// Update markers on the map
const updateMarkers = () => {
	if (!map.value) return;

	// Remove existing markers
	markerInstances.value.forEach((marker) => marker.remove());
	markerInstances.value = [];

	// Add new markers
	props.markers.forEach((markerDef) => {
		const marker = new maplibregl.Marker({
			color: markerDef.color || "#FF0000",
			draggable: markerDef.draggable ?? true,
		})
			.setLngLat(markerDef.location)
			.addTo(map.value!);

		if (markerDef.draggable ?? true) {
			marker.on("dragend", () => {
				const lngLat = marker.getLngLat();
				emit("markerDragEnd", {
					lat: lngLat.lat,
					lng: lngLat.lng,
				});
			});
		}

		markerInstances.value.push(marker);
	});

	// Frame markers if there are any
	if (props.markers.length > 0) {
		frameMarkers();
	}
};

// Frame all markers in view
const frameMarkers = () => {
	if (!map.value || props.markers.length === 0) return;

	if (props.markers.length === 1) {
		// Single marker: pan to it
		map.value.panTo(
			props.markers[0].location,
			{ duration: 1000, zoom: props.modelValue?.zoom },
			{ isInternalUpdate: true },
		);
	} else {
		// Multiple markers: fit bounds
		const bounds = new maplibregl.LngLatBounds();
		props.markers.forEach((marker) => {
			bounds.extend([marker.location.lng, marker.location.lat]);
		});
		map.value.fitBounds(
			bounds,
			{ padding: 10, duration: 1000 },
			{ isInternalUpdate: true },
		);
	}
};

// Watch for modelValue changes to pan to location
watch(
	() => props.modelValue,
	(newLocation) => {
		if (map.value && newLocation) {
			map.value.panTo(
				newLocation.location,
				{ duration: 1000 },
				{ isInternalUpdate: true },
			);
		}
	},
	{ deep: true },
);

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
	if (map.value) {
		map.value.remove();
	}
});
</script>
