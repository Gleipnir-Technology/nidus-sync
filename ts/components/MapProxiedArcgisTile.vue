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
import type { LngLatBoundsLike, Map as MapLibreMap } from "maplibre-gl";
import maplibregl from "maplibre-gl";
import {
	ref,
	watch,
	onMounted,
	onBeforeUnmount,
	shallowRef,
	type Ref,
} from "vue";
import { Location, MapClickEvent, Marker, Point } from "@/types";

interface Emits {
	(e: "map-click", event: MapClickEvent): void;
}
interface Props {
	location: Location;
	markers: Marker[];
	organizationId: Number;
	tegola: string;
	urlTiles: string;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();

const error = ref<string | null>(null);
const mapContainer = ref<HTMLElement | null>(null);
const map: Ref<MapLibreMap | null> = shallowRef(null);
const markerInstances = ref<Map<string, maplibregl.Marker>>(new Map());
const markers = ref<Map<string, maplibregl.Marker>>(new Map());

// Watch for latitude/longitude changes
watch(
	() => [props.location],
	([newLocation]) => {
		if (map.value) {
			map.value.jumpTo({
				center: [newLocation.lng, newLocation.lat],
				zoom: 19,
			});
		}
	},
);

const initializeMap = () => {
	if (!mapContainer.value) return;

	try {
		map.value = new maplibregl.Map({
			center: [props.location.lng, props.location.lat],
			container: mapContainer.value,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
			zoom: 19,
		});
		const mapInstance = map.value;

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

			mapInstance.addSource("flyover", {
				type: "raster",
				tiles: [props.urlTiles],
			});

			mapInstance.addLayer({
				id: "flyover-layer",
				source: "flyover",
				type: "raster",
			});

			mapInstance.on("click", (e) => {
				emit("map-click", {
					location: {
						lat: e.lngLat.lat,
						lng: e.lngLat.lng,
					},
					map: mapInstance,
					point: e.point,
				});
			});
		});
	} catch (e) {
		console.error("hey dummy", e);
	}
};

onMounted(() => {
	setTimeout(() => initializeMap(), 0);
});

onBeforeUnmount(() => {
	if (map.value) {
		map.value.remove();
	}
});
</script>
