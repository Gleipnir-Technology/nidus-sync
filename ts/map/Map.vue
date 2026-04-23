<style scoped>
.map {
	position: absolute;
	top: 0;
	bottom: 0;
	left: 0;
	right: 0;
	height: 100%;
	width: 100%;
}
</style>

<template>
	<div ref="mapDiv" class="map" v-bind="$attrs"></div>
	<slot />
</template>

<script setup lang="ts">
import "maplibre-gl/dist/maplibre-gl.css";
import maplibregl from "maplibre-gl";
import {
	onBeforeUnmount,
	onMounted,
	provide,
	ref,
	type Ref,
	shallowRef,
} from "vue";

import type { Bounds } from "@/type/api";

interface Props {
	bounds?: Bounds;
	center?: maplibregl.LngLatLike;
	zoom?: number;
}

const props = withDefaults(defineProps<Props>(), {});

const mapDiv = ref<HTMLElement | null>(null);
const map: Ref<maplibregl.Map | null> = shallowRef(null);

// Provide the map instance to children
provide("map", map);

// Registry for tracking child components
const sources = new Map();
const layers = new Map();

provide("registerSource", (id: string, config: any) => {
	console.log("register source", id, config);
	sources.set(id, config);
	if (map.value && map.value.loaded()) {
		if (!map.value.getSource(id)) {
			map.value.addSource(id, config);
		}
	}
});

provide("unregisterSource", (id: string) => {
	console.log("unregister source", id);
	sources.delete(id);
	if (map.value && map.value?.getSource(id)) {
		map.value.removeSource(id);
	}
});

provide("registerLayer", (id: string, config: any) => {
	console.log("register layer", id, config);
	layers.set(id, config);
	if (map.value && map.value.loaded()) {
		if (!map.value.getLayer(id)) {
			map.value.addLayer(config);
		}
	}
});

provide("unregisterLayer", (id: string) => {
	console.log("unregister layer", id);
	layers.delete(id);
	if (map.value?.getLayer(id)) {
		map.value.removeLayer(id);
	}
});

function initializeMap() {
	if (!mapDiv.value) return;

	console.log("initializing map...");
	const _map = new maplibregl.Map({
		container: mapDiv.value,
		center: props.center,
		style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
		zoom: props.zoom,
	});

	// When map loads, add all registered sources/layers
	_map.on("load", () => {
		console.log("map loaded.");
		sources.forEach((config, id) => {
			console.log("adding source", id, config);
			if (!_map.getSource(id)) {
				_map.addSource(id, config);
			}
		});

		layers.forEach((config, id) => {
			console.log("adding layer", id, config);
			if (!_map.getLayer(id)) {
				_map.addLayer(config);
			}
		});
	});
	map.value = _map;
}
onMounted(() => {
	initializeMap();
});

onBeforeUnmount(() => {
	if (map.value) {
		map.value.remove();
	}
});
</script>
