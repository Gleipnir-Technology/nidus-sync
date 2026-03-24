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
		<p>{{error}}</p>
	</div>
</template>

<script setup lang="ts">
import "maplibre-gl/dist/maplibre-gl.css";
import {
	onMounted,
	onUnmounted,
	ref,
	watch,
} from "vue";
import { Bounds, Marker } from "@/types";
import maplibregl from "maplibre-gl";

interface Emits {
	(e: "load"): void;
}
interface Props {
	bounds?: Bounds;
	markers: Marker[];
	"organization-id": int;
	tegola: string;
}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
	// default bounds cover a bunch of the continental US
	bounds: {
		max: { x: -70, y: 50 },
		min: { x: -125, y: 25 },
	},
});

const error = ref<string | null>(null);
const mapContainer = ref<HTMLElement | null>(null);
const map = ref<maplibregl.Map | null>(null);
const markerInstances = ref<Map<string, maplibregl.Marker>>(new Map());
const markers = ref<Map<string, maplibregl.Marker>>(new Map());
watch(
	() => props.bounds,
	(newBounds) => {
		const bounds = new maplibregl.LngLatBounds(
			new maplibregl.LngLat(newBounds.min.lng, newBounds.min.lat),
			new maplibregl.LngLat(newBounds.max.lng, newBounds.max.lat),
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

const _bounds = () => {
	return [
		[props.bounds.min.x, props.bounds.min.y],
		[props.bounds.max.y, props.bounds.max.y],
	];
};

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

		// Wait for map to load, then add the markers
		map.value.on("load", () => {
			if (props.organizationId !== 0) {
				map.value.addSource("tegola", {
					type: "vector",
					tiles: [
						`${props.tegola}maps/nidus/{z}/{x}/{y}?id=${props.organizationId}&organization_id=${props.organizationId}`,
					],
				});

				map.value.addLayer({
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
		error.value = e;
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
			marker.setLngLat([markerData.location.lng, markerData.location.lat]);
			console.log("updated", markerData);
		} else {
			// Create a new marker
			const marker = new maplibregl.Marker({
				color: markerData.color,
				draggable: markerData.draggable,
			})
				.setLngLat([markerData.location.lng, markerData.location.lat])
				.addTo(map.value!);

			markerInstances.value.set(markerData.id, marker);
			console.log("added", markerData);
		}
	});
}
</script>
