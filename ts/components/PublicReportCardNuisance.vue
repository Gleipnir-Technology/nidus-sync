<template>
	<div class="card-header bg-danger bg-opacity-10">
		<i class="bi bi-exclamation-triangle"></i> Nuisance Details
	</div>
	<div class="card-body">
		<div class="row g-3">
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-clock"></i> Time of Day Encountered
				</label>
				<ul>
					<li v-if="report.time_of_day_early">Early</li>
					<li v-if="report.time_of_day_day">Daytime</li>
					<li v-if="report.time_of_day_evening">Evening</li>
					<li v-if="report.time_of_day_night">Night</li>
				</ul>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-house"></i> Property Area
				</label>
				<div>
					<ul>
						<li v-if="report.is_location_backyard">Backyard</li>
						<li v-if="report.is_location_frontyard">Frontyard</li>
						<li v-if="report.is_location_garden">Garden</li>
						<li v-if="report.is_location_other">Other</li>
						<li v-if="report.is_location_pool">Pool</li>
					</ul>
				</div>
			</div>
			<div
				v-if="
					report.source_container ||
					report.source_gutter ||
					report.source_stagnant
				"
				class="col-md-6"
			>
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-droplet"></i> Sources
				</label>
				<ul>
					<li v-if="report.source_container">Container</li>
					<li v-if="report.source_gutter">Gutter</li>
					<li v-if="report.source_stagnant">Sprinklers & Gutters</li>
				</ul>
			</div>
			<div v-if="report.source_description" class="col-12">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-chat-text"></i> Source Description
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.source_description || "none" }}
				</div>
			</div>
			<div class="col-12">
				<label class="form-label text-mudet small mb-0">
					<i class="bi bi-clock"></i> Duration
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.duration }}
				</div>
			</div>
			<div class="col-12">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-chat-text"></i> Additional Notes
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.additional_info || "No additional notes" }}
				</div>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { PublicReportNuisance } from "@/type/api";
interface Props {
	report: PublicReportNuisance;
}
const props = defineProps<Props>();
</script>
