<template>
	<div ref="mapContainer" class="map-container mb-4"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

interface Props {
	organizationId: string;
	tegola: string;
	xmin: number;
	ymin: number;
	xmax: number;
	ymax: number;
}

const props = defineProps<Props>();
const mapContainer = ref<HTMLDivElement | null>(null);
let map: maplibregl.Map | null = null;

const initMap = () => {
	if (!mapContainer.value) return;

	map = new maplibregl.Map({
		container: mapContainer.value,
		style: `${props.tegola}/maps/basic/style.json`,
		center: [(props.xmin + props.xmax) / 2, (props.ymin + props.ymax) / 2],
		zoom: 10,
	});

	// Add service area bounds
	map.on("load", () => {
		if (!map) return;

		map.addSource("service-area", {
			type: "geojson",
			data: {
				type: "Feature",
				geometry: {
					type: "Polygon",
					coordinates: [
						[
							[props.xmin, props.ymin],
							[props.xmax, props.ymin],
							[props.xmax, props.ymax],
							[props.xmin, props.ymax],
							[props.xmin, props.ymin],
						],
					],
				},
				properties: {},
			},
		});

		map.addLayer({
			id: "service-area-fill",
			type: "fill",
			source: "service-area",
			paint: {
				"fill-color": "#088",
				"fill-opacity": 0.2,
			},
		});

		map.addLayer({
			id: "service-area-outline",
			type: "line",
			source: "service-area",
			paint: {
				"line-color": "#088",
				"line-width": 2,
			},
		});
	});
};

onMounted(() => {
	initMap();
});

// Clean up on unmount
onUnmounted(() => {
	if (map) {
		map.remove();
	}
});
</script>

<style scoped>
.map-container {
	width: 100%;
	height: 400px;
	border-radius: 8px;
	overflow: hidden;
}
</style>
