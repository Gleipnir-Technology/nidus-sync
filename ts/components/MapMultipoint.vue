<style scoped>
#map {
	height: 100%;
	width: 100%;
	margin-bottom: 10px;
}
.map-multipoint {
	height: 100%;
	width: 100%;
}
</style>
<template>
	<div v-if="error == null">
		<div ref="mapContainer" class="map-multipoint"></div>
	</div>
	<div v-else>
		<h1>Map failed to load</h1>
		<p>{{ error }}</p>
	</div>
</template>

<script setup lang="ts">
import "maplibre-gl/dist/maplibre-gl.css";
import type { LngLatBoundsLike, Map as MapLibreMap } from "maplibre-gl";
import maplibregl from "maplibre-gl";
import {
	computed,
	onMounted,
	onUnmounted,
	ref,
	shallowRef,
	watch,
	type Ref,
} from "vue";
import { Marker } from "@/types";
import type { Bounds } from "@/type/api";

interface Emits {}
interface Props {
	bounds?: Bounds;
	markers: Marker[];
	organizationId: number;
	tegola: string;
}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
	// default bounds cover a bunch of the continental US
	bounds: () => {
		return {
			max: { longitude: -70, latitude: 50 },
			min: { longitude: -125, latitude: 25 },
		};
	},
});

const boundsSafe = props.bounds as Bounds;
const error = ref<string | null>(null);
const mapContainer = ref<HTMLElement | null>(null);
const map: Ref<MapLibreMap | null> = shallowRef(null);
const markerInstances = shallowRef<Map<string, maplibregl.Marker>>(new Map());
watch(
	() => props.bounds,
	(newBounds) => {
		const bounds = new maplibregl.LngLatBounds(
			new maplibregl.LngLat(newBounds.min.longitude, newBounds.min.latitude),
			new maplibregl.LngLat(newBounds.max.longitude, newBounds.max.latitude),
		);
		if (map.value == null) {
			return;
		}
		map.value.fitBounds(bounds, {
			padding: 50,
		});
	},
	{ deep: true },
);
watch(
	() => props.markers,
	(newMarkers) => {
		updateMarkers(newMarkers);
	},
	{ deep: true },
);

function _bounds(): LngLatBoundsLike {
	return new maplibregl.LngLatBounds(
		new maplibregl.LngLat(boundsSafe.min.longitude, boundsSafe.min.latitude),
		new maplibregl.LngLat(boundsSafe.max.longitude, boundsSafe.max.latitude),
	);
}

const _initializeMap = () => {};

// Lifecycle
onMounted(() => {
	if (!mapContainer.value) return;
	const bounds = _bounds();

	try {
		map.value = new maplibregl.Map({
			bounds: bounds,
			container: mapContainer.value,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
		});
		const mapInstance = map.value;

		// Wait for map to load, then add the markers
		mapInstance.on("load", () => {
			if (props.organizationId !== 0) {
				mapInstance.addSource("tegola", {
					type: "vector",
					tiles: [
						`${props.tegola}maps/nidus/{z}/{x}/{y}?id=${props.organizationId}&organization_id=${props.organizationId}`,
					],
				});

				mapInstance.addLayer({
					id: "service-area",
					source: "tegola",
					"source-layer": "service-area-bounds",
					type: "line",
					paint: {
						"line-color": "#f00",
					},
				});
			}
			updateMarkers(props.markers);
		});
	} catch (e) {
		error.value = e instanceof Error ? e.message : "an error ocurred";
		console.error("Error on map multipoint init", e);
		//throw new Error(error.value);
	}
});

onUnmounted(() => {
	// Remove all markers
	markerInstances.value.forEach((marker) => marker.remove());
	markerInstances.value.clear();

	// Free OpenGL context
	map.value?.remove();
	map.value = null;
});

function updateMarkers(markers: Marker[]) {
	const newMarkerIds = new Set(markers.map((m) => m.id));

	if (map.value == null) {
		console.log("refusing to add markers until map is set");
		return;
	}

	// Remove markers that no longer exist
	markerInstances.value.forEach((marker, id) => {
		if (!newMarkerIds.has(id)) {
			marker.remove();
			markerInstances.value.delete(id);
		}
	});

	// Add or update markers
	markers.forEach((markerData) => {
		if (markerInstances.value.has(markerData.id)) {
			// Update existing marker position
			const marker = markerInstances.value.get(markerData.id)!;
			marker.setLngLat([
				markerData.location.longitude,
				markerData.location.latitude,
			]);
			console.log("updated", markerData);
		} else {
			// Create a new marker
			const marker = new maplibregl.Marker({
				color: markerData.color,
				draggable: markerData.draggable,
			})
				.setLngLat([
					markerData.location.longitude,
					markerData.location.latitude,
				])
				.addTo(map.value!);

			markerInstances.value.set(markerData.id, marker);
			console.log("added", markerData);
		}
	});
}
</script>
