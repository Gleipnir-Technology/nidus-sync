<style scoped>
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
	background: rgba(255, 255, 255, 0.85);
	backdrop-filter: blur(2px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	cursor: pointer;
	transition: background 0.2s ease;
}

.map-overlay:hover {
	background: rgba(255, 255, 255, 0.9);
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
	z-index: 999;
	background: #ffc107;
	border: none;
	color: #000;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
	font-size: 0.875rem;
	padding: 0.375rem 0.75rem;
	transition: all 0.2s;
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
	<div class="map-wrapper" ref="mapWrapper">
		<!-- Tap-to-activate overlay -->
		<div
			class="map-overlay"
			@click="activateMap"
			ref="mapOverlay"
			@touchstart.prevent="activateMap"
			v-if="!mapInteractive"
		>
			<div class="overlay-content">
				<i class="bi bi-hand-index-thumb"></i>
				<p class="mb-0">Tap to select location</p>
			</div>
		</div>

		<!-- Map container -->
		<div
			ref="mapContainer"
			class="map-container"
			:class="{ 'map-inactive': !mapInteractive }"
		></div>

		<!-- Lock/unlock indicator button -->
		<button
			v-if="mapInteractive"
			type="button"
			class="btn btn-sm map-status-btn"
			@click="deactivateMap"
			title="Lock map to enable page scrolling"
		>
			<i class="bi bi-unlock-fill"></i>
			<span class="d-none d-md-inline ms-1">Map Active</span>
		</button>
	</div>
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
interface Emits {
	(e: "update:modelValue", value: Camera): void;
	(e: "click", location: Location): void;
	(e: "marker-drag-end", location: Location): void;
}

// Props
interface Props {
	modelValue: Camera | null;
	markers?: Marker[];
}

const props = withDefaults(defineProps<Props>(), {
	markers: () => [],
});

const emit = defineEmits<Emits>();

// Refs
const clickTimeout = ref<number | null>(null);
const map: Ref<maplibregl.Map | null> = shallowRef(null);
const mapContainer = ref<HTMLDivElement | null>(null);
const mapInteractive = ref(false);
const mapMarkers: Ref<Map<string, maplibregl.Marker>> = shallowRef<
	Map<string, maplibregl.Marker>
>(new Map());
const mapOverlay = useTemplateRef("mapOverlay");
const mapWrapper = useTemplateRef("mapWrapper");

function activateMap() {
	mapInteractive.value = true;
	console.log("activated map");
	if (!map.value) {
		return;
	}
	map.value.scrollZoom.enable();
	map.value.dragPan.enable();
	map.value.touchZoomRotate.enable();
	map.value.doubleClickZoom.enable();
}

// Deactivate map interaction
function deactivateMap() {
	mapInteractive.value = false;
	if (!map.value) {
		return;
	}
	map.value.scrollZoom.disable();
	map.value.dragPan.disable();
	map.value.touchZoomRotate.disable();
	map.value.doubleClickZoom.disable();
}

// Handle clicks outside the map to deactivate
function handleOutsideClick(event: MouseEvent | TouchEvent) {
	if (!event.target) {
		console.log("Didn't click on anything");
		return;
	}
	if (!mapWrapper.value) {
		console.log("Somehow map wrapper is null");
		return;
	}
	const target = event.target as HTMLElement;
	const cls = target.className ?? "";
	if (
		mapInteractive.value &&
		mapContainer.value &&
		!(mapWrapper.value.contains(target) || cls == "map-overlay")
	) {
		console.log("deactivate map: outside map click", target, cls);
		deactivateMap();
	} else {
		console.log("click is inside the map, ignoring");
	}
}
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
	_map.on("click", (e: maplibregl.MapLayerMouseEvent) => {
		e.preventDefault();
		if (!mapInteractive.value) {
			activateMap();
			return;
		}

		// Use timeout to distinguish between click and drag
		if (clickTimeout.value) {
			clearTimeout(clickTimeout.value);
		}

		clickTimeout.value = window.setTimeout(() => {
			const location: Location = {
				latitude: e.lngLat.lat,
				longitude: e.lngLat.lng,
			};
			emit("click", location);
		}, 100);
	});

	_map.on("load", () => {
		updateModel(_map);
	});

	_map.on("moveend", (evt: MoveEndEventInternal) => {
		if (_map && !evt.isInternalUpdate) {
			updateModel(_map);
		}
	});

	_map.on("zoomend", (evt: MoveEndEventInternal) => {
		if (_map && !evt.isInternalUpdate) {
			updateModel(_map);
		}
	});

	// Listen for clicks outside the map
	document.addEventListener("mousedown", handleOutsideClick);
	document.addEventListener("touchstart", handleOutsideClick);
}

function updateModel(_map: maplibregl.Map) {
	const center = _map.getCenter();
	const newCamera: Camera = {
		location: {
			latitude: center.lat,
			longitude: center.lng,
		},
		zoom: _map.getZoom(),
	};
	emit("update:modelValue", newCamera);
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
					emit("marker-drag-end", location);
				});
			}

			mapMarkers.value.set(markerData.id, marker);
		}
	});
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
			{ duration: 1000, zoom: props.modelValue?.zoom },
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

// Watch for modelValue changes to pan to location
watch(
	() => props.modelValue,
	(newCamera) => {
		if (map.value && newCamera) {
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

// Lifecycle hooks
onMounted(() => {
	setTimeout(() => {
		initializeMap();
		updateMarkers();
	}, 0);
});

onUnmounted(() => {
	if (clickTimeout.value) {
		clearTimeout(clickTimeout.value);
	}
	document.removeEventListener("click", handleOutsideClick);
	document.removeEventListener("touchstart", handleOutsideClick);

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
