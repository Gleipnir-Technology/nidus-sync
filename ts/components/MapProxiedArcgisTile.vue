<style scoped>
#map {
	height: 100%;
	width: 100%;
	margin-bottom: 10px;
}
.map-container {
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
import { ref, watch, onMounted, onBeforeUnmount } from "vue";
import { Point } from "@/types";
import maplibregl from "maplibre-gl";

interface Emits {
	(e: "map-click", latitude: Number, longitude: Number): void;
}
interface Props {
	location: Point;
	organizationId: Number;
	tegola: string;
	urlTiles: string;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();

const error = ref<string | null>(null);
const mapContainer = ref<HTMLElement | null>(null);
const map = ref<maplibregl.Map | null>(null);
const markerInstances = ref<Map<string, maplibrgl.Marker>>(new Map());
const markers = ref<Map<string, maplibrgl.Marker>>(new Map());

// Watch for latitude/longitude changes
watch(
	() => [props.location],
	([newLocation]) => {
		if (map.value) {
			map.value.jumpTo({
				center: [newLocation.longitude, newLocation.latitude],
				zoom: 19,
			});
		}
	},
);

const initializeMap = () => {
	if (!mapContainer.value) return;

	try {
		map.value = new maplibregl.Map({
			center: [props.location.longitude, props.location.latitude],
			container: mapContainer.value,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
			zoom: 19,
		});

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

			map.value.addSource("flyover", {
				type: "raster",
				tiles: [props.urlTiles],
			});

			map.value.addLayer({
				id: "flyover-layer",
				source: "flyover",
				type: "raster",
			});

			emit("load", { map: getCurrentInstance() });

			map.value.on("click", (e) => {
				emit("map-click", {
					lng: e.lngLat.lng,
					lat: e.lngLat.lat,
					map: getCurrentInstance(),
					point: e.point,
				});
			});
		});
	} catch (e) {
		console.error("hey dummy", e);
	}
};

const addLayer = (a) => {
	return map.value?.addLayer(a);
};

const addSource = (a, b) => {
	return map.value?.addSource(a, b);
};

const jumpTo = (args) => {
	return map.value?.jumpTo(args);
};

const once = (a, b) => {
	return map.value?.once(a, b);
};

const queryRenderedFeatures = (a) => {
	return map.value?.queryRenderedFeatures(a);
};

const fitBounds = (bounds, options) => {
	return map.value?.fitBounds(bounds, options);
};

const setLayoutProperty = (layout, property, value) => {
	return map.value?.setLayoutProperty(layout, property, value);
};

const setMarkers = (newMarkers) => {
	console.log("Setting map markers", newMarkers);
	markers.value.forEach((marker) => marker.remove());
	markers.value = newMarkers;
	for (let m of newMarkers) {
		m.addTo(map.value);
	}
};

const getCurrentInstance = () => {
	// Return an object with the public methods
	return {
		addLayer,
		addSource,
		jumpTo,
		on,
		once,
		queryRenderedFeatures,
		fitBounds,
		setLayoutProperty,
		setMarkers,
	};
};

// Expose methods to parent components
defineExpose({
	addLayer,
	addSource,
	jumpTo,
	once,
	queryRenderedFeatures,
	fitBounds,
	setLayoutProperty,
	setMarkers,
	map,
});

onMounted(() => {
	setTimeout(() => initializeMap(), 0);
});

onBeforeUnmount(() => {
	if (map.value) {
		map.value.remove();
	}
});
</script>

<style scoped>
.map-container {
	height: 100%;
	width: 100%;
}
</style>
