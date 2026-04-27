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
	background-color: #e9ecef;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 400px;
	margin-top: 20px;
	position: relative;
}
</style>

<template>
	<div class="card shadow-sm mb-3">
		<div class="card-header bg-white pane-header">Communication Workbench</div>
		<div class="card-body">
			<div
				v-if="loading || session.self == null || session.organization == null"
				class="loading"
			>
				Loading...
			</div>
			<div v-else>
				<div class="map-container">
					<Map
						:bounds="mapBounds"
						:markers="mapMarkers"
						:organizationId="session.organization?.id"
					>
						<Layer
							id="parcel"
							:minzoom="14"
							:paint="{ 'line-color': '#0f0' }"
							source="tegola"
							sourceLayer="parcel"
							type="line"
						/>
						<Layer
							id="service-area"
							:paint="{ 'line-color': '#f00' }"
							source="tegola"
							sourceLayer="service-area-bounds"
							type="line"
						/>
						<Source
							id="tegola"
							type="vector"
							:tiles="[
								session.urls?.tegola +
									'maps/nidus/{z}/{x}/{y}?id=' +
									session.organization?.id,
							]"
						/>
					</Map>
				</div>
				<div
					v-if="!selectedCommunication"
					class="d-flex flex-column align-items-center justify-content-center text-muted"
				>
					<i class="bi bi-hand-index fs-1"></i>
					<p class="mt-2">Select a report to view details</p>
				</div>

				<div v-if="selectedCommunication" class="h-100 d-flex flex-column">
					<PublicReportCard
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
import Map, { LngLatBounds } from "@/map/Map.vue";
import Layer from "@/map/Layer.vue";
import Source from "@/map/Source.vue";
import PublicReportCard from "@/components/PublicReportCard.vue";
import TimeRelative from "@/components/TimeRelative.vue";
import type { Marker } from "@/types";
import type { Bounds, Communication, User } from "@/type/api";
import { useSessionStore } from "@/store/session";

interface Emits {
	(e: "viewImage", index: number): void;
}
interface Props {
	loading: boolean;
	mapBounds?: LngLatBounds;
	mapMarkers: Marker[];
	selectedCommunication: Communication | null;
}

const emit = defineEmits<Emits>();
const props = defineProps<Props>();
const session = useSessionStore();
function openPhotoViewer(index: number) {
	emit("viewImage", index);
}
</script>
