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
	<div class="map-wrapper" ref="mapWrapper">
		<!-- Map container -->
		<div ref="mapContainer" class="map-container"></div>

		<!-- Lock/unlock indicator button -->
		<button
			v-if="mapInteractive"
			type="button"
			class="btn btn-sm map-status-btn unlocked"
			@click="deactivateMap"
			title="Lock map to enable page scrolling"
		>
			<i class="bi bi-unlock-fill"></i>
			<span class="d-none d-md-inline ms-1">Map Active</span>
		</button>
		<button
			v-if="!mapInteractive"
			type="button"
			class="btn btn-sm map-status-btn locked"
			@click="activateMap"
			title="Unlock map to enable map pan/zoom"
		>
			<i class="bi bi-lock-fill"></i>
			<span class="d-none d-md-inline ms-1">Map Locked</span>
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
import { boundsMarkers, boundsDefault } from "@/map/util";
import type { MapClickEvent, Marker } from "@/types";
import type { Location } from "@/type/api";
import type { Camera, MoveEndEventInternal } from "@/type/map";

// Emits interface
interface Emits {
	(e: "update:modelValue", value: Camera): void;
	(e: "click", event: MapClickEvent): void;
	(e: "marker-drag-end", location: Location): void;
}

// Props
interface Props {
	initialCamera?: Camera;
	modelValue: Camera;
	markers?: Marker[];
	useSatellite?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
	markers: () => [],
	useSatellite: false,
});

const emit = defineEmits<Emits>();

// Refs
const clickTimeout = ref<number | null>(null);
const isLoaded = ref<boolean>(false);
const map: Ref<maplibregl.Map | null> = shallowRef(null);
const mapContainer = ref<HTMLDivElement | null>(null);
const mapInteractive = ref(false);
const mapMarkers: Ref<Map<string, maplibregl.Marker>> = shallowRef<
	Map<string, maplibregl.Marker>
>(new Map());
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

// Initialize map
function initializeMap() {
	if (!mapContainer.value) return;

	let style = "https://tiles.stadiamaps.com/styles/alidade_smooth.json";
	if (props.useSatellite) {
		style = "https://tiles.stadiamaps.com/styles/alidade_satellite.json";
	}
	let map_options: maplibregl.MapOptions = {
		container: mapContainer.value,
		style: style,
		// Disable interactions by default
		doubleClickZoom: false,
		dragPan: false,
		scrollZoom: false,
		touchZoomRotate: false,
	};
	if (props.markers.length > 0) {
		if (props.markers.length == 1) {
			const m = props.markers[0];
			map_options.center = [m.location.longitude, m.location.latitude];
			map_options.zoom = 15;
			console.log(
				"initial map fitting single marker:",
				m,
				"location:",
				m.location,
				"zoom:",
				15,
			);
		} else {
			const bounds = boundsMarkers(props.markers);
			console.log(
				"initial map fitting initial markers:",
				props.markers,
				"bounds:",
				bounds,
			);
			map_options.bounds = bounds;
		}
	} else if (
		props.initialCamera &&
		(props.initialCamera.location.latitude ||
			props.initialCamera.location.longitude)
	) {
		console.log("initial map jump to initial camera", props.initialCamera);
		map_options.center = [
			props.initialCamera.location.longitude,
			props.initialCamera.location.latitude,
		];
		map_options.zoom = props.initialCamera.zoom;
	} else if (
		props.modelValue.location.latitude != 0 ||
		props.modelValue.location.longitude != 0
	) {
		console.log("initial map jump to initial model", props.modelValue);
		map_options.center = [
			props.modelValue.location.longitude,
			props.modelValue.location.latitude,
		];
		map_options.zoom = props.modelValue.zoom;
	} else {
		const bounds = boundsDefault();
		console.log("initial map fitting default bounds", bounds);
		map_options.bounds = bounds;
	}
	const _map = new maplibregl.Map(map_options);
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
			emit("click", {
				location: {
					latitude: e.lngLat.lat,
					longitude: e.lngLat.lng,
				},
				map: _map,
				point: e.point,
			});
		}, 100);
	});

	_map.on("load", () => {
		// It's possible at this point that the client changed the camera while the map
		// was loading. If that's the case we need to handle that change now.
		console.log("map load complete");
		isLoaded.value = true;
		// Delay this by a tick so that other load handlers fire first
		// This allows updates to the camera model that happened during the load to fire
		// and jump the camera to a new location before doing this update.
		setTimeout(() => {
			updateModel(_map);
		}, 1);
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
	frameMarkers();
};

// Frame all markers in view
function frameMarkers() {
	if (!map.value || props.markers.length === 0 || !isLoaded.value) return;

	if (props.markers.length === 1) {
		// Single marker: pan to it
		// If we are zoomed way out we are likely in the default state antd therefore should zoom in a bunch
		// for the framing.
		const zoom = props.modelValue.zoom > 1 ? props.modelValue.zoom : 15;
		console.log(
			"framing single marker. location:",
			props.markers[0].location,
			"model zoom: ",
			props.modelValue.zoom,
			"calculated zoom: ",
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
		if (map.value) {
			console.log("framing multiple markers", isLoaded.value, props.markers);
			if (isLoaded.value) {
				panToMarkers(props.markers);
			} else {
				map.value.on("load", () => {
					panToMarkers(props.markers);
				});
			}
		} else {
			console.error("Can't frame multiple markers before the map is created");
		}
	}
}
function panToMarkers(markers: Marker[]) {
	setTimeout(() => {
		if (!map.value) return;
		const bounds = boundsMarkers(markers);
		map.value.fitBounds(
			bounds,
			{ padding: 100, duration: 1000 },
			{ isInternalUpdate: true },
		);
		console.log("fitting map to bounds", bounds);
	}, 1);
}
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

// Watch for modelValue changes to pan to location
watch(
	() => props.modelValue,
	(newCamera) => {
		if (map.value) {
			if (isLoaded.value) {
				console.log("panning based on model change", newCamera);
				map.value.panTo(
					{
						lat: newCamera.location.latitude,
						lng: newCamera.location.longitude,
					},
					{ duration: 1000, zoom: newCamera.zoom },
					{ isInternalUpdate: true },
				);
			} else {
				console.log("delaying jump until loaded", newCamera);
				map.value.once("load", () => {
					if (!map.value) return;
					map.value.jumpTo(
						{
							center: {
								lat: newCamera.location.latitude,
								lng: newCamera.location.longitude,
							},
							zoom: newCamera.zoom,
						},
						{ isInternalUpdate: true },
					);
				});
			}
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
