<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import maplibregl from "maplibre-gl";

const mapContainer = ref<HTMLDivElement | null>(null);
let map: maplibregl.Map | null = null;

onMounted(() => {
	if (!mapContainer.value) return;

	map = new maplibregl.Map({
		container: mapContainer.value,
		style: "https://demotiles.maplibre.org/style.json",
		center: [-74.5, 40], // [lng, lat]
		zoom: 9,
	});

	// Add a marker as an example
	new maplibregl.Marker().setLngLat([-74.5, 40]).addTo(map);
});

onUnmounted(() => {
	map?.remove();
});
</script>

<template>
	<div ref="mapContainer" class="map-container"></div>
</template>

<style scoped>
.map-container {
	width: 100%;
	height: 500px;
}
</style>
