<template>
	<div ref="mapContainer" class="map-container"></div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from "vue";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

const props = defineProps({
	latitude: {
		type: Number,
		required: true,
	},
	longitude: {
		type: Number,
		required: true,
	},
	organizationId: {
		type: Number,
		default: 0,
	},
	tegola: {
		type: String,
		default: "",
	},
	urlTiles: {
		type: String,
		default: "",
	},
});

const emit = defineEmits(["load", "map-click"]);

const mapContainer = ref(null);
const map = ref(null);
const markers = ref([]);

// Watch for latitude/longitude changes
watch(
	() => [props.latitude, props.longitude],
	([newLat, newLng]) => {
		if (map.value) {
			map.value.jumpTo({
				center: [newLng, newLat],
				zoom: 19,
			});
		}
	},
);

const initializeMap = () => {
	if (!mapContainer.value) return;

	try {
		map.value = new maplibregl.Map({
			center: [props.longitude, props.latitude],
			container: mapContainer.value,
			//style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
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
	} catch(e) {
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
