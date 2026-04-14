<style scoped>
.photo-thumbnail {
	width: 100px;
	height: 100px;
	object-fit: cover;
	cursor: pointer;
	border-radius: 4px;
	transition: transform 0.2s;
}
.photo-thumbnail:hover {
	transform: scale(1.05);
}
</style>
<template>
	<div class="details-section p-3 border-top">
		<div class="d-flex justify-content-between align-items-start mb-3">
			<div>
				<h5 class="mb-1">
					<span v-if="report.type === 'nuisance'">
						<i class="bi bi-mosquito icon-nuisance"></i>
						Nuisance Report
					</span>
					<span v-if="report.type === 'water'">
						<i class="bi bi-droplet-fill icon-standing-water"></i>
						Standing Water Report
					</span>
				</h5>
				<small class="text-muted">Report ID: #{{ report.public_id }}</small>
			</div>
			<span class="badge bg-secondary">
				<TimeRelative :time="report.created" />
			</span>
		</div>

		<!-- Common Fields -->
		<div class="card mb-3">
			<div class="card-body">
				<div class="row g-3">
					<div class="col-12">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-geo-alt"></i> Address
						</label>
						<div class="fw-medium">
							{{ formatAddress(report.address) }}
						</div>
					</div>
					<div class="col-md-6">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-person"></i> Reporter Name
						</label>
						<div class="fw-medium">
							{{ report.reporter.name || "not given" }}
						</div>
					</div>
					<div class="col-md-6">
						<label
							v-if="report.reporter.has_email"
							class="form-label text-muted small mb-0"
						>
							<i class="bi bi-envelope"></i>
						</label>
						<label
							v-if="report.reporter.has_phone"
							class="form-label text-muted small mb-0"
						>
							<i class="bi bi-phone"></i>
						</label>
					</div>
				</div>
				<div v-if="report instanceof PublicReportWater" class="row g-3">
					<div class="col-12">
						<ul>
							<li v-if="report.is_reporter_owner">
								Reporter is the owner of the property
							</li>
							<li v-if="report.is_reporter_confidential">
								Reporter has asked to be kept confidential
							</li>
						</ul>
					</div>
				</div>
			</div>
		</div>

		<div v-if="report instanceof PublicReportNuisance">
			<PublicReportCardNuisance :report="report" />
		</div>

		<div v-if="report instanceof PublicReportWater">
			<PublicReportCardWater :report="report" />
		</div>

		<!-- Images Section -->
		<div class="card">
			<div
				class="card-header d-flex justify-content-between align-items-center"
			>
				<span><i class="bi bi-images"></i> Attached Photos</span>
				<span class="badge bg-primary">
					{{ report.images?.length || 0 }}
				</span>
			</div>
			<div class="card-body">
				<div
					v-if="report.images && report.images.length > 0"
					class="d-flex flex-wrap gap-2"
				>
					<img
						v-for="(photo, index) in report.images"
						:key="index"
						:src="photo.url_content"
						class="photo-thumbnail"
						@click="openPhotoViewer(index)"
						:alt="'Photo ' + (index + 1)"
					/>
				</div>
				<div
					v-if="!report.images || report.images.length === 0"
					class="text-muted text-center py-3"
				>
					<i class="bi bi-camera-slash fs-4"></i>
					<p class="mb-0 small">No images attached</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import MapMultipoint from "@/components/MapMultipoint.vue";
import TimeRelative from "@/components/TimeRelative.vue";
import PublicReportCardNuisance from "@/components/PublicReportCardNuisance.vue";
import PublicReportCardWater from "@/components/PublicReportCardWater.vue";
import { formatAddress } from "@/format";
import {
	PublicReport,
	PublicReportNuisance,
	PublicReportWater,
} from "@/type/api";

interface Emits {
	(e: "viewImage", index: number): void;
}
interface Props {
	report: PublicReport;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
function openPhotoViewer(index: number) {
	emit("viewImage", index);
}
</script>
