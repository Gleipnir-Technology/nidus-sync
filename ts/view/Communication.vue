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
				@viewImage="openImageViewer"
			/>
		</template>
		<template #right>
			<CommunicationColumnAction
				:loading="loading"
				@markInvalid="markInvalid"
				@markSignal="markSignal"
				@sendMessage="sendMessage"
				:selectedCommunication="selectedCommunication"
			/>
		</template>
	</ThreeColumn>
	<ImageViewerModal
		@close="showImageModal = false"
		@imageNext="imageNext()"
		@imagePrevious="imagePrevious()"
		:images="currentImages"
		:currentImageIndex="currentImageIndex"
		:show="showImageModal"
	/>
	<ToastNotification
		:message="toastMessage"
		:show="toastShow"
		:title="toastTitle"
	/>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import maplibregl from "maplibre-gl";

import CommunicationColumnAction from "@/components/CommunicationColumnAction.vue";
import CommunicationColumnDetail from "@/components/CommunicationColumnDetail.vue";
import CommunicationColumnList from "@/components/CommunicationColumnList.vue";
import ImageViewerModal from "@/components/ImageViewerModal.vue";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ToastNotification from "@/components/ToastNotification.vue";
import { SSEManager } from "@/SSEManager";
import { useCommunicationStore } from "@/store/communication";
import { useSessionStore } from "@/store/session";
import type { Bounds, Marker } from "@/types";
import type { Communication } from "@/type/api";

const communication = useCommunicationStore();
const session = useSessionStore();
onMounted(() => {
	fetchCommunications();
});

// Refs
const currentImageIndex = ref<number>(0);
const error = ref<string | null>(null);
const loading = ref<boolean>(true);
const mapBounds = ref<Bounds | null>(null);
const mapMarkers = ref<Marker[]>([]);
const selectedId = ref<string | null>(null);
const showImageModal = ref(false);
const toastMessage = ref("");
const toastShow = ref(false);
const toastTitle = ref("");

const currentImage = computed(() => {
	const comm = selectedCommunication.value;
	return comm?.public_report?.images[currentImageIndex.value] ?? null;
});
const currentImages = computed(() => {
	const comm = selectedCommunication.value;
	if (comm == null || comm.public_report == null) {
		return [];
	}
	return comm.public_report.images ?? [];
});
const selectedCommunication = computed<Communication | null>(
	(): Communication | null => {
		if (selectedId.value == null) {
			return null;
		}
		if (communication.all == null) {
			return null;
		}
		const result = communication.all.find((c) => c.id == selectedId.value);
		return result || null;
	},
);
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
}
function imageNext() {
	currentImageIndex.value = Math.min(
		currentImages.value.length - 1,
		currentImageIndex.value + 1,
	);
}
function imagePrevious() {
	currentImageIndex.value = Math.max(0, currentImageIndex.value - 1);
}
async function loadFromAPI() {
	loading.value = true;
	error.value = null;
	try {
		await Promise.all([fetchCommunications()]);
	} catch (err) {
		error.value = err instanceof Error ? err.message : "fetch error";
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
}

function openImageViewer(index: number) {
	currentImageIndex.value = index;
	showImageModal.value = true;
}

async function markInvalid() {
	if (selectedCommunication.value == null) {
		return;
	}
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
	if (selectedCommunication.value == null) {
		return;
	}
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
		error.value = err instanceof Error ? err.message : "fetch error";
		console.error("Error creating lead:", err);
	}
}

function removeCurrentFromList() {
	if (communication.all == null) {
		return;
	}
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
	if (selectedCommunication.value == null) return;
	if (session.urls == null) return;
	console.log("Sending message reporter:", message);

	const payload = {
		message: message,
		reportID: selectedCommunication.value.id,
	};
	const response = await fetch(session.urls?.api.publicreport_message, {
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
		`Message successfully sent to ${selectedCommunication.value.public_report?.reporter.name}`,
	);
}
function showNotification(title: string, message: string) {
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
			location: loc,
		},
	];
	console.log("markers now", mapMarkers.value);

	let min = loc;
	let max = loc;

	for (const i of selectedCommunication.value?.public_report?.images ?? []) {
		if (
			i.location != null &&
			i.location.latitude != 0 &&
			i.location.longitude != 0
		) {
			mapMarkers.value.push({
				color: "#00FF00",
				draggable: false,
				id: new Date().toISOString(),
				location: i.location,
			});
			min.latitude = Math.min(min.latitude, i.location.latitude);
			min.longitude = Math.min(min.longitude, i.location.longitude);
			max.latitude = Math.max(max.latitude, i.location.latitude);
			max.longitude = Math.max(max.longitude, i.location.longitude);
		}
	}

	mapBounds.value = {
		max: {
			latitude: max.latitude + 0.01,
			longitude: max.longitude + 0.01,
		},
		min: {
			latitude: min.latitude - 0.01,
			longitude: min.longitude - 0.01,
		},
	};
}
// Lifecycle hooks
onMounted(async () => {
	await loadFromAPI();

	// Setup map layer after next tick to ensure map is mounted
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
