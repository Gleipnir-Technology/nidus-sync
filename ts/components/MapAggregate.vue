<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

interface Props {
	centroid: [number, number];
	organizationId: number;
	tegola: string;
	xmin: number;
	ymin: number;
	xmax: number;
	ymax: number;
}

interface CellClickDetail {
	cell: string;
}

const props = withDefaults(defineProps<Props>(), {
	organizationId: 0,
});

const emit = defineEmits<{
	"cell-click": [detail: CellClickDetail];
}>();

const mapContainer = ref<HTMLElement | null>(null);
const map = ref<maplibregl.Map | null>(null);

const initializeMap = () => {
	if (!mapContainer.value) return;

	const bounds: [[number, number], [number, number]] = [
		[props.xmin, props.ymin],
		[props.xmax, props.ymax],
	];

	map.value = new maplibregl.Map({
		bounds: bounds,
		container: mapContainer.value,
		style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
	});

	console.log("Initializing map to bounds", bounds);

	map.value.on("load", () => {
		if (!map.value) return;

		map.value.addSource("tegola", {
			type: "vector",
			tiles: [
				`${props.tegola}maps/nidus/{z}/{x}/{y}?id=${props.organizationId}&organization_id=${props.organizationId}`,
			],
		});

		map.value.addLayer({
			id: "mosquito_source",
			type: "fill",
			filter: ["==", ["zoom"], ["+", 2, ["to-number", ["get", "resolution"]]]],
			source: "tegola",
			"source-layer": "mosquito_source",
			paint: {
				"fill-opacity": 0.4,
				"fill-color": "#dc3545",
			},
		});

		map.value.addLayer({
			id: "service_request",
			type: "fill",
			filter: ["==", ["zoom"], ["+", 2, ["to-number", ["get", "resolution"]]]],
			source: "tegola",
			"source-layer": "service_request",
			paint: {
				"fill-opacity": 0.4,
				"fill-color": "#ffc107",
			},
		});

		map.value.addLayer({
			id: "trap",
			type: "fill",
			filter: ["==", ["zoom"], ["+", 2, ["to-number", ["get", "resolution"]]]],
			source: "tegola",
			"source-layer": "trap",
			paint: {
				"fill-opacity": 0.4,
				"fill-color": "#0dcaf0",
			},
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

		map.value.on("mouseenter", "mosquito_source", () => {
			if (map.value) {
				map.value.getCanvas().style.cursor = "pointer";
			}
		});

		map.value.on("mouseleave", "mosquito_source", () => {
			if (map.value) {
				map.value.getCanvas().style.cursor = "";
			}
		});

		const handleClick = (e: maplibregl.MapLayerMouseEvent) => {
			if (!e.features || e.features.length === 0) return;

			const feature = e.features[0];
			const properties = feature.properties;

			emit("cell-click", {
				cell: properties.cell,
			});
		};

		map.value.on("click", "mosquito_source", handleClick);
		map.value.on("click", "service_request", handleClick);
		map.value.on("click", "trap", handleClick);
	});
};

const jumpTo = (args: maplibregl.JumpToOptions) => {
	if (map.value) {
		map.value.jumpTo(args);
	}
};

onMounted(() => {
	setTimeout(() => initializeMap(), 0);
});

onUnmounted(() => {
	if (map.value) {
		map.value.remove();
		map.value = null;
	}
});

// Expose public methods
defineExpose({
	jumpTo,
});
</script>

<template>
	<div class="map-container">
		<div ref="mapContainer" class="map"></div>
	</div>
</template>

<style scoped>
.map-container {
	background-color: #e9ecef;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	margin-top: 20px;
	position: relative;
}

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
