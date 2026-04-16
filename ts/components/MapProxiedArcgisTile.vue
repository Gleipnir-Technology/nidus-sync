<style scoped lang="scss">
@import url("https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
.map-wrapper {
	position: relative;
	width: 100%;
	height: 100%;
	border-radius: 10px;
	overflow: hidden;
}

.map {
	width: 100%;
	height: 100%;
	transition: filter 0.2s ease;
}

.map-inactive {
	filter: brightness(0.95);
}

.map-overlay {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(255, 255, 255, 0.45);
	backdrop-filter: blur(2px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 2;
	cursor: pointer;
	transition: background 0.2s ease;
}

.map-overlay:hover {
	background: rgba(255, 255, 255, 0.65);
}

.overlay-content {
	text-align: center;
	color: #0d6efd;
	font-size: 1.1rem;
	font-weight: 500;
	user-select: none;
	pointer-events: none;
	animation: pulse 2s ease-in-out infinite;
}

.overlay-content i {
	display: block;
	font-size: 3rem;
	margin-bottom: 0.5rem;
}

@keyframes pulse {
	0%,
	100% {
		opacity: 1;
		transform: scale(1);
	}
	50% {
		opacity: 0.8;
		transform: scale(1.05);
	}
}

.map-status-btn {
	position: absolute;
	top: 10px;
	right: 10px;
	z-index: 2;
	border: none;
	color: #000;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
	font-size: 0.875rem;
	padding: 0.375rem 0.75rem;
	transition: all 0.2s;
}
.map-status-btn.locked {
	background: $warning;
}
.map-status-btn.unlocked {
	background: $primary;
}

.map-status-btn:hover {
	background: #ffb300;
	box-shadow: 0 3px 8px rgba(0, 0, 0, 0.3);
	transform: translateY(-1px);
}

.map-status-btn:active {
	transform: translateY(0);
}

/* Mobile optimizations */
@media (max-width: 768px) {
	.overlay-content {
		font-size: 1rem;
	}

	.overlay-content i {
		font-size: 2.5rem;
	}

	.map-status-btn {
		padding: 0.5rem;
		font-size: 1rem;
	}
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
	<div ref="mapContainer" class="map-container"></div>
</template>

<script setup lang="ts">
import "maplibre-gl/dist/maplibre-gl.css";
import maplibregl from "maplibre-gl";
import {
	onMounted,
	onBeforeUnmount,
	ref,
	type Ref,
	shallowRef,
	watch,
} from "vue";

import LayersControl from "@/components/LayersControl";
import { boundsMarkers, boundsDefault } from "@/map-utils";
import { MapClickEvent, Marker, Point } from "@/types";
import type { Location } from "@/type/api";
import type { Camera, MoveEndEventInternal } from "@/type/map";

interface Emits {
	(e: "map-click", event: MapClickEvent): void;
	(e: "update:modelValue", value: Camera): void;
}
interface Props {
	initialCamera?: Camera;
	modelValue: Camera;
	markers?: Marker[];
	organizationId: Number;
	tegola: string;
	urlTiles: string;
}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
	markers: () => [],
});

const error = ref<string | null>(null);
const isLoaded = ref<boolean>(false);
const map: Ref<maplibregl.Map | null> = shallowRef(null);
const mapContainer = ref<HTMLElement | null>(null);
const markerInstances = ref<Map<string, maplibregl.Marker>>(new Map());
const mapMarkers: Ref<Map<string, maplibregl.Marker>> = shallowRef<
	Map<string, maplibregl.Marker>
>(new Map());

// Frame all markers in view
function frameMarkers() {
	if (!map.value || props.markers.length === 0) return;

	if (props.markers.length === 1) {
		// Single marker: pan to it
		// If we are zoomed way out we are likely in the default state antd therefore should zoom in a bunch
		// for the framing.
		const zoom = props.modelValue.zoom > 1 ? props.modelValue.zoom : 15;
		console.log(
			"framing single marker",
			isLoaded.value,
			props.markers[0].location,
			props.modelValue.zoom,
			zoom,
		);

		// Defer this until the map is loaded or we'll drop updates
		if (map.value) {
			if (isLoaded.value) {
				panToLocation(props.markers[0].location, zoom);
			} else {
				map.value.on("load", () => {
					panToLocation(props.markers[0].location, zoom);
				});
			}
		} else {
			console.error("Can't frame markers before the map is created");
		}
	} else {
		// Multiple markers: fit bounds
		console.log("framing multiple markers", props.markers);
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
}
const initializeMap = () => {
	if (!mapContainer.value) return;

	try {
		const _map = new maplibregl.Map({
			container: mapContainer.value,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
			zoom: 19,
		});
		if (props.markers.length > 0) {
			console.log("initial map fitting initial markers", props.markers);
			_map.fitBounds(boundsMarkers(props.markers));
		} else if (
			props.initialCamera &&
			(props.initialCamera.location.latitude ||
				props.initialCamera.location.longitude)
		) {
			console.log("initial map jump to initial camera", props.initialCamera);
			_map.jumpTo({
				center: [
					props.initialCamera.location.longitude,
					props.initialCamera.location.latitude,
				],
				zoom: props.initialCamera.zoom,
			});
		} else if (
			props.modelValue.location.latitude != 0 ||
			props.modelValue.location.longitude != 0
		) {
			console.log("initial map jump to initial model", props.modelValue);
			_map.jumpTo({
				center: [
					props.modelValue.location.longitude,
					props.modelValue.location.latitude,
				],
				zoom: props.modelValue.zoom,
			});
		} else {
			const bounds = boundsDefault();
			console.log("initial map fitting default bounds", bounds);
			_map.fitBounds(bounds);
		}
		map.value = _map;
		_map.on("load", () => {
			console.log("proxied tile loaded");
			isLoaded.value = true;
			if (props.organizationId !== 0) {
				_map.addSource("tegola", {
					type: "vector",
					tiles: [
						`${props.tegola}maps/nidus/{z}/{x}/{y}?id=${props.organizationId}&organization_id=${props.organizationId}`,
					],
				});
				_map.addLayer({
					id: "service-area",
					source: "tegola",
					"source-layer": "service-area-bounds",
					type: "line",
					paint: {
						"line-color": "#f00",
					},
				});
			}

			_map.addSource("flyover", {
				type: "raster",
				tiles: [props.urlTiles],
			});

			_map.addLayer({
				id: "flyover-layer",
				source: "flyover",
				type: "raster",
			});

			_map.on("click", (e) => {
				emit("map-click", {
					location: {
						latitude: e.lngLat.lat,
						longitude: e.lngLat.lng,
					},
					map: _map,
					point: e.point,
				});
			});
			console.log("MapProxiedArcgisTile loaded");
		});
		_map.addControl(new maplibregl.NavigationControl(), "top-left");
		_map.addControl(
			new LayersControl({
				title: "layers",
				customLabels: { "countries-fill": "Countries" },
			}),
		);
		console.log("MapProxiedArcgisTile initialized");
	} catch (e) {
		console.error("hey dummy", e);
	}
};
function panToLocation(location: Location, zoom: number) {
	if (!map.value) return;
	map.value.panTo(
		{
			lat: props.markers[0].location.latitude,
			lng: props.markers[0].location.longitude,
		},
		{ duration: 1000, zoom: zoom },
		{ isInternalUpdate: true },
	);
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
			marker.setDraggable(markerData.draggable ?? false);
		} else {
			// Create new marker
			const el = document.createElement("div");
			el.className = "custom-marker";
			el.style.backgroundColor = markerData.color ?? "#FF0000";
			el.style.width = "25px";
			el.style.height = "25px";
			el.style.borderRadius = "50%";
			el.style.border = "3px solid white";
			el.style.boxShadow = "0 2px 4px rgba(0,0,0,0.3)";
			el.style.cursor = markerData.draggable ? "move" : "pointer";

			marker = new maplibregl.Marker({
				element: el,
				draggable: markerData.draggable ?? false,
			})
				.setLngLat([
					markerData.location.longitude,
					markerData.location.latitude,
				])
				.addTo(map.value!);

			// Handle marker drag end
			if (markerData.draggable) {
				marker.on("dragend", () => {
					const lngLat = marker!.getLngLat();
					const location: Location = {
						latitude: lngLat.lat,
						longitude: lngLat.lng,
					};
					//emit("marker-drag-end", location);
				});
			}

			mapMarkers.value.set(markerData.id, marker);
		}
	});
	frameMarkers();
};

onMounted(() => {
	setTimeout(() => initializeMap(), 0);
});

onBeforeUnmount(() => {
	if (map.value) {
		map.value.remove();
	}
});
watch(
	() => props.modelValue,
	(newCamera) => {
		if (map.value && newCamera) {
			console.log("panning based on model change", newCamera);
			map.value.panTo(
				{
					lat: newCamera.location.latitude,
					lng: newCamera.location.longitude,
				},
				{ duration: 1000, zoom: newCamera.zoom },
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
</script>
