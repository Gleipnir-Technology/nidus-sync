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
				:all="storeCommunication.all"
				@deselect="handleDeselect"
				:loading="storeCommunication.loading"
				:selected-id="selectedId"
				@select="handleSelect"
			/>
		</template>
		<template #center>
			<CommunicationColumnDetail
				:loading="storePublicReport.loading || storeCommunication.loading"
				:mapBounds="mapBounds || undefined"
				:mapMarkers="mapMarkers"
				:selectedCommunication="selectedCommunication"
				:selectedReport="selectedReport"
				@viewImage="openImageViewer"
			/>
		</template>
		<template #right>
			<CommunicationColumnAction
				:loading="storePublicReport.loading || storeCommunication.loading"
				@markInvalid="markInvalid"
				@markSignal="markSignal"
				@sendMessage="sendMessage"
				:selectedCommunication="selectedCommunication"
				:selectedReport="selectedReport"
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
import { computed, onMounted, ref, watch } from "vue";
import { computedAsync } from "@vueuse/core";
import maplibregl from "maplibre-gl";

import CommunicationColumnAction from "@/components/CommunicationColumnAction.vue";
import CommunicationColumnDetail from "@/components/CommunicationColumnDetail.vue";
import CommunicationColumnList from "@/components/CommunicationColumnList.vue";
import ImageViewerModal from "@/components/ImageViewerModal.vue";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ToastNotification from "@/components/ToastNotification.vue";
import { useQueryParam } from "@/composable/use-query-param";
import { SSEManager } from "@/SSEManager";
import { useCommunicationStore } from "@/store/communication";
import { useSessionStore } from "@/store/session";
import type { Marker } from "@/types";
import { Bounds, type Communication, PublicReport } from "@/type/api";
import type { LngLatBounds } from "@/map/Map.vue";
import { boundsForServiceArea, boundsWithPadding } from "@/map/util";
import { useStorePublicReport } from "@/store/publicreport";

const session = useSessionStore();

// Refs
const currentImageIndex = ref<number>(0);
const paramCommunication = useQueryParam("communication");
const selectedId = ref<string | null>(null);
const showImageModal = ref(false);
const storeCommunication = useCommunicationStore();
const storePublicReport = useStorePublicReport();
const toastMessage = ref("");
const toastShow = ref(false);
const toastTitle = ref("");

const currentImage = computed(() => {
	const comm = selectedCommunication.value;
	return selectedReport.value?.images[currentImageIndex.value] ?? null;
});
const currentImages = computed(() => {
	const comm = selectedCommunication.value;
	if (comm == null || comm.public_report == null) {
		return [];
	}
	return selectedReport.value?.images ?? [];
});
const mapBounds = computed<LngLatBounds | null>((): LngLatBounds | null => {
	let bounds = new Bounds();
	const loc = selectedReport.value?.location;
	console.log("updating for loc", loc);
	if (loc && loc.latitude != 0 && loc.longitude != 0) {
		bounds.addLocation(loc);
	}
	const address_loc = selectedReport.value?.address.location;
	if (address_loc && address_loc.latitude != 0 && address_loc.longitude != 0) {
		bounds.addLocation(address_loc);
	}

	for (const [i, image] of (selectedReport.value?.images ?? []).entries()) {
		if (
			image.location != null &&
			image.location.latitude != 0 &&
			image.location.longitude != 0
		) {
			bounds.addLocation(image.location);
		}
	}
	if (bounds.isEmpty()) {
		return boundsForServiceArea();
	}
	return boundsWithPadding(bounds.min, bounds.max, 0.01);
});
const mapMarkers = computed<Marker[]>((): Marker[] => {
	const loc = selectedReport.value?.location;
	let markers: Marker[] = [];
	if (loc && loc.latitude != 0 && loc.longitude != 0) {
		markers.push({
			color: "#0000FF",
			draggable: false,
			id: "reporter",
			location: loc,
		});
	}
	const address_loc = selectedReport.value?.address.location;
	if (address_loc && address_loc.latitude != 0 && address_loc.longitude != 0) {
		markers.push({
			color: "#FF0000",
			draggable: false,
			id: "address",
			location: address_loc,
		});
	}

	for (const [i, image] of (selectedReport.value?.images ?? []).entries()) {
		if (
			image.location != null &&
			image.location.latitude != 0 &&
			image.location.longitude != 0
		) {
			markers.push({
				color: "#00FF00",
				draggable: false,
				id: `image-${i}`,
				location: image.location,
			});
		}
	}
	return markers;
});
const selectedCommunication = computed<Communication | null>(
	(): Communication | null => {
		if (selectedId.value == null) {
			return null;
		}
		if (storeCommunication.all == null) {
			return null;
		}
		const result = storeCommunication.all.find((c) => c.id == selectedId.value);
		return result || null;
	},
);
const selectedReport = computedAsync(
	async (): Promise<PublicReport | undefined> => {
		if (
			!(
				selectedCommunication.value && selectedCommunication.value.public_report
			)
		)
			return;
		return await storePublicReport.fetchByURI(
			selectedCommunication.value.public_report,
		);
	},
);
const handleDeselect = (id: string) => {
	selectedId.value = null;
};
const handleSelect = (id: string) => {
	selectedId.value = id;
	paramCommunication.setValue(id);
};
function imageNext() {
	currentImageIndex.value = Math.min(
		currentImages.value.length - 1,
		currentImageIndex.value + 1,
	);
}
function imagePrevious() {
	currentImageIndex.value = Math.max(0, currentImageIndex.value - 1);
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
	const response = await fetch("/api/publicreport/invalid", {
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
	await storeCommunication.fetchAll();
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
		await storeCommunication.fetchAll();
	} catch (err) {
		console.error("Error creating lead:", err);
	}
}

function removeCurrentFromList() {
	if (storeCommunication.all == null) {
		return;
	}
	const index = storeCommunication.all.findIndex(
		(c) => c.id === selectedId.value,
	);
	if (index > -1) {
		storeCommunication.all.splice(index, 1);
	}
	if (storeCommunication.all.length > 0) {
		const nextIndex = Math.min(index, storeCommunication.all.length - 1);
		selectedId.value = storeCommunication.all[nextIndex].id;
	} else {
		selectedId.value = null;
	}
}
async function sendMessage(message: string) {
	if (!message.trim()) return;
	if (selectedCommunication.value == null) return;
	if (selectedReport.value == null) return;
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
		`Message successfully sent to ${selectedReport.value.reporter.name}`,
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

// Lifecycle hooks
onMounted(async () => {
	await storeCommunication.fetchAll();
});
watch(
	paramCommunication.value,
	(communication_id) => {
		if (communication_id) {
			handleSelect(communication_id);
		}
	},
	{ immediate: true },
);
</script>
