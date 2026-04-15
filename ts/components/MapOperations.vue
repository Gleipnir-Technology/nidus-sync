<style scoped lang="scss">
@import url("https://unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
.map {
	width: 100%;
	height: 100%;
	transition: filter 0.2s ease;
}
.map-operations {
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
}

.map-header {
	padding: 1rem;
	background-color: #f8f9fa;
	border-bottom: 1px solid #dee2e6;
}

.map-header h2 {
	margin: 0 0 1rem 0;
	font-size: 1.5rem;
	color: #333;
}

.legend {
	display: flex;
	gap: 1.5rem;
	flex-wrap: wrap;
}

.legend-item {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	font-size: 0.875rem;
}

.legend-color {
	width: 20px;
	height: 4px;
	border-radius: 2px;
}

.map-container {
	flex: 1;
	min-height: 500px;
	position: relative;
	width: 100%;
	height: 100%;
	border-radius: 10px;
	overflow: hidden;
}
/* Ensure map fills container on all devices */
:deep(.maplibregl-map) {
	width: 100%;
	height: 100%;
}

:deep(.maplibregl-canvas) {
	width: 100%;
	height: 100%;
}
</style>
<template>
	<div class="map-operations">
		<div class="map-header">
			<h2>Technician Routes</h2>
			<div class="legend">
				<div class="legend-item">
					<span class="legend-color" style="background-color: #3b82f6"></span>
					<span>Technician A (5 stops)</span>
				</div>
				<div class="legend-item">
					<span class="legend-color" style="background-color: #ef4444"></span>
					<span>Technician B (4 stops)</span>
				</div>
				<div class="legend-item">
					<span class="legend-color" style="background-color: #10b981"></span>
					<span>Technician C (6 stops)</span>
				</div>
			</div>
		</div>
		<div ref="mapContainer" class="map-container"></div>
	</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

const mapContainer = ref(null);
let map = null;

// Mock route data for three technicians
const routes = [
	{
		id: "route-a",
		name: "Technician A",
		color: "#3b82f6",
		stops: [
			[-122.4194, 37.7749], // San Francisco
			[-122.4284, 37.7849],
			[-122.4374, 37.7949],
			[-122.4264, 37.8049],
			[-122.4154, 37.7949],
		],
	},
	{
		id: "route-b",
		name: "Technician B",
		color: "#ef4444",
		stops: [
			[-122.4094, 37.7649],
			[-122.4184, 37.7549],
			[-122.4274, 37.7649],
			[-122.4184, 37.7749],
		],
	},
	{
		id: "route-c",
		name: "Technician C",
		color: "#10b981",
		stops: [
			[-122.4494, 37.7849],
			[-122.4584, 37.7949],
			[-122.4494, 37.8049],
			[-122.4394, 37.8149],
			[-122.4294, 37.8249],
			[-122.4194, 37.8149],
		],
	},
];

const initMap = () => {
	console.log("init ops map");
	map = new maplibregl.Map({
		center: [-122.4194, 37.7849],
		container: mapContainer.value,
		style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
		zoom: 12,
	});

	map.addControl(new maplibregl.NavigationControl(), "top-left");
	map.on("load", () => {
		console.log("ops map loaded");
		// Add routes and stops for each technician
		routes.forEach((route, routeIndex) => {
			// Add route line
			map.addSource(route.id, {
				type: "geojson",
				data: {
					type: "Feature",
					properties: {},
					geometry: {
						type: "LineString",
						coordinates: route.stops,
					},
				},
			});

			map.addLayer({
				id: route.id,
				type: "line",
				source: route.id,
				layout: {
					"line-join": "round",
					"line-cap": "round",
				},
				paint: {
					"line-color": route.color,
					"line-width": 4,
					"line-opacity": 0.8,
				},
			});

			// Add stops as circles
			map.addSource(`${route.id}-stops`, {
				type: "geojson",
				data: {
					type: "FeatureCollection",
					features: route.stops.map((coord, index) => ({
						type: "Feature",
						properties: {
							title: `${route.name} - Stop ${index + 1}`,
							stopNumber: index + 1,
							isStart: index === 0,
							isEnd: index === route.stops.length - 1,
						},
						geometry: {
							type: "Point",
							coordinates: coord,
						},
					})),
				},
			});

			// Add circle layer for stops
			map.addLayer({
				id: `${route.id}-stops`,
				type: "circle",
				source: `${route.id}-stops`,
				paint: {
					"circle-radius": [
						"case",
						["get", "isStart"],
						8,
						["get", "isEnd"],
						8,
						6,
					],
					"circle-color": route.color,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#ffffff",
				},
			});

			// Add labels for stop numbers
			map.addLayer({
				id: `${route.id}-labels`,
				type: "symbol",
				source: `${route.id}-stops`,
				layout: {
					"text-field": ["get", "stopNumber"],
					"text-size": 10,
					"text-offset": [0, 0],
					"text-anchor": "center",
				},
				paint: {
					"text-color": "#ffffff",
				},
			});

			// Add popups on click
			map.on("click", `${route.id}-stops`, (e) => {
				const coordinates = e.features[0].geometry.coordinates.slice();
				const description = e.features[0].properties.title;

				new maplibregl.Popup()
					.setLngLat(coordinates)
					.setHTML(`<strong>${description}</strong>`)
					.addTo(map);
			});

			// Change cursor on hover
			map.on("mouseenter", `${route.id}-stops`, () => {
				map.getCanvas().style.cursor = "pointer";
			});

			map.on("mouseleave", `${route.id}-stops`, () => {
				map.getCanvas().style.cursor = "";
			});
		});
	});
};

onMounted(() => {
	setTimeout(() => {
		initMap();
	}, 0);
});

onUnmounted(() => {
	if (map) {
		map.remove();
	}
});
</script>
