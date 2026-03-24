<style scoped>
.badge-larvae {
	background-color: #ffc107;
	color: #000;
}

.badge-pupae {
	background-color: #fd7e14;
	color: #fff;
}

.badge-adult {
	background-color: #dc3545;
	color: #fff;
}
.details-section {
	overflow-y: auto;
}
.icon-standing-water {
	color: #0dcaf0;
}

.icon-nuisance {
	color: #dc3545;
}
.map-container {
	height: 400px;
	width: 100%;
}
</style>

<template>
	<div class="card shadow-sm mb-3">
		<div class="card-header bg-white pane-header">Communication Workbench</div>
		<div class="card-body">
			<div class="map-container">
				<MapMultipoint
					id="map"
					:bounds="mapBounds"
					:markers="mapMarkers"
					:organization-id="user.organization.id"
					:tegola="user.urls.tegola"
					:xmin="user.organization.service_area?.min.x ?? 0"
					:ymin="user.organization.service_area?.min.y ?? 0"
					:xmax="user.organization.service_area?.max.x ?? 0"
					:ymax="user.organization.service_area?.max.y ?? 0"
				/>
			</div>
			<div v-if="loading" class="loading">Loading...</div>
			<div v-else>
				<div
					v-if="!selectedCommunication"
					class="d-flex flex-column align-items-center justify-content-center text-muted"
				>
					<i class="bi bi-hand-index fs-1"></i>
					<p class="mt-2">Select a report to view details</p>
				</div>

				<div v-if="selectedCommunication" class="h-100 d-flex flex-column">
					<PublicreportCard :report="selectedCommunication.public_report" @viewImage="openPhotoViewer" />
					<!-- Report Details -->
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import MapMultipoint from "@/components/MapMultipoint.vue";
import PublicreportCard from "@/components/PublicreportCard.vue";
import TimeRelative from "@/components/TimeRelative.vue";

interface Emits {
	(e: "viewImage", index: int): void;
}
interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapMarkers: Marker[];
	selectedCommunication: Communication | null;
	user: User | null;
}

const emit = defineEmits<Emits>();
const props = defineProps<Props>();
const nuisance = computed(() => {
	return props.selectedCommunication?.value?.public_report?.nuisance || null;
});
const water = computed(() => {
	return props.selectedCommunication?.value?.public_report?.water || null;
});
function openPhotoViewer(index) {
	emit("viewImage", index);
}
</script>
