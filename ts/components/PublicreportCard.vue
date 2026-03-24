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
					<span
						v-if="
							report.type === 'nuisance'
						"
					>
						<i class="bi bi-mosquito icon-nuisance"></i>
						Nuisance Report
					</span>
					<span
						v-if="report.type === 'water'"
					>
						<i class="bi bi-droplet-fill icon-standing-water"></i>
						Standing Water Report
					</span>
				</h5>
				<small class="text-muted"
					>Report ID: #{{ report.id }}</small
				>
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
							{{
								formatAddress(
									report.address,
								)
							}}
						</div>
					</div>
					<div class="col-md-6">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-person"></i> Reporter Name
						</label>
						<div class="fw-medium">
							{{
								report.reporter.name ||
								"not given"
							}}
						</div>
					</div>
					<div class="col-md-6">
						<label
							v-if="
								report.reporter.has_email
							"
							class="form-label text-muted small mb-0"
						>
							<i class="bi bi-envelope"></i>
						</label>
						<label
							v-if="
								report.reporter.has_phone
							"
							class="form-label text-muted small mb-0"
						>
							<i class="bi bi-phone"></i>
						</label>
					</div>
				</div>
				<div v-if="report.water != null" class="row g-3">
					<div class="col-12">
						<ul>
							<li v-if="report.water?.is_reporter_owner">
								Reporter is the owner of the property
							</li>
							<li v-if="report.water?.is_reporter_confidential">
								Reporter has asked to be kept confidential
							</li>
						</ul>
					</div>
				</div>
			</div>
		</div>

		<!-- Nuisance-specific Fields -->
		<div v-if="report.nuisance" class="card mb-3">
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
							<li v-if="report.nuisance?.time_of_day_early">Early</li>
							<li v-if="report.nuisance?.time_of_day_day">Daytime</li>
							<li v-if="report.nuisance?.time_of_day_evening">Evening</li>
							<li v-if="report.nuisance?.time_of_day_night">Night</li>
						</ul>
					</div>
					<div class="col-md-6">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-house"></i> Property Area
						</label>
						<div>
							<ul>
								<li v-if="report.nuisance?.is_location_backyard">Backyard</li>
								<li v-if="report.nuisance?.is_location_frontyard">
									Frontyard
								</li>
								<li v-if="report.nuisance?.is_location_garden">Garden</li>
								<li v-if="report.nuisance?.is_location_other">Other</li>
								<li v-if="report.nuisance?.is_location_pool">Pool</li>
							</ul>
						</div>
					</div>
					<div
						v-if="
							report.nuisance?.source_container ||
							report.nuisance?.source_gutter ||
							report.nuisance?.source_stagnant
						"
						class="col-md-6"
					>
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-droplet"></i> Sources
						</label>
						<ul>
							<li v-if="report.nuisance?.source_container">Container</li>
							<li v-if="report.nuisance?.source_gutter">Gutter</li>
							<li v-if="report.nuisance?.source_stagnant">
								Sprinklers & Gutters
							</li>
						</ul>
					</div>
					<div v-if="report.nuisance?.source_description" class="col-12">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-chat-text"></i> Source Description
						</label>
						<div class="p-2 bg-light rounded">
							{{ report.nuisance?.source_description || "none" }}
						</div>
					</div>
					<div class="col-12">
						<label class="form-label text-mudet small mb-0">
							<i class="bi bi-clock"></i> Duration
						</label>
						<div class="p-2 bg-light rounded">
							{{ report.nuisance?.duration }}
						</div>
					</div>
					<div class="col-12">
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-chat-text"></i> Additional Notes
						</label>
						<div class="p-2 bg-light rounded">
							{{ report.nuisance?.additional_info || "No additional notes" }}
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Standing Water-specific Fields -->
		<div v-if="report.water" class="card mb-3">
			<div class="card-header bg-info bg-opacity-10">
				<i class="bi bi-droplet"></i> Standing Water Details
			</div>
			<div class="card-body">
				<div
					v-if="
						report.water?.access_gate ||
						report.water?.access_fence ||
						report.water?.access_locked ||
						report.water?.access_dog ||
						report.water?.access_other
					"
					class="col-md-6"
				>
					<label class="form-label text-muted small mb-0">
						<i class="bi bi-droplet"></i> Access
					</label>
					<div>
						<ul>
							<li v-if="report.water?.access_gate">Gate</li>
							<li v-if="report.water?.access_fence">Fence</li>
							<li v-if="report.water?.access_locked">Locked</li>
							<li v-if="report.water?.access_dog">Dog</li>
							<li v-if="report.water?.access_other">Other access obstacle</li>
						</ul>
					</div>
				</div>
				<div v-if="report.water?.access_comments" class="col-12">
					<label class="form-label text-muted small mb-0">
						<i class="bi bi-chat-text"></i> Access Comments
					</label>
					<div class="p-2 bg-light rounded">
						{{ report.water?.access_comments }}
					</div>
				</div>
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-eye"></i> Mosquito Life Stages Observed
				</label>
				<div class="mt-2">
					<span
						class="badge me-2"
						:class="
							report.water?.has_larvae ? 'badge-larvae' : 'bg-light text-muted'
						"
					>
						<i
							class="bi"
							:class="
								report.water?.has_larvae ? 'bi-check-circle' : 'bi-circle'
							"
						></i>
						Larvae
					</span>
					<span
						class="badge me-2"
						:class="
							report.water?.has_pupae ? 'badge-pupae' : 'bg-light text-muted'
						"
					>
						<i
							class="bi"
							:class="
								report.water?.has_pupae ? 'bi-check-circle' : 'bi-circle'
							"
						></i>
						Pupae
					</span>
					<span
						class="badge"
						:class="
							report.water?.has_adult ? 'badge-adult' : 'bg-light text-muted'
						"
					>
						<i
							class="bi"
							:class="
								report.water?.has_adult ? 'bi-check-circle' : 'bi-circle'
							"
						></i>
						Adult Mosquitoes
					</span>
				</div>
				<div v-if="report.water?.comments" class="col-12">
					<label class="form-label text-muted small mb-0">
						<i class="bi bi-chat-text"></i> Comments
					</label>
					<div class="p-2 bg-light rounded">
						{{ report.water?.comments }}
					</div>
				</div>
				<div class="col-md-6">
					<label class="form-label text-muted small mb-0">
						<i class="bi bi-person"></i> Owner Name
					</label>
					<div class="fw-medium">
						{{ report.water?.owner.name || "not given" }}
					</div>
				</div>
				<div class="col-md-6">
					<label
						v-if="report.water?.owner.has_email"
						class="form-label text-muted small mb-0"
					>
						<i class="bi bi-envelope"></i>
					</label>
					<label
						v-if="report.water?.owner.has_phone"
						class="form-label text-muted small mb-0"
					>
						<i class="bi bi-phone"></i>
					</label>
				</div>
			</div>
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
					v-if="
						report.images &&
						report.images.length > 0
					"
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
					v-if="
						!report.images ||
						report.images.length === 0
					"
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
import PublicreportCard from "@/components/PublicreportCard.vue";
import TimeRelative from "@/components/TimeRelative.vue";

interface Emits {
	(e: "viewImage", index: int): void;
}
interface Props {
	report: Publicreport;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
function formatAddress(a) {
	if (a.number === "" && a.street === "") {
		return "no address provided";
	}
	return `${a.number} ${a.street}, ${a.locality}`;
}
</script>
