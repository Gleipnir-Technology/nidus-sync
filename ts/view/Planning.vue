<style scoped>
/* Add any component-specific styles here */
.pane-header {
	font-weight: 600;
	border-bottom: 2px solid #dee2e6;
}
</style>

<template>
	<ThreeColumn>
		<template #header>
			<div class="col">
				<h3 class="mb-1">Daily Planning Workbench</h3>
				<div class="text-muted small">
					Signals and leads enter from the left, are investigated in the center,
					and transformed into structured field assignments using tools on the
					right.
				</div>
			</div>
		</template>
		<template #left>
			<PlanningColumnList
				:error="error"
				:leads="leads"
				:loading="loading"
				:planFollowups="planFollowups"
				@refresh="refresh"
				:selectedSignalIDs="selectedSignalIDs"
				@signalDeselect="signalDeselect"
				@signalSelect="signalSelect"
				:signals="signal.all"
			/>
		</template>
		<template #center>
			<PlanningColumnDetail
				:markers="markers"
				:selectedSignals="selectedSignals"
				:signals="signal.all"
			/>
		</template>
		<template #right>
			<PlanningColumnAction
				:creating="creating"
				:selectedSignalIDs="selectedSignalIDs"
			/>
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, nextTick } from "vue";

import MapMultipoint from "../components/MapMultipoint.vue";
import PlanningColumnAction from "../components/PlanningColumnAction.vue";
import PlanningColumnDetail from "../components/PlanningColumnDetail.vue";
import PlanningColumnList from "../components/PlanningColumnList.vue";
import ThreeColumn from "../components/layout/ThreeColumn.vue";
import TimeRelative from "../components/TimeRelative.vue";
import { useSignalStore } from "../store/signal";
import { useUserStore } from "../store/user";

// Refs
const mapTile = ref(null);

// State
const apiBase = ref("/api");
const creating = ref(false);
const error = ref(null);
const leads = ref([]);
const loading = ref(false);
const planFollowups = ref([]);
const poolLocations = ref({});
const selectedSignalIDs = ref(new Set<int>([]));
const signal = useSignalStore();
const user = useUserStore();

// Helper functions (outside component)
const getBoundingBox = (points) => {
	if (!points || points.length === 0) {
		return null;
	}

	let minLat = points[0].latitude;
	let maxLat = points[0].latitude;
	let minLng = points[0].longitude;
	let maxLng = points[0].longitude;

	for (const point of points) {
		if (point.latitude < minLat) minLat = point.latitude;
		if (point.latitude > maxLat) maxLat = point.latitude;
		if (point.longitude < minLng) minLng = point.longitude;
		if (point.longitude > maxLng) maxLng = point.longitude;
	}

	return new window.maplibregl.LngLatBounds(
		new window.maplibregl.LngLat(minLng, minLat),
		new window.maplibregl.LngLat(maxLng, maxLat),
	);
};
const markers = computed(() => {
	return [];
});
const selectedSignals = computed(() => {
	if (selectedSignalIDs.value.size == 0 || signal.all == null) {
		return [];
	}
	const result = signal.all.filter((s) => selectedSignalIDs.value.has(s));
	return result;
});
const updateMap = (signals) => {
	const locations = signals.map((s) => s.location);
	const markers = locations.map((l) =>
		new window.maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([l.longitude, l.latitude]),
	);

	/*
	const bounds = getBoundingBox(locations);
	if (bounds != null) {
		mapMultipoint.value.FitBounds(bounds, {
			padding: 50,
		});
	}
*/
};

// Methods
const loadData = async () => {
	loading.value = true;
	error.value = null;

	try {
		await signal.fetchAll();
	} catch (err) {
		error.value = err.message;
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
};

const loadPlanFollowups = async () => {
	try {
		const response = await fetch(`${apiBase.value}/plan-followups`);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		planFollowups.value = data.followups || data;
	} catch (err) {
		console.error("Error loading plan followups:", err);
		throw err;
	}
};

const clearSelection = () => {
	selectedSignalIDs.value.clear();
};

const createLead = async () => {
	if (selectedSignalIDs.value.size === 0) return;

	creating.value = true;

	try {
		const response = await fetch(`${apiBase.value}/leads`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				pool_locations: poolLocations.value,
				signal_ids: selectedSignalIDs,
			}),
		});

		if (!response.ok) {
			const errorData = await response.json();
			throw new Error(
				errorData.message || `HTTP error! status: ${response.status}`,
			);
		}

		const newLead = await response.json();
		leads.value.unshift(newLead);
		clearSelection();
		await loadData();
	} catch (err) {
		console.error("Error creating lead:", err);
		alert(`Failed to create lead: ${err.message}`);
	} finally {
		creating.value = false;
	}
};

const markAsAddressed = async () => {
	if (selectedSignalIDs.value.size === 0) return;

	try {
		const response = await fetch(`${apiBase.value}/signal/mark-addressed`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				signal_ids: selectedSignalIDs,
			}),
		});

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		signals.value = signals.value.filter(
			(signal) => !selectedSignalIDs.value.has(s.id),
		);

		clearSelection();
		alert("Signals marked as addressed");
	} catch (err) {
		console.error("Error marking signals as addressed:", err);
		alert(`Failed to mark signals: ${err.message}`);
	}
};
const refresh = () => {
	loadData();
};
const signalDeselect = (id: int) => {
	selectedSignalIDs.value.delete(id);
};
const signalSelect = (id: int) => {
	selectedSignalIDs.value.add(id);
};
// Lifecycle
onMounted(() => {
	loadData();

	// Configure map multipoint
	/*
	const map = mapMultipoint.value;
	if (map) {
		map.on("load", () => {
			map.addLayer({
				id: "parcel",
				minzoom: 14,
				paint: {
					"line-color": "#0f0",
				},
				source: "tegola",
				"source-layer": "parcel",
				type: "line",
			});
			map.addLayer({
				id: "signal-point",
				paint: {
					"circle-color": "#0D6EfD",
					"circle-radius": 7,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#024AB6",
				},
				source: "tegola",
				"source-layer": "signal-point",
				type: "circle",
			});
			console.log("Added parcel and signal layers");
		});
	}
*/
});
</script>
