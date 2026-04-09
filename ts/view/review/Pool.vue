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
				@doSelectTask="selectTask"
				:error="error"
				:selectedTaskID="selectedTaskID"
				:tasks="reviewTask.all()"
				:total="totalPending"
			/>
			<div v-else>
				<p>Loading</p>
			</div>
		</template>
		<template #center>
			<ReviewPoolColumnDetail
				:loading="loading"
				:mapBounds="mapBounds || undefined"
				:mapMarkers="mapMarkers"
				:newPoolCondition="newPoolCondition"
				:newPoolLocation="newPoolLocation"
				:selectedTask="selectedTask"
			/>
		</template>
		<template #right>
			<ReviewPoolColumnAction :changes="changes" :submitting="submitting" />
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useReviewTaskStore } from "@/store/review-task";
import { useSessionStore } from "@/store/session";
import maplibregl from "maplibre-gl";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ReviewPoolColumnAction from "@/components/ReviewPoolColumnAction.vue";
import ReviewPoolColumnDetail from "@/components/ReviewPoolColumnDetail.vue";
import ReviewPoolColumnList from "@/components/ReviewPoolColumnList.vue";
import type { Changes } from "@/types";
import { Contact, Location, ReviewTask } from "@/type/api";
import { Bounds, MapClickEvent, Marker } from "@/types";

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

// Props (you can pass these from parent component or environment)
interface Props {
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
const newPoolCondition = ref<string>("");
const newPoolLocation = ref<Location>({ latitude: 0, longitude: 0 });
const newOwnerName = ref<string>("");
const newResidentName = ref<string>("");
const error = ref<string | null>(null);
const loading = ref<boolean>(true);
const mapBounds = ref<Bounds | null>(null);
const mapMarkers = ref<Marker[]>([]);
const selectedTaskID = ref<number | null>(null);
const submitting = ref<boolean>(false);
const totalPending = ref<number>(0);

const reviewTask = useReviewTaskStore();
const session = useSessionStore();

// Refs for map components
const mapMultipoint = ref<any>(null);
const mapTile = ref<any>(null);

// Computed: track which fields have changed
const changes = computed<Changes>(() => {
	const pool = selectedTask.value?.pool;
	if (!pool) return { updated: [], unchanged: [] };

	const updated: string[] = [];
	const unchanged: string[] = [];

	const fields: FieldConfig[] = [
		{ key: "latitude", label: "Latitude" },
		{ key: "longitude", label: "Longitude" },
		{ key: "condition", label: "Pool condition" },
		{ key: "ownerContact", label: "Owner contact" },
		{ key: "residentContact", label: "Resident contact" },
	];

	if (newPoolCondition.value != pool.condition) {
		updated.push("condition");
	} else {
		unchanged.push("condition");
	}
	if (newPoolLocation.value.latitude != pool.site.location.latitude) {
		updated.push("latitude");
	} else {
		unchanged.push("latitude");
	}
	if (newPoolLocation.value.longitude != pool.site.location.longitude) {
		updated.push("longitude");
	} else {
		unchanged.push("longitude");
	}
	if (newOwnerName.value != pool.site.owner?.name) {
		updated.push("ownerContact");
	} else {
		unchanged.push("ownerContact");
	}
	if (newResidentName.value != pool.site.resident?.name) {
		updated.push("residentContact");
	} else {
		unchanged.push("residentContact");
	}

	return { updated, unchanged };
});

const selectedTask = computed<ReviewTask | undefined>(() => {
	if (selectedTaskID.value == null) {
		return undefined;
	}
	return reviewTask.byID(selectedTaskID.value);
});
async function fetchTasks() {
	await reviewTask.fetchAll();
}
// Helper Functions
// Task Selection
function selectTask(id: number): void {
	console.log("Selected task", id);
	selectedTaskID.value = id;

	const task = reviewTask.byID(id);
	if (!task) {
		console.log("no task", id);
		return;
	}
	const pool = task.pool;
	if (!pool) {
		console.log("no pool for selected task");
		return;
	}
	newPoolCondition.value = pool.condition;
	newPoolLocation.value = pool.location;
	newOwnerName.value = pool.site.owner?.name ?? "";
	newResidentName.value = pool.site.resident?.name ?? "";

	// Update map
	updateMap(task);
}

// Map Update
function updateMap(task: ReviewTask): void {
	console.log("Updating map for task:", task.id);

	const map = mapMultipoint.value;
	if (!map) return;

	const loc = task.pool?.location;
	if (!loc) {
		map.SetMarkers([]);
		return;
	}
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
	console.log("map click", selectedTask.value?.id, event);

	/*
	const map = event.map;
	const loc = event.location;

	map.SetMarkers([
		new maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([loc.lng, loc.lat]),
	]);

	selectedTaskChanges.latitude = event.detail.lat;
	selectedTaskChanges.longitude = event.detail.lng;
	*/
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
			//updates: selectedTaskChanges,
		};

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
		reviewTask.remove(selectedTask.value!.id);

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
