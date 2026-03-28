<style scoped>
body {
	background-color: #f5f5f5;
}

.left-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	border-right: 1px solid #dee2e6;
}

.middle-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	padding: 20px;
}

.right-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	border-left: 1px solid #dee2e6;
	padding: 20px;
}

.entry-item {
	padding: 15px;
	border-bottom: 1px solid #e9ecef;
	cursor: pointer;
	transition: background-color 0.2s;
}

.entry-item:hover {
	background-color: #f8f9fa;
}

.entry-item.active {
	background-color: #e7f3ff;
	border-left: 4px solid #0d6efd;
}

.placeholder-box {
	background-color: #e9ecef;
	border: 2px dashed #adb5bd;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #6c757d;
	font-size: 18px;
	margin-bottom: 20px;
}

.map-placeholder {
	height: 300px;
}

.image-placeholder {
	height: 400px;
}

.action-btn {
	width: 100%;
	margin-bottom: 10px;
	padding: 12px;
	font-size: 16px;
}

.status-badge {
	font-size: 11px;
}

.map-container {
	margin-bottom: 20px;
}
</style>
<template>
	<ThreeColumn>
		<template #left>
			<ReviewPoolColumnList
				v-if="reviewTask.all"
				:error="error"
				:selectedTaskID="selectedTaskID"
				:tasks="reviewTask.all"
				:total="totalPending"
			/>
			<div v-else>
				<p>Loading</p>
			</div>
		</template>
		<template #center>
			<ReviewPoolColumnDetail :selectedTask="selectedTask" />
		</template>
		<template #right>
			<ReviewPoolColumnAction :submitting="submitting" />
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useReviewTaskStore } from "@/store/review-task";
import { useUserStore } from "@/store/user";
import maplibregl from "maplibre-gl";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ReviewTask from "@/types";
import ReviewPoolColumnAction from "@/components/ReviewPoolColumnAction.vue";
import ReviewPoolColumnDetail from "@/components/ReviewPoolColumnDetail.vue";
import ReviewPoolColumnList from "@/components/ReviewPoolColumnList.vue";

// TypeScript Interfaces
interface Address {
	number: string;
	street: string;
	locality: string;
}

interface Location {
	latitude: number;
	longitude: number;
}

interface Task {
	id: number;
	location: Location;
	condition: string;
	ownerContact?: string;
	residentContact?: string;
	poolShape?: string;
	address?: Address;
}

interface FormData {
	latitude: number;
	longitude: number;
	condition: string;
	ownerContact: string;
	residentContact: string;
	poolShape: string;
}

interface FieldConfig {
	key: keyof FormData;
	label: string;
}

interface Changes {
	updated: string[];
	unchanged: string[];
}

interface MapClickEvent {
	detail: {
		map: any;
		lat: number;
		lng: number;
	};
}

// Props (you can pass these from parent component or environment)
interface Props {
	organizationId?: string;
	tegolaUrl?: string;
	tilesUrl?: string;
	serviceArea?: {
		xmin: number;
		ymin: number;
		xmax: number;
		ymax: number;
	};
}

const props = withDefaults(defineProps<Props>(), {
	organizationId: "",
	tegolaUrl: "",
	tilesUrl: "",
	serviceArea: () => ({
		xmin: 0,
		ymin: 0,
		xmax: 0,
		ymax: 0,
	}),
});

// State
const totalPending = ref<number>(0);
const selectedTaskID = ref<int | null>(null);
const originalValues = ref<Partial<FormData>>({});
const loading = ref<boolean>(true);
const submitting = ref<boolean>(false);
const error = ref<string | null>(null);

const reviewTask = useReviewTaskStore();
const user = useUserStore();

// Refs for map components
const mapMultipoint = ref<any>(null);
const mapTile = ref<any>(null);

// Computed: track which fields have changed
const changes = computed<Changes>(() => {
	if (!selectedTask.value) return { updated: [], unchanged: [] };

	const updated: string[] = [];
	const unchanged: string[] = [];

	const fields: FieldConfig[] = [
		{ key: "latitude", label: "Latitude" },
		{ key: "longitude", label: "Longitude" },
		{ key: "condition", label: "Pool condition" },
		{ key: "ownerContact", label: "Owner contact" },
		{ key: "residentContact", label: "Resident contact" },
		{ key: "poolShape", label: "Pool shape" },
	];

	fields.forEach((field) => {
		if (form[field.key] !== originalValues.value[field.key]) {
			updated.push(field.label);
		} else {
			unchanged.push(field.label);
		}
	});

	return { updated, unchanged };
});

const selectedTask = computed<ReviewTask | null>(() => {
	if (selectedTaskID.value == null) {
		return null;
	}
	return reviewTask.byID(selectedTaskID.value);
});
async function fetchTasks() {
	await reviewTask.fetchAll();
}
// Helper Functions
// Task Selection
function selectTask(task: Task): void {
	console.log("Selected task", task);
	selectedTask.value = task;

	// Populate form with task values
	form.latitude = task.location.latitude;
	form.longitude = task.location.longitude;
	form.condition = task.condition || "";
	form.ownerContact = task.ownerContact || "";
	form.residentContact = task.residentContact || "";
	form.poolShape = task.poolShape || "";

	// Store original values for change tracking
	originalValues.value = { ...form };

	// Update map
	updateMap(task);
}

// Map Update
function updateMap(task: Task): void {
	console.log("Updating map for task:", task.id);

	const map = mapMultipoint.value;
	if (!map) return;

	const loc = task.location;
	const markers = [
		new maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([loc.longitude, loc.latitude]),
	];

	map.SetMarkers(markers);

	const bounds = new maplibregl.LngLatBounds(
		new maplibregl.LngLat(loc.longitude - 0.005, loc.latitude - 0.005),
		new maplibregl.LngLat(loc.longitude + 0.005, loc.latitude + 0.005),
	);

	map.FitBounds(bounds, {
		padding: 50,
	});
}

// Map Click Handler
function updatePoolLocation(event: MapClickEvent): void {
	console.log("map click", selectedTask.value?.id, event.detail);

	const map = event.detail.map;
	const loc = {
		latitude: event.detail.lat,
		longitude: event.detail.lng,
	};

	map.SetMarkers([
		new maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([event.detail.lng, event.detail.lat]),
	]);

	form.latitude = event.detail.lat;
	form.longitude = event.detail.lng;
}

// Submit Review
async function submitReview(action: "committed" | "discarded"): Promise<void> {
	if (!selectedTask.value || submitting.value) return;

	submitting.value = true;
	error.value = null;

	try {
		const payload: any = {
			task_id: selectedTask.value.id,
			status: action,
			updates: {},
		};

		// Include changed fields in the payload
		if (action === "committed") {
			(Object.keys(form) as Array<keyof FormData>).forEach((key) => {
				if (form[key] !== originalValues.value[key]) {
					payload.updates[key] = form[key];
				}
			});
		}

		const response = await fetch("/api/review/pool", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(payload),
		});

		if (!response.ok) {
			throw new Error("Failed to submit review");
		}

		// Remove task from list
		const taskIndex = reviewTask.all.value.findIndex(
			(t) => t.id === selectedTask.value!.id,
		);
		if (taskIndex > -1) {
			reviewTask.all.value.splice(taskIndex, 1);
			totalPending.value = Math.max(0, totalPending.value - 1);
		}

		// Select next task or clear selection
		if (reviewTask.all.length > 0) {
			const nextIndex = Math.min(taskIndex, reviewTask.all.length - 1);
			selectTask(reviewTask.all[nextIndex]);
		} else {
			selectedTask.value = null;
			form.condition = "";
			form.ownerContact = "";
			form.residentContact = "";
			form.poolShape = "";
			form.latitude = 0;
			form.longitude = 0;
			originalValues.value = {};
		}

		// Update list of tasks
		await fetchTasks();
	} catch (err) {
		error.value = err instanceof Error ? err.message : "Unknown error";
		console.error("Error submitting review:", err);
	} finally {
		submitting.value = false;
	}
}

// Action Handlers
function markReviewed(): void {
	submitReview("committed");
}

function discardEntry(): void {
	submitReview("discarded");
}

// Initialize Maps
function initializeMaps(): void {
	const mapElement = mapMultipoint.value;
	const mapTileElement = mapTile.value;

	if (mapElement) {
		mapElement.addEventListener("load", () => {
			mapElement.addLayer({
				id: "parcel",
				minzoom: 14,
				paint: {
					"line-color": "#0f0",
				},
				source: "tegola",
				"source-layer": "parcel",
				type: "line",
			});

			mapElement.addLayer({
				id: "pools",
				paint: {
					"circle-color": "#0D6EfD",
					"circle-radius": 7,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#024AB6",
				},
				source: "tegola",
				"source-layer": "feature-pool",
				type: "circle",
			});

			// Create a popup
			const popup = new maplibregl.Popup({
				closeButton: false,
				closeOnClick: false,
			});

			let currentFeatureCoordinates: string | undefined;

			mapElement.addEventListener("mousemove", "pools", (e: any) => {
				const featureCoordinates =
					e.features[0].geometry.coordinates.toString();
				if (currentFeatureCoordinates !== featureCoordinates) {
					currentFeatureCoordinates = featureCoordinates;
					mapElement.getCanvas().style.cursor = "pointer";

					const coordinates = e.features[0].geometry.coordinates.slice();
					const condition = e.features[0].properties.condition;

					while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
						coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
					}

					popup
						.setLngLat(coordinates)
						.setHTML(condition)
						.addTo(mapElement._map);
				}
			});

			mapElement.addEventListener("mouseleave", "pools", () => {
				currentFeatureCoordinates = undefined;
				mapElement.getCanvas().style.cursor = "";
				popup.remove();
			});
		});
	}

	if (mapTileElement) {
		mapTileElement.addEventListener("load", () => {
			mapTileElement.addLayer({
				id: "parcel",
				minzoom: 14,
				paint: {
					"line-color": "#0f0",
				},
				source: "tegola",
				"source-layer": "parcel",
				type: "line",
			});

			mapTileElement.addLayer({
				id: "pools",
				paint: {
					"circle-color": "#0D6EfD",
					"circle-radius": 7,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#024AB6",
				},
				source: "tegola",
				"source-layer": "feature-pool",
				type: "circle",
			});
		});
	}
}

// Lifecycle
onMounted(async () => {
	initializeMaps();
	await fetchTasks();
});
</script>
