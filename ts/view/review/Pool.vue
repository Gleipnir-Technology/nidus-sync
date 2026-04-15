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
				v-show="storeReviewTask.all"
				@doSelectTask="selectTask"
				:error="error"
				:selectedTaskID="selectedTaskID"
				:tasks="storeReviewTask.all()"
				:total="totalPending"
			/>
			<div v-show="!storeReviewTask.all">
				<p>Loading</p>
			</div>
		</template>
		<template #center>
			<ReviewPoolColumnDetail
				:loading="loading"
				:mapBounds="mapBounds || undefined"
				:mapFlyoverCamera="mapFlyoverCamera"
				:mapMarkers="mapMarkers"
				:selectedTask="selectedTask"
				v-model="reviewForm"
			/>
		</template>
		<template #right>
			<ReviewPoolColumnAction
				:changes="changes"
				@doComplete="doComplete"
				@doDiscard="doDiscard"
				:selectedTask="selectedTask"
				:submitting="submitting"
			/>
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useStoreReviewTask } from "@/store/review-task";
import { useSessionStore } from "@/store/session";
import maplibregl from "maplibre-gl";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ReviewPoolColumnAction from "@/components/ReviewPoolColumnAction.vue";
import ReviewPoolColumnDetail, {
	ReviewTaskPoolForm,
} from "@/components/ReviewPoolColumnDetail.vue";
import ReviewPoolColumnList from "@/components/ReviewPoolColumnList.vue";
import { formatAddress } from "@/format";
import type { Changes } from "@/types";
import { Bounds, Contact, Location, ReviewTask } from "@/type/api";
import { Camera } from "@/type/map";
import { MapClickEvent, Marker } from "@/types";

interface FieldConfig {
	key: keyof ReviewTaskPoolForm;
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
const error = ref<string | null>(null);
const loading = ref<boolean>(true);
const mapBounds = ref<Bounds | null>(null);
const mapFlyoverCamera = ref<Camera>(new Camera());
const poolLocation = ref<Location>({
	latitude: 0,
	longitude: 0,
});
const reviewForm = ref<ReviewTaskPoolForm>({
	address: "",
	condition: "",
	location: new Location(),
	owner: "",
	resident: "",
});
const selectedTaskID = ref<number | null>(null);
const session = useSessionStore();
const storeReviewTask = useStoreReviewTask();
const submitting = ref<boolean>(false);
const totalPending = ref<number>(0);

// Refs for map components
const mapMultipoint = ref<any>(null);
const mapTile = ref<any>(null);

// Computed: track which fields have changed
const changes = computed<Changes>(() => {
	const task = selectedTask.value;
	const pool = task?.pool;
	if (!pool) return { updated: [], unchanged: [] };

	const updated: string[] = [];
	const unchanged: string[] = [];

	const fields: FieldConfig[] = [
		{ key: "location", label: "Location" },
		{ key: "condition", label: "Pool condition" },
		{ key: "owner", label: "Owner contact" },
		{ key: "resident", label: "Resident contact" },
	];

	if (reviewForm.value.condition != pool.condition) {
		updated.push("condition");
	} else {
		unchanged.push("condition");
	}
	if (
		reviewForm.value.location.latitude != pool.location?.latitude ||
		reviewForm.value.location.longitude != pool.location?.longitude
	) {
		updated.push("location");
	} else {
		unchanged.push("location");
	}
	if (reviewForm.value.owner != (pool.site.owner?.name ?? "")) {
		updated.push("owner");
	} else {
		unchanged.push("owner");
	}
	if (reviewForm.value.resident != (pool.site.resident?.name ?? "")) {
		updated.push("resident");
	} else {
		unchanged.push("resident");
	}

	return { updated, unchanged };
});
const mapMarkers = computed<Marker[]>(() => {
	const form = reviewForm.value;
	const task = selectedTask.value;
	const loc =
		reviewForm.value.location.latitude != 0
			? reviewForm.value.location
			: task?.pool?.location;
	if (!loc) {
		return [];
	}
	const markers = {
		color: "#FF0000",
		draggable: false,
		id: "x",
		location: loc,
	};
	return [markers];
});
const selectedTask = computed<ReviewTask | undefined>(() => {
	if (selectedTaskID.value == null) {
		return undefined;
	}
	return storeReviewTask.byID(selectedTaskID.value);
});
// Helper Functions
// Task Selection
function selectTask(id: number): void {
	selectedTaskID.value = id;

	const task = storeReviewTask.byID(id);
	if (!task) {
		console.log("no task", id);
		return;
	}
	const pool = task.pool;
	if (!pool) {
		console.log("no pool for selected task");
		return;
	}
	console.log("selecting task", id, task);
	mapFlyoverCamera.value = new Camera(pool.location, 20);
	reviewForm.value = {
		address: formatAddress(task.address),
		condition: pool.condition,
		location: pool.location,
		owner: pool.site.owner?.name ?? "",
		resident: pool.site.resident?.name ?? "",
	};
}

// Submit Review
async function submitReview(action: "committed" | "discarded"): Promise<void> {
	if (!selectedTask.value || submitting.value) return;

	submitting.value = true;
	error.value = null;

	try {
		const payload: any = {
			status: action,
			task_id: selectedTask.value.id,
			updates: {
				condition: reviewForm.value.condition,
				latitude: reviewForm.value.location.latitude,
				longitude: reviewForm.value.location.longitude,
			},
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
		// Save the current item's index for setting the newly selected item
		const index = storeReviewTask
			.all()
			.findIndex((t) => t.id == selectedTaskID.value);

		// Remove task from list
		storeReviewTask.remove(selectedTask.value!.id);

		// Update list of tasks
		const all_tasks: ReviewTask[] = await storeReviewTask.fetchAll();

		// Select the next item in the list
		let new_index = index < all_tasks.length ? index : all_tasks.length - 1;
		const new_id = all_tasks[new_index].id;
		selectTask(new_id);
	} catch (err) {
		error.value = err instanceof Error ? err.message : "Unknown error";
		console.error("Error submitting review:", err);
	} finally {
		submitting.value = false;
	}
}

// Action Handlers
function doComplete(): void {
	submitReview("committed");
}

function doDiscard(): void {
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

watch(
	() => reviewForm.value,
	(newReviewForm: ReviewTaskPoolForm) => {
		console.log("new review form", newReviewForm);
	},
);
// Lifecycle
onMounted(async () => {
	initializeMaps();

	loading.value = true;
	await storeReviewTask.fetchAll();
	loading.value = false;
});
</script>
