<style scoped>
.capitalized {
	text-transform: capitalize;
}
.map-container {
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	margin-bottom: 20px;
	margin-top: 20px;
	align-items: center;
	justify-content: center;
	/* Prevent touch scrolling issues */
	touch-action: pan-y pinch-zoom;
}
#map {
	width: 100%;
	height: 100%;
}
.status-badge {
	font-size: 1rem;
}
.timeline {
	border-left: 3px solid #dee2e6;
	padding-left: 20px;
	margin-left: 10px;
}
.timeline-item {
	position: relative;
	margin-bottom: 25px;
}
.timeline-item:before {
	content: "";
	position: absolute;
	left: -29px;
	top: 0;
	width: 16px;
	height: 16px;
	border-radius: 50%;
	background-color: #0d6efd;
}
.timeline-date {
	font-size: 0.85rem;
	color: #6c757d;
}
</style>
<template>
	<HeaderDistrict :district="district" v-if="district" />
	<Header v-else />
	<div class="container my-4" v-if="report">
		<!-- Report ID and Status Section -->
		<div class="card mb-4">
			<div
				class="card-header bg-primary text-white d-flex justify-content-between align-items-center"
			>
				<h5 class="mb-0">Report {{ formatReportID(id) }}</h5>
				<span class="badge bg-warning capitalized status-badge text-dark">
					{{ report.status }}
				</span>
			</div>
			<div class="card-body">
				<div class="row">
					<div class="col-md-4 mb-3">
						<strong><i class="bi bi-tag me-2"></i>Type:</strong>
						<span class="capitalized">{{ report.type }}</span>
					</div>
					<div class="col-md-4 mb-3">
						<strong><i class="bi bi-calendar me-2"></i>Created:</strong>
						<span>{{ formatTimeRelative(report.created) }}</span>
					</div>
					<div class="col-md-4 mb-3" v-if="district">
						<strong><i class="bi bi-crosshair me-2"></i>District:</strong>
						<span>
							{{ district.name }}
						</span>
					</div>
				</div>
				<div class="row">
					<div class="col-md-12">
						<strong><i class="bi bi-pin-map me-2"></i>Location:</strong>
						<span>{{ report.address.raw }}</span>
					</div>
				</div>
				<div class="row">
					<div class="col-md-12">
						<strong><i class="bi bi-images me-2"></i>Images:</strong>
						<span>
							{{
								report.images.length > 0
									? report.images.length
									: "None provided"
							}}
						</span>
					</div>
				</div>
				<ReportDetailNuisance
					:nuisance="report as PublicReportNuisance"
					v-if="report instanceof PublicReportNuisance"
				/>
				<ReportDetailWater
					:water="report as PublicReportWater"
					v-if="report instanceof PublicReportWater"
				/>
			</div>
		</div>

		<!-- Map Section -->
		<div class="card mb-4">
			<div class="card-header bg-info text-white">
				<h5 class="mb-0">
					<i class="bi bi-pin-map-fill me-2"></i>Location Map
				</h5>
			</div>
			<div class="card-body p-0">
				<div class="map-container">
					<MapLocatorDisplay id="map" :markers="markers"></MapLocatorDisplay>
				</div>
			</div>
		</div>
		<!-- History Timeline -->
		<div class="card">
			<div class="card-header bg-success text-white">
				<h5 class="mb-0">
					<i class="bi bi-clock-history me-2"></i>Request History
				</h5>
			</div>
			<div class="card-body">
				<div class="timeline">
					<div
						v-for="(item, index) in report.log"
						:key="index"
						class="timeline-item"
					>
						<div class="timeline-date">
							{{ formatTimeRelative(item.created) }}
						</div>
						<h5 class="capitalized mb-1">{{ item.type }}</h5>
						<p class="mb-0">{{ item.message }}</p>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div class="container my-4" v-else>
		<p>loading...</p>
	</div>
</template>
<script setup lang="ts">
import { ref, computed } from "vue";
import { computedAsync } from "@vueuse/core";
import Header from "@/rmo/components/Header.vue";
import HeaderDistrict from "@/components/HeaderDistrict.vue";
import MapLocatorDisplay from "@/components/MapLocatorDisplay.vue";
import ReportDetailNuisance from "@/rmo/components/ReportDetailNuisance.vue";
import ReportDetailWater from "@/rmo/components/ReportDetailWater.vue";
import { useStoreDistrict } from "@/rmo/store/district";
import { useStorePublicReport } from "@/rmo/store/publicreport";
import type { Marker } from "@/types";
import {
	type District,
	PublicReport,
	PublicReportNuisance,
	PublicReportWater,
} from "@/type/api";
import { formatReportID, formatTimeRelative } from "@/format";

// Props
interface Props {
	id: string;
}

const props = defineProps<Props>();
const storeDistrict = useStoreDistrict();
const storePublicReport = useStorePublicReport();
// Computed
const report = computedAsync(async (): Promise<PublicReport | undefined> => {
	return await storePublicReport.byID(props.id);
});
const district = computedAsync(async (): Promise<District | undefined> => {
	if (!(report.value && report.value.district)) {
		return undefined;
	}
	return await storeDistrict.byURI(report.value.district);
});
const markers = computed((): Marker[] => {
	if (
		!(
			report.value &&
			report.value.address.location &&
			(report.value.address.location.latitude ||
				report.value.address.location.longitude)
		)
	) {
		return [];
	}
	return [
		{
			id: props.id,
			location: report.value.address.location,
		},
	];
});
</script>
