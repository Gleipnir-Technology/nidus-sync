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
	<div class="p-3">
		<div class="map-container">
			<MapMultipoint
				id="map"
				ref="mapRef"
				:organization-id="user.organization.id"
				:tegola="user.urls.tegola"
				:xmin="user.organization.service_area?.min.x ?? 0"
				:ymin="user.organization.service_area?.min.y ?? 0"
				:xmax="user.organization.service_area?.max.x ?? 0"
				:ymax="user.organization.service_area?.max.y ?? 0"
			/>
		</div>
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
			<!-- Report Details -->
			<div class="details-section p-3 border-top">
				<div class="d-flex justify-content-between align-items-start mb-3">
					<div>
						<h5 class="mb-1">
							<span
								v-if="selectedCommunication.type === 'publicreport.nuisance'"
							>
								<i class="bi bi-mosquito icon-nuisance"></i>
								Nuisance Report
							</span>
							<span v-if="selectedCommunication.type === 'publicreport.water'">
								<i class="bi bi-droplet-fill icon-standing-water"></i>
								Standing Water Report
							</span>
						</h5>
						<small class="text-muted"
							>Report ID: #{{ selectedCommunication.id }}</small
						>
					</div>
					<span class="badge bg-secondary">
						<TimeRelative :time="selectedCommunication.created" />
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
										formatAddress(selectedCommunication.public_report.address)
									}}
								</div>
							</div>
							<div class="col-md-6">
								<label class="form-label text-muted small mb-0">
									<i class="bi bi-person"></i> Reporter Name
								</label>
								<div class="fw-medium">
									{{
										selectedCommunication.public_report.reporter.name ||
										"not given"
									}}
								</div>
							</div>
							<div class="col-md-6">
								<label
									v-if="selectedCommunication.public_report.reporter.has_email"
									class="form-label text-muted small mb-0"
								>
									<i class="bi bi-envelope"></i>
								</label>
								<label
									v-if="selectedCommunication.public_report.reporter.has_phone"
									class="form-label text-muted small mb-0"
								>
									<i class="bi bi-phone"></i>
								</label>
							</div>
						</div>
						<div v-if="water" class="row g-3">
							<div class="col-12">
								<ul>
									<li v-if="water?.is_reporter_owner">
										Reporter is the owner of the property
									</li>
									<li v-if="water?.is_reporter_confidential">
										Reporter has asked to be kept confidential
									</li>
								</ul>
							</div>
						</div>
					</div>
				</div>

				<!-- Nuisance-specific Fields -->
				<div v-if="nuisance" class="card mb-3">
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
									<li v-if="nuisance?.time_of_day_early">Early</li>
									<li v-if="nuisance?.time_of_day_day">Daytime</li>
									<li v-if="nuisance?.time_of_day_evening">Evening</li>
									<li v-if="nuisance?.time_of_day_night">Night</li>
								</ul>
							</div>
							<div class="col-md-6">
								<label class="form-label text-muted small mb-0">
									<i class="bi bi-house"></i> Property Area
								</label>
								<div>
									<ul>
										<li v-if="nuisance?.is_location_backyard">Backyard</li>
										<li v-if="nuisance?.is_location_frontyard">Frontyard</li>
										<li v-if="nuisance?.is_location_garden">Garden</li>
										<li v-if="nuisance?.is_location_other">Other</li>
										<li v-if="nuisance?.is_location_pool">Pool</li>
									</ul>
								</div>
							</div>
							<div
								v-if="
									nuisance?.source_container ||
									nuisance?.source_gutter ||
									nuisance?.source_stagnant
								"
								class="col-md-6"
							>
								<label class="form-label text-muted small mb-0">
									<i class="bi bi-droplet"></i> Sources
								</label>
								<ul>
									<li v-if="nuisance?.source_container">Container</li>
									<li v-if="nuisance?.source_gutter">Gutter</li>
									<li v-if="nuisance?.source_stagnant">Sprinklers & Gutters</li>
								</ul>
							</div>
							<div v-if="nuisance?.source_description" class="col-12">
								<label class="form-label text-muted small mb-0">
									<i class="bi bi-chat-text"></i> Source Description
								</label>
								<div class="p-2 bg-light rounded">
									{{ nuisance?.source_description || "none" }}
								</div>
							</div>
							<div class="col-12">
								<label class="form-label text-mudet small mb-0">
									<i class="bi bi-clock"></i> Duration
								</label>
								<div class="p-2 bg-light rounded">
									{{ nuisance?.duration }}
								</div>
							</div>
							<div class="col-12">
								<label class="form-label text-muted small mb-0">
									<i class="bi bi-chat-text"></i> Additional Notes
								</label>
								<div class="p-2 bg-light rounded">
									{{ nuisance?.additional_info || "No additional notes" }}
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- Standing Water-specific Fields -->
				<div v-if="water" class="card mb-3">
					<div class="card-header bg-info bg-opacity-10">
						<i class="bi bi-droplet"></i> Standing Water Details
					</div>
					<div class="card-body">
						<div
							v-if="
								water?.access_gate ||
								water?.access_fence ||
								water?.access_locked ||
								water?.access_dog ||
								water?.access_other
							"
							class="col-md-6"
						>
							<label class="form-label text-muted small mb-0">
								<i class="bi bi-droplet"></i> Access
							</label>
							<div>
								<ul>
									<li v-if="water?.access_gate">Gate</li>
									<li v-if="water?.access_fence">Fence</li>
									<li v-if="water?.access_locked">Locked</li>
									<li v-if="water?.access_dog">Dog</li>
									<li v-if="water?.access_other">Other access obstacle</li>
								</ul>
							</div>
						</div>
						<div v-if="water?.access_comments" class="col-12">
							<label class="form-label text-muted small mb-0">
								<i class="bi bi-chat-text"></i> Access Comments
							</label>
							<div class="p-2 bg-light rounded">
								{{ water?.access_comments }}
							</div>
						</div>
						<label class="form-label text-muted small mb-0">
							<i class="bi bi-eye"></i> Mosquito Life Stages Observed
						</label>
						<div class="mt-2">
							<span
								class="badge me-2"
								:class="
									water?.has_larvae ? 'badge-larvae' : 'bg-light text-muted'
								"
							>
								<i
									class="bi"
									:class="water?.has_larvae ? 'bi-check-circle' : 'bi-circle'"
								></i>
								Larvae
							</span>
							<span
								class="badge me-2"
								:class="
									water?.has_pupae ? 'badge-pupae' : 'bg-light text-muted'
								"
							>
								<i
									class="bi"
									:class="water?.has_pupae ? 'bi-check-circle' : 'bi-circle'"
								></i>
								Pupae
							</span>
							<span
								class="badge"
								:class="
									water?.has_adult ? 'badge-adult' : 'bg-light text-muted'
								"
							>
								<i
									class="bi"
									:class="water?.has_adult ? 'bi-check-circle' : 'bi-circle'"
								></i>
								Adult Mosquitoes
							</span>
						</div>
						<div v-if="water?.comments" class="col-12">
							<label class="form-label text-muted small mb-0">
								<i class="bi bi-chat-text"></i> Comments
							</label>
							<div class="p-2 bg-light rounded">
								{{ water?.comments }}
							</div>
						</div>
						<div class="col-md-6">
							<label class="form-label text-muted small mb-0">
								<i class="bi bi-person"></i> Owner Name
							</label>
							<div class="fw-medium">
								{{ water?.owner.name || "not given" }}
							</div>
						</div>
						<div class="col-md-6">
							<label
								v-if="water?.owner.has_email"
								class="form-label text-muted small mb-0"
							>
								<i class="bi bi-envelope"></i>
							</label>
							<label
								v-if="water?.owner.has_phone"
								class="form-label text-muted small mb-0"
							>
								<i class="bi bi-phone"></i>
							</label>
						</div>
					</div>
				</div>

				<!-- Photos Section -->
				<div class="card">
					<div
						class="card-header d-flex justify-content-between align-items-center"
					>
						<span><i class="bi bi-images"></i> Attached Photos</span>
						<span class="badge bg-primary">
							{{ selectedCommunication.public_report.images?.length || 0 }}
						</span>
					</div>
					<div class="card-body">
						<div
							v-if="
								selectedCommunication.public_report.images &&
								selectedCommunication.public_report.images.length > 0
							"
							class="d-flex flex-wrap gap-2"
						>
							<img
								v-for="(photo, index) in selectedCommunication.public_report
									.images"
								:key="index"
								:src="photo.url_content"
								class="photo-thumbnail"
								@click="openPhotoViewer(index)"
								:alt="'Photo ' + (index + 1)"
							/>
						</div>
						<div
							v-if="
								!selectedCommunication.public_report.images ||
								selectedCommunication.public_report.images.length === 0
							"
							class="text-muted text-center py-3"
						>
							<i class="bi bi-camera-slash fs-4"></i>
							<p class="mb-0 small">No images attached</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import MapMultipoint from "../components/MapMultipoint.vue";
import TimeRelative from "../components/TimeRelative.vue";

interface Props {
	loading: boolean;
	selectedCommunication: Communication | null;
	user: User | null;
}

const props = defineProps<Props>();
const nuisance = computed(() => {
	return props.selectedCommunication?.value?.public_report?.nuisance || null;
});
const water = computed(() => {
	return props.selectedCommunication?.value?.public_report?.water || null;
});

function formatAddress(a) {
	if (a.number === "" && a.street === "") {
		return "no address provided";
	}
	return `${a.number} ${a.street}, ${a.locality}`;
}
</script>
