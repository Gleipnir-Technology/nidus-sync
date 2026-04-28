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
	watch,
} from "vue";

import { boundsDefault } from "@/map/util";
import type { MapClickEvent, Marker } from "@/types";
import type { Location } from "@/type/api";

export type LngLatLike = maplibregl.LngLatLike;
export type LngLatBounds = maplibregl.LngLatBounds;
interface Emits {
	(e: "marker-drag-end", location: Location): void;
}
interface Props {
	bounds?: LngLatBounds;
	center?: LngLatLike;
	cursor?: string;
	markers?: Marker[];
	zoom?: number;
}

const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
	bounds: boundsDefault,
	cursor: "",
	markers: () => [],
});

const mapDiv = ref<HTMLElement | null>(null);
const map: Ref<maplibregl.Map | null> = shallowRef(null);
const mapMarkers: Ref<Map<string, maplibregl.Marker>> = shallowRef<
	Map<string, maplibregl.Marker>
>(new Map());

// Provide the map instance to children
provide("map", map);

// Registry for tracking child components
const ons = new Map();
const onces = new Map();
const sources = new Map();
const layers = new Map();

type OnCallbackFunc = () => void;
provide(
	"registerOn",
	(
		eventname: keyof maplibregl.MapLayerEventType,
		layerid: string,
		callback: OnCallbackFunc,
	) => {
		console.log("register map.on", eventname, layerid);
		ons.set(`${eventname}.${layerid}`, {
			callback: callback,
			eventname: eventname,
			layerid: layerid,
		});
		if (map.value && map.value.loaded()) {
			map.value.on(eventname, layerid, callback);
		}
	},
);
provide(
	"registerOnce",
	(
		eventname: keyof maplibregl.MapLayerEventType,
		layerid: string,
		callback: OnCallbackFunc,
	) => {
		console.log("register map.once", eventname, layerid);
		onces.set(`${eventname}.${layerid}`, {
			callback: callback,
			eventname: eventname,
			layerid: layerid,
		});
		if (map.value && map.value.loaded()) {
			map.value.once(eventname, layerid, callback);
		}
	},
);
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
	/*
	sources.delete(id);
	if (map.value && map.value?.getSource(id)) {
		map.value.removeSource(id);
	}
*/
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
	/*
	layers.delete(id);
	if (map.value?.getLayer(id)) {
		map.value.removeLayer(id);
	}
*/
});

function initializeMap() {
	if (!mapDiv.value) return;

	console.log("initializing map...", props.bounds, props.center, props.zoom);
	const _map = new maplibregl.Map({
		bounds: props.bounds,
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

		ons.forEach((config, id) => {
			console.log("adding map.on", config.eventname, config.layerid);
			_map.on(config.eventname, config.layerid, config.callback);
		});
	});
	onces.forEach((config, id) => {
		console.log("adding map.on", config.eventname, config.layerid);
		_map.once(config.eventname, config.layerid, config.callback);
	});
	map.value = _map;
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
			marker = new maplibregl.Marker({
				color: markerData.color,
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

onMounted(() => {
	initializeMap();
});

onBeforeUnmount(() => {
	if (map.value) {
		map.value.remove();
	}
});
watch(
	() => props.bounds,
	(newBounds) => {
		if (!map.value) return;
		if (map.value.loaded()) {
			map.value.fitBounds(
				newBounds,
				{ padding: 50, duration: 1000 },
				{ isInternalUpdate: true },
			);
		} else {
			map.value.once("load", () => {
				if (!map.value) return;
				map.value.fitBounds(
					newBounds,
					{ padding: 50, duration: 0 },
					{ isInternalUpdate: true },
				);
			});
		}
	},
);
watch(
	() => props.cursor,
	(newCursor) => {
		if (map.value && map.value.loaded()) {
			map.value.getCanvas().style.cursor = newCursor;
		}
	},
);
watch(
	() => props.markers,
	() => {
		updateMarkers();
	},
	{ deep: true },
);
</script>
