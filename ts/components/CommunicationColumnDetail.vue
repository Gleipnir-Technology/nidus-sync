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
			<div v-if="loading || session.user == null" class="loading">
				Loading...
			</div>
			<div v-else>
				<div class="map-container">
					<MapMultipoint
						id="map"
						:bounds="mapBounds"
						:markers="mapMarkers"
						:organizationId="session.user?.organization.id"
						:tegola="session.urls?.tegola ?? ''"
					/>
				</div>
				<div
					v-if="!selectedCommunication"
					class="d-flex flex-column align-items-center justify-content-center text-muted"
				>
					<i class="bi bi-hand-index fs-1"></i>
					<p class="mt-2">Select a report to view details</p>
				</div>

				<div v-if="selectedCommunication" class="h-100 d-flex flex-column">
					<PublicreportCard
						v-if="selectedCommunication?.public_report"
						:report="selectedCommunication?.public_report"
						@viewImage="openPhotoViewer"
					/>
					<p v-else>No public report</p>
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
import { Bounds, Communication, Marker, User } from "@/types";
import { useSessionStore } from "@/store/session";

interface Emits {
	(e: "viewImage", index: number): void;
}
interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapMarkers: Marker[];
	selectedCommunication: Communication | null;
}

const emit = defineEmits<Emits>();
const props = defineProps<Props>();
const nuisance = computed(() => {
	return props.selectedCommunication?.public_report?.nuisance || null;
});
const session = useSessionStore();
const water = computed(() => {
	return props.selectedCommunication?.public_report?.water || null;
});
function openPhotoViewer(index: number) {
	emit("viewImage", index);
}
</script>
