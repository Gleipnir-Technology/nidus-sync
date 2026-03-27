<style scoped></style>

<template>
	<ThreeColumn>
		<template #header>
			<div class="col">
				<h3 class="mb-1">Communication Workbench</h3>
				<div class="text-muted small">
					Communications from various sources come in at the left, are
					investigated in the center, and labeled as valuable signal or invalid
					on the right.
				</div>
			</div>
		</template>
		<template #left>
			<CommunicationColumnList
				:all="communication.all"
				@deselect="handleDeselect"
				:loading="loading"
				:selected-id="selectedId"
				@select="handleSelect"
			/>
		</template>
		<template #center>
			<CommunicationColumnDetail
				:loading="loading"
				:mapBounds="mapBounds || undefined"
				:mapMarkers="mapMarkers"
				:selectedCommunication="selectedCommunication"
				:user="user"
				@viewImage="openPhotoViewer"
			/>
		</template>
		<template #right>
			<CommunicationColumnAction
				:loading="loading"
				@markInvalid="markInvalid"
				@markSignal="markSignal"
				@sendMessage="sendMessage"
				:selectedCommunication="selectedCommunication"
				:user="user"
			/>
		</template>
	</ThreeColumn>
	<PhotoViewerModal
		@close="showPhotoModal = false"
		@imageNext="imageNext()"
		@imagePrevious="imagePrevious()"
		:images="currentImages"
		:currentPhotoIndex="currentPhotoIndex"
		:show="showPhotoModal"
	/>
	<ToastNotification
		:message="toastMessage"
		:show="toastShow"
		:title="toastTitle"
	/>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from "vue";
import maplibregl from "maplibre-gl";

import { useCommunicationStore } from "../store/communication";
import { useUserStore } from "../store/user";
import CommunicationColumnAction from "../components/CommunicationColumnAction.vue";
import CommunicationColumnDetail from "../components/CommunicationColumnDetail.vue";
import CommunicationColumnList from "../components/CommunicationColumnList.vue";
import PhotoViewerModal from "../components/PhotoViewerModal.vue";
import ThreeColumn from "../components/layout/ThreeColumn.vue";
import ToastNotification from "../components/ToastNotification.vue";

const communication = useCommunicationStore();
const user = useUserStore();
onMounted(() => {
	fetchCommunications();
});

// Refs
const currentPhotoIndex = ref(0);
const error = ref(null);
const loading = ref(true);
const mapBounds = ref<Bounds | null>(null);
const mapMarkers = ref<Marker[]>([]);
const selectedId = ref<string | null>(null);
const showPhotoModal = ref(false);
const toastMessage = ref("");
const toastShow = ref(false);
const toastTitle = ref("");

const currentPhoto = computed(() => {
	const comm = selectedCommunication.value;
	return comm.public_report?.images[currentPhotoIndex] ?? null;
});
const currentImages = computed(() => {
	const comm = selectedCommunication.value;
	if (comm == null || comm.public_report == null) {
		return [];
	}
	return comm.public_report.images ?? [];
});
const selectedCommunication = computed<Communication | null>(() => {
	if (selectedId.value == null) {
		return null;
	}
	if (communication.all == null) {
		return null;
	}
	const result = communication.all.find((c) => c.id == selectedId.value);
	return result;
});
const handleDeselect = (id: string) => {
	selectedId.value = null;
	updateMap();
};
const handleSelect = (id: string) => {
	selectedId.value = id;
	updateMap();
};
async function fetchCommunications() {
	await communication.fetchAll();
	// if we already had something selected, reset it using the new data
	if (selectedCommunication.value) {
		const matching = communication.all.filter((c) => {
			return c.id === selectedCommunication.value.id;
		});
		if (matching.length > 0) {
			selectedCommunication.value = matching[0];
		}
	}
}
function imageNext() {
	currentPhotoIndex.value = Math.min(
		currentImages.value.length - 1,
		currentPhotoIndex.value + 1,
	);
}
function imagePrevious() {
	currentPhotoIndex.value = Math.max(0, currentPhotoIndex.value - 1);
}
async function loadFromAPI() {
	loading.value = true;
	error.value = null;
	try {
		await Promise.all([fetchCommunications()]);
	} catch (err) {
		error.value = err.message;
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
}

function openPhotoViewer(index) {
	currentPhotoIndex.value = index;
	showPhotoModal.value = true;
}

async function markInvalid() {
	console.log("Marking report as invalid:", selectedCommunication.value.id);
	const payload = {
		reportID: selectedCommunication.value.id,
	};
	const response = await fetch("api/publicreport/invalid", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(payload),
	});

	showNotification(
		"Report Marked Invalid",
		`Report #${selectedCommunication.value.id} has been marked as invalid`,
	);
	removeCurrentFromList();
	await fetchCommunications();
}

async function markSignal() {
	console.log("Marking report as signal:", selectedCommunication.value.id);
	try {
		const report_id = selectedCommunication.value.id;
		const payload = {
			reportID: report_id,
		};
		removeCurrentFromList();
		const response = await fetch("api/publicreport/signal", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(payload),
		});
		if (!response.ok) {
			throw new Error("Failed to submit signal");
		}
		showNotification(
			"Report Marked Signal",
			`Report #${report_id} has been marked as useful signal`,
		);
		await fetchCommunications();
	} catch (err) {
		error.value = err.message;
		console.error("Error creating lead:", err);
	}
}

function removeCurrentFromList() {
	const index = communication.all.findIndex((c) => c.id === selectedId.value);
	if (index > -1) {
		communication.all.splice(index, 1);
	}
	if (communication.all.length > 0) {
		const nextIndex = Math.min(index, communication.all.length - 1);
		selectedId.value = communication.all[nextIndex].id;
	} else {
		selectedId.value = null;
	}
}
async function sendMessage(message: string) {
	if (!message.trim()) return;

	console.log("Sending message reporter:", message.value);

	const payload = {
		message: message.value,
		reportID: selectedCommunication.value.id,
	};
	const response = await fetch(user.urls.api.publicreport_message, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(payload),
	});

	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}

	showNotification(
		"Message Sent",
		`Message successfully sent to ${selectedCommunication.value.public_report.reporter.name}`,
	);
	messageText.value = "";
}
function showNotification(title, message) {
	toastTitle.value = title;
	toastMessage.value = message;
	toastShow.value = true;

	setTimeout(() => {
		toastShow.value = false;
	}, 3000);
}

function updateMap() {
	const loc = selectedCommunication.value?.public_report?.location;
	console.log("updating for loc", loc);
	if (loc == null) {
		mapMarkers.value = [];
		return;
	}

	mapMarkers.value = [
		{
			color: "#FF0000",
			draggable: false,
			id: String(Date.now()),
			location: {
				lng: loc.longitude,
				lat: loc.latitude,
			},
		},
	];
	console.log("markers now", mapMarkers.value);

	let min = { lat: loc.latitude, lng: loc.longitude };
	let max = { lat: loc.latitude, lng: loc.longitude };

	for (const i of selectedCommunication.value.public_report.images) {
		if (
			i.location != null &&
			i.location.latitude != 0 &&
			i.location.longitude != 0
		) {
			mapMarkers.value.push({
				color: "#00FF00",
				draggable: false,
				location: {
					lat: i.location.latitude,
					lng: i.location.longitude,
				},
			});
			min.lat = Math.min(min.lat, i.location.latitude);
			min.lng = Math.min(min.lng, i.location.longitude);
			max.lat = Math.max(max.lat, i.location.latitude);
			max.lng = Math.max(max.lng, i.location.longitude);
		}
	}

	mapBounds.value = {
		max: {
			lat: max.lat + 0.01,
			lng: max.lng + 0.01,
		},
		min: {
			lat: min.lat - 0.01,
			lng: min.lng - 0.01,
		},
	};
}
function onFilterChange(filters) {
	console.log("Filters changed");
}
// Lifecycle hooks
onMounted(async () => {
	await loadFromAPI();

	// Subscribe to SSE events
	if (window.SSEManager) {
		window.SSEManager.subscribe("*", (e) => {
			if (e.resource.startsWith("rmo:")) {
				fetchCommunications();
			}
		});
	}

	// Setup map layer after next tick to ensure map is mounted
	await nextTick();

	/*
	if (mapRef.value) {
		const mapEl = mapRef.value.$el || mapRef.value;
		mapEl.addEventListener("load", () => {
			mapEl.addLayer({
				id: "parcel",
				minzoom: 14,
				paint: {
					"line-color": "#0f0",
				},
				source: "tegola",
				"source-layer": "parcel",
				type: "line",
			});
		});
	}
	*/
});
</script>
