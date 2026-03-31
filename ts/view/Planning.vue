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
				@doAddToLead="doAddToLead"
				@doAddLeadsToAssignment="doAddLeadsToAssignment"
				@doCreateLead="doCreateLead"
				@doCreateProposedAssignment="doCreateProposedAssignment"
				@doEstimateEffort="doEstimateEffort"
				@doMarkSignalAddressed="doMarkSignalAddressed"
				@doSetPriority="doSetPriority"
				@doSendToOperations="doSendToOperations"
				@doSplitLead="doSplitLead"
				:creating="creating"
				:selectedSignalIDs="selectedSignalIDs"
			/>
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, nextTick } from "vue";

import MapMultipoint from "../components/MapMultipoint.vue";
import PlanningColumnAction from "@/components/PlanningColumnAction.vue";
import PlanningColumnDetail from "@/components/PlanningColumnDetail.vue";
import PlanningColumnList from "@/components/PlanningColumnList.vue";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import TimeRelative from "@/components/TimeRelative.vue";
import { useSignalStore } from "@/store/signal";
import { useSessionStore } from "@/store/session";
import { Lead, Location, Point, Signal } from "@/types";

// Refs
const mapTile = ref(null);

// State
const creating = ref(false);
const error = ref<string | null>(null);
const leads = ref<Lead[]>([]);
const loading = ref<boolean>(false);
const planFollowups = ref([]);
const poolLocations = ref({});
const selectedSignalIDs = ref(new Set<number>([]));
const signal = useSignalStore();
const session = useSessionStore();

function doAddToLead() {
	console.log("doAddToLead");
}
function doAddLeadsToAssignment() {
	console.log("doAddLeadsToAssignment");
}
function doCreateLead() {
	console.log("doCreateLead");
}
function doCreateProposedAssignment() {
	console.log("doCreateProposedAssignment");
}
function doEstimateEffort() {
	console.log("doEstimateEffort");
}
function doMarkSignalAddressed() {
	console.log("doMarkSignalAddressed");
}
function doSetPriority() {
	console.log("doSetPriority");
}
function doSendToOperations() {
	console.log("doSendToOperations");
}
function doSplitLead() {
	console.log("doSplitLead");
}
// Helper functions (outside component)
const getBoundingBox = (points: Location[]) => {
	if (!points || points.length === 0) {
		return null;
	}

	let minLat = points[0].lat;
	let maxLat = points[0].lat;
	let minLng = points[0].lng;
	let maxLng = points[0].lng;

	for (const point of points) {
		if (point.lat < minLat) minLat = point.lat;
		if (point.lat > maxLat) maxLat = point.lat;
		if (point.lng < minLng) minLng = point.lng;
		if (point.lng > maxLng) maxLng = point.lng;
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
	const result = signal.all.filter((s: Signal) =>
		selectedSignalIDs.value.has(s.id),
	);
	return result;
});
const updateMap = (signals: Signal[]) => {
	const locations = signals.map((s) => s.location);
	const markers = locations.map((l) =>
		new window.maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([l.lng, l.lat]),
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
		error.value = err instanceof Error ? err.message : "fetch error";
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
};

const loadPlanFollowups = async () => {
	try {
		const response = await fetch("/api/plan-followups");

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
		const response = await fetch("api/leads", {
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
	} finally {
		creating.value = false;
	}
};

const markAsAddressed = async () => {
	if (selectedSignalIDs.value.size === 0) return;

	try {
		const response = await fetch("/api/signal/mark-addressed", {
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

		/*
		signals.value = signals.value.filter(
			(signal) => !selectedSignalIDs.value.has(s.id),
		);
		*/

		clearSelection();
	} catch (err) {
		console.error("Error marking signals as addressed:", err);
	}
};
const refresh = () => {
	loadData();
};
const signalDeselect = (id: number) => {
	selectedSignalIDs.value.delete(id);
};
const signalSelect = (id: number) => {
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
