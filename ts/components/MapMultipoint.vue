<template>
	<div ref="mapContainer" class="map-multipoint"></div>
</template>

<script setup>
import {
	ref,
	onMounted,
	onUnmounted,
	defineProps,
	defineEmits,
	defineExpose,
} from "vue";
import maplibregl from "maplibre-gl";

const props = defineProps({
	xmin: {
		type: Number,
		default: 0,
	},
	ymin: {
		type: Number,
		default: 0,
	},
	xmax: {
		type: Number,
		default: 0,
	},
	ymax: {
		type: Number,
		default: 0,
	},
	organizationId: {
		type: Number,
		default: 0,
	},
	tegola: {
		type: String,
		default: "",
	},
});

const emit = defineEmits(["load"]);

const mapContainer = ref(null);
const _map = ref(null);
const _markers = ref([]);
const _preOns = ref([]);

const _bounds = () => {
	let bounds = [
		[props.xmin, props.ymin],
		[props.xmax, props.ymax],
	];

	if (
		props.xmin === 0 ||
		props.xmax === 0 ||
		props.ymin === 0 ||
		props.ymax === 0
	) {
		bounds = [
			[-125, 25],
			[-70, 50],
		];
	}

	return bounds;
};

const _initializeMap = () => {
	const bounds = _bounds();

	_map.value = new maplibregl.Map({
		bounds: bounds,
		container: mapContainer.value,
		style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
	});

	_map.value.on("load", () => {
		if (props.organizationId !== 0) {
			_map.value.addSource("tegola", {
				type: "vector",
				tiles: [
					`${props.tegola}maps/nidus/{z}/{x}/{y}?id=${props.organizationId}&organization_id=${props.organizationId}`,
				],
			});

			_map.value.addLayer({
				id: "service-area",
				source: "tegola",
				"source-layer": "service-area-bounds",
				type: "line",
				paint: {
					"line-color": "#f00",
				},
			});
		}

		emit("load", { map: _map.value });
	});

	for (const on of _preOns.value) {
		_map.value.on(...on);
	}
};

// Map wrapper methods
const addLayer = (a) => {
	return _map.value?.addLayer(a);
};

const addSource = (a, b) => {
	return _map.value?.addSource(a, b);
};

const flyTo = (a, b) => {
	return _map.value?.flyTo(a, b);
};

const getCanvas = (...args) => {
	return _map.value?.getCanvas(...args);
};

const getContainer = (...args) => {
	return _map.value?.getContainer(...args);
};

const jumpTo = (args) => {
	return _map.value?.jumpTo(args);
};

const on = (...args) => {
	if (_map.value != null) {
		return _map.value.on(...args);
	} else {
		_preOns.value.push(args);
	}
};

const once = (a, b) => {
	return _map.value?.once(a, b);
};

const panTo = (a, b) => {
	return _map.value?.panTo(a, b);
};

const queryRenderedFeatures = (a) => {
	return _map.value?.queryRenderedFeatures(a);
};

const ClearMarkers = () => {
	_markers.value.forEach((marker) => marker.remove());
};

const FitBounds = (bounds, options) => {
	return _map.value?.fitBounds(bounds, options);
};

const ResetCamera = () => {
	const bounds = _bounds();
	FitBounds(bounds, {
		linear: false,
	});
};

const SetLayoutProperty = (layout, property, value) => {
	return _map.value?.setLayoutProperty(layout, property, value);
};

const SetMarkers = (markers) => {
	console.log("Setting map markers", markers);
	_markers.value.forEach((marker) => marker.remove());
	_markers.value = markers;
	for (let m of markers) {
		m.addTo(_map.value);
	}
};

// Lifecycle
onMounted(() => {
	setTimeout(() => _initializeMap(), 0);
});

onUnmounted(() => {
	if (_map.value) {
		_map.value.remove();
	}
});

// Expose methods to parent component
defineExpose({
	addLayer,
	addSource,
	flyTo,
	getCanvas,
	getContainer,
	jumpTo,
	on,
	once,
	panTo,
	queryRenderedFeatures,
	ClearMarkers,
	FitBounds,
	ResetCamera,
	SetLayoutProperty,
	SetMarkers,
});
</script>

<style scoped>
@import url("//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");

.map-multipoint {
	height: 100%;
	width: 100%;
}
</style>
