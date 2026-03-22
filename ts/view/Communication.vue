<style scoped>
/* Add your component-specific styles here */
.reports-list {
	overflow-y: auto;
	max-height: 100vh;
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

.details-section {
	overflow-y: auto;
}

.actions-panel {
	height: 100%;
	overflow-y: auto;
}

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

.icon-standing-water {
	color: #0dcaf0;
}

.icon-nuisance {
	color: #dc3545;
}

.modal.show {
	background-color: rgba(0, 0, 0, 0.5);
}
</style>
<template>
	<div class="h-100">
		<div class="container-fluid h-100">
			<div class="row h-100">
				<!-- Left Column - Communications List -->
				<CommunicationColumnList
					:all="communication.all"
					:loading="loading"
					:selected-id="selectedId"
					@select="handleSelect"
				/>

				<!-- Middle Column - Report Details -->
				<div class="col-md-6 p-0">
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
							<div
								class="d-flex justify-content-between align-items-start mb-3"
							>
								<div>
									<h5 class="mb-1">
										<span
											v-if="
												selectedCommunication.type === 'publicreport.nuisance'
											"
										>
											<i class="bi bi-mosquito icon-nuisance"></i>
											Nuisance Report
										</span>
										<span
											v-if="selectedCommunication.type === 'publicreport.water'"
										>
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
													formatAddress(
														selectedCommunication.public_report.address,
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
													selectedCommunication.public_report.reporter.name ||
													"not given"
												}}
											</div>
										</div>
										<div class="col-md-6">
											<label
												v-if="
													selectedCommunication.public_report.reporter.has_email
												"
												class="form-label text-muted small mb-0"
											>
												<i class="bi bi-envelope"></i>
											</label>
											<label
												v-if="
													selectedCommunication.public_report.reporter.has_phone
												"
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
													<li v-if="nuisance?.is_location_backyard">
														Backyard
													</li>
													<li v-if="nuisance?.is_location_frontyard">
														Frontyard
													</li>
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
												<li v-if="nuisance?.source_stagnant">
													Sprinklers & Gutters
												</li>
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
												<li v-if="water?.access_other">
													Other access obstacle
												</li>
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
												water?.has_larvae
													? 'badge-larvae'
													: 'bg-light text-muted'
											"
										>
											<i
												class="bi"
												:class="
													water?.has_larvae ? 'bi-check-circle' : 'bi-circle'
												"
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
												:class="
													water?.has_pupae ? 'bi-check-circle' : 'bi-circle'
												"
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
												:class="
													water?.has_adult ? 'bi-check-circle' : 'bi-circle'
												"
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
										{{
											selectedCommunication.public_report.images?.length || 0
										}}
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
											v-for="(photo, index) in selectedCommunication
												.public_report.images"
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

				<!-- Right Column - Actions -->
				<div class="col-md-3 border-start p-0">
					<div
						v-if="!selectedCommunication"
						class="h-100 d-flex flex-column align-items-center justify-content-center text-muted p-3"
					>
						<i class="bi bi-gear fs-1"></i>
						<p class="mt-2 text-center">
							Actions will appear here when a report is selected
						</p>
					</div>

					<div
						v-if="selectedCommunication"
						class="actions-panel d-flex flex-column"
					>
						<div class="p-3 bg-light border-bottom">
							<h6 class="mb-0">
								<i class="bi bi-lightning"></i> Quick Actions
							</h6>
						</div>

						<div class="p-3 flex-grow-1">
							<!-- Create Signal -->
							<div class="d-grid mb-3">
								<button class="btn btn-success btn-lg" @click="createSignal()">
									<i class="bi bi-plus-circle me-2"></i>Mark Signal
								</button>
								<small class="text-muted mt-1"
									>This report is useful signal</small
								>
							</div>

							<!-- Mark Invalid -->
							<div class="d-grid mb-3">
								<button class="btn btn-outline-danger" @click="markInvalid()">
									<i class="bi bi-x-circle me-2"></i>Mark Invalid
								</button>
								<small class="text-muted mt-1">This report isn't useful</small>
							</div>

							<hr />

							<!-- Message Reporter -->
							<div
								v-if="
									!(
										selectedCommunication?.public_report.reporter.has_email ||
										selectedCommunication?.public_report.reporter.has_phone
									)
								"
								class="mb-3"
							>
								<h6>
									<i class="bi bi-chat-dots"></i> No Reporter Communications
									Available
								</h6>
							</div>
							<div
								v-if="
									selectedCommunication?.public_report.reporter.has_email ||
									selectedCommunication?.public_report.reporter.has_phone
								"
								class="mb-3"
							>
								<h6><i class="bi bi-chat-dots"></i> Message Reporter</h6>
								<div class="mb-2">
									<label class="form-label small text-muted"
										>Quick Templates</label
									>
									<select
										class="form-select form-select-sm"
										@change="applyMessageTemplate($event.target.value)"
									>
										<option value="">Select a template...</option>
										<option value="received">Report Received</option>
										<option value="scheduled">Service Scheduled</option>
										<option value="completed">Service Completed</option>
										<option value="need_info">Need More Information</option>
									</select>
								</div>
								<textarea
									class="form-control mb-2"
									rows="5"
									v-model="messageText"
									placeholder="Type your message to the reporter..."
								></textarea>
								<div class="d-grid">
									<button
										class="btn btn-primary"
										@click="sendMessage()"
										:disabled="!messageText.trim()"
									>
										<i class="bi bi-send me-2"></i>Send Message
									</button>
								</div>
							</div>

							<hr />

							<!-- Report History -->
							<div>
								<h6><i class="bi bi-clock-history"></i> Activity Log</h6>
								<div class="small">
									<div
										v-for="entry in selectedCommunication.public_report.log ||
										[]"
										:key="entry.created"
										class="border-start border-2 ps-2 mb-2"
									>
										<div v-if="entry.type === 'created'">
											<div class="text-muted">Initial Report</div>
											<small class="text-muted">{{
												formatDate(entry.created)
											}}</small>
										</div>
										<div v-else-if="entry.type === 'message-text'">
											<div class="text-muted">Text Message</div>
											<div>{{ entry.message }}</div>
											<small class="text-muted">{{
												formatDate(entry.created)
											}}</small>
										</div>
										<div v-else>{{ entry.type }}</div>
									</div>
									<div
										v-if="
											!selectedCommunication.public_report.log ||
											selectedCommunication.public_report.log.length === 0
										"
										class="text-muted"
									>
										No activity yet
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Photo Viewer Modal -->
		<div
			class="modal fade"
			:class="{ 'show d-block': showPhotoModal }"
			tabindex="-1"
			v-show="showPhotoModal"
			@click.self="showPhotoModal = false"
		>
			<div class="modal-dialog modal-lg modal-dialog-centered">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">
							Photo {{ currentPhotoIndex + 1 }} of
							{{ selectedCommunication?.public_report.images.length || 0 }}
						</h5>
						<button
							type="button"
							class="btn-close"
							@click="showPhotoModal = false"
						></button>
					</div>
					<div class="modal-body text-center">
						<div v-if="selectedCommunication && showPhotoModal">
							<img
								:src="
									selectedCommunication.public_report.images[currentPhotoIndex]
										.url_content
								"
								class="img-fluid rounded"
								style="max-height: 60vh"
							/>

							<!-- EXIF Data Section -->
							<div class="mt-4 pt-3 border-top text-start">
								<h6 class="text-muted mb-3">Photo Information</h6>
								<div class="row g-3">
									<div class="col-md-4">
										<small class="text-muted d-block">Date Taken</small>
										<span>
											{{
												selectedCommunication.public_report.images[
													currentPhotoIndex
												].exif?.created || "N/A"
											}}
										</span>
									</div>
									<div class="col-md-4">
										<small class="text-muted d-block">Camera</small>
										<span>
											{{
												(selectedCommunication.public_report.images[
													currentPhotoIndex
												].exif?.make || "") +
													" " +
													(selectedCommunication.public_report.images[
														currentPhotoIndex
													].exif?.model || "") || "N/A"
											}}
										</span>
									</div>
									<div class="col-md-4">
										<small class="text-muted d-block"
											>Distance from Reporter</small
										>
										<span
											v-if="
												selectedCommunication.public_report.images[
													currentPhotoIndex
												].location != null
											"
										>
											{{
												formatDistance(
													selectedCommunication.public_report.images[
														currentPhotoIndex
													].distance_from_reporter_meters,
												)
											}}
										</span>
										<span v-else>No location data in image</span>
									</div>
								</div>
							</div>
						</div>
					</div>
					<div class="modal-footer justify-content-between">
						<button
							class="btn btn-outline-secondary"
							@click="currentPhotoIndex = Math.max(0, currentPhotoIndex - 1)"
							:disabled="currentPhotoIndex === 0"
						>
							<i class="bi bi-chevron-left"></i> Previous
						</button>
						<button
							class="btn btn-outline-secondary"
							@click="
								currentPhotoIndex = Math.min(
									selectedCommunication.public_report.images.length - 1,
									currentPhotoIndex + 1,
								)
							"
							:disabled="
								currentPhotoIndex >=
								(selectedCommunication?.public_report.images?.length || 1) - 1
							"
						>
							Next <i class="bi bi-chevron-right"></i>
						</button>
					</div>
				</div>
			</div>
		</div>
		<div
			class="modal-backdrop fade show"
			v-show="showPhotoModal"
			@click="showPhotoModal = false"
		></div>

		<!-- Toast Notifications -->
		<div class="toast-container position-fixed bottom-0 end-0 p-3">
			<div class="toast" :class="{ show: showToast }" role="alert">
				<div class="toast-header">
					<i class="bi bi-check-circle text-success me-2"></i>
					<strong class="me-auto">{{ toastTitle }}</strong>
					<button
						type="button"
						class="btn-close"
						@click="showToast = false"
					></button>
				</div>
				<div class="toast-body">{{ toastMessage }}</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from "vue";
import maplibregl from "maplibre-gl";

import { useCommunicationStore } from "../store/communication";
import { useUserStore } from "../store/user";
import CommunicationColumnList from "../components/CommunicationColumnList.vue";
import MapMultipoint from "../components/MapMultipoint.vue";
import TimeRelative from "../components/TimeRelative.vue";

const communication = useCommunicationStore();
const user = useUserStore();
onMounted(() => {
	fetchCommunications();
});

// Refs
const apiBase = ref("/api");
const messageText = ref("");
const showPhotoModal = ref(false);
const selectedId = ref<string | null>(null);
const currentPhotoIndex = ref(0);
const showToast = ref(false);
const toastTitle = ref("");
const toastMessage = ref("");
const loading = ref(true);
const error = ref(null);
const mapRef = ref(null);

const nuisance = computed(() => {
	return selectedCommunication.value?.public_report?.nuisance || null;
});
const selectedCommunication = computed<Communication | null>(() => {
	if (selectedId.value == null) {
		return null;
	}
	if (communication.all == null) {
		return null;
	}
	const result = communication.all.find((c) => c.id == selectedId.value);
	return result;
});
const water = computed(() => {
	return selectedCommunication.value?.public_report?.water || null;
});

const handleSelect = (id: string) => {
	selectedId.value = id;
};
async function fetchCommunications() {
	await communication.fetchAll();
	// if we already had something selected, reset it using the new data
	if (selectedCommunication.value) {
		const matching = communication.all.filter((c) => {
			return c.id === selectedCommunication.value.id;
		});
		if (matching.length > 0) {
			selectedCommunication.value = matching[0];
		}
	}
}
function formatAddress(a) {
	if (a.number === "" && a.street === "") {
		return "no address provided";
	}
	return `$${a.number} $${a.street}, ${a.locality}`;
}

function formatDistance(meters) {
	if (meters === undefined || meters === null) {
		return "unknown";
	}
	if (meters < 1) {
		const mm = Math.round(meters * 1000);
		return `${mm} mm`;
	} else if (meters >= 1000) {
		const km = Math.round(meters / 1000);
		return `${km} km`;
	} else {
		const m = Math.round(meters);
		return `${m} m`;
	}
}

function formatDate(date) {
	return new Date(date).toLocaleString();
}

async function loadFromAPI() {
	loading.value = true;
	error.value = null;
	try {
		await Promise.all([fetchCommunications()]);
	} catch (err) {
		error.value = err.message;
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
}

function openPhotoViewer(index) {
	currentPhotoIndex.value = index;
	showPhotoModal.value = true;
}

function applyMessageTemplate(template) {
	const templates = {
		received: `Dear ${selectedCommunication.value?.public_report.reporter.name || "Resident"},\n\nThank you for submitting your report to the Mosquito Control District. We have received your communication and it has been assigned to our team for review.\n\nWe will be in touch if we need any additional information.\n\nBest regards,\nMosquito Control District`,
		scheduled: `Dear ${selectedCommunication.value?.public_report.reporter.name || "Resident"},\n\nGood news! Based on your report, we have scheduled a service visit to your area. Our technicians will be conducting mosquito control operations within the next 3-5 business days.\n\nNo action is required on your part.\n\nBest regards,\nMosquito Control District`,
		completed: `Dear ${selectedCommunication.value?.public_report.reporter.name || "Resident"},\n\nWe wanted to let you know that our team has completed mosquito control operations in your area based on your recent report.\n\nIf you continue to experience issues, please don't hesitate to submit a new report.\n\nBest regards,\nMosquito Control District`,
		need_info: `Dear ${selectedCommunication.value?.public_report.reporter.name || "Resident"},\n\nThank you for your recent report. In order to better assist you, we need some additional information:\n\n- [Specify what information is needed]\n\nPlease reply to this message with the requested details.\n\nBest regards,\nMosquito Control District`,
	};

	if (templates[template]) {
		messageText.value = templates[template];
	}
}

async function createSignal() {
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
		error.value = err.message;
		console.error("Error creating lead:", err);
	}
}

async function markInvalid() {
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

function removeCurrentFromList() {
	const index = communications.value.findIndex(
		(c) => c.id === selectedCommunication.value.id,
	);
	if (index > -1) {
		communications.value.splice(index, 1);
	}
	if (communications.value.length > 0) {
		const nextIndex = Math.min(index, communications.value.length - 1);
		selectedCommunication.value = communications.value[nextIndex];
	} else {
		selectedCommunication.value = null;
	}
}

async function sendMessage() {
	if (!messageText.value.trim()) return;

	console.log("Sending message reporter:", messageText.value);

	const payload = {
		message: messageText.value,
		reportID: selectedCommunication.value.id,
	};
	const response = await fetch(`${apiBase.value}/publicreport/message`, {
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
		`Message successfully sent to ${selectedCommunication.value.public_report.reporter.name}`,
	);
	messageText.value = "";
}

function showNotification(title, message) {
	toastTitle.value = title;
	toastMessage.value = message;
	showToast.value = true;

	setTimeout(() => {
		showToast.value = false;
	}, 3000);
}

function updateMap() {
	if (!mapRef.value) return;

	const map = mapRef.value.$el || mapRef.value;
	const loc = selectedCommunication.value.public_report.location;

	if (loc == null) {
		map.ClearMarkers();
		map.ResetCamera();
		return;
	}

	let markers = [
		new maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([loc.longitude, loc.latitude]),
	];

	let min = { lat: loc.latitude, lng: loc.longitude };
	let max = { lat: loc.latitude, lng: loc.longitude };

	for (const i of selectedCommunication.value.public_report.images) {
		if (
			i.location != null &&
			i.location.latitude != 0 &&
			i.location.longitude != 0
		) {
			markers.push(
				new maplibregl.Marker({
					color: "#00FF00",
					draggable: false,
				}).setLngLat([i.location.longitude, i.location.latitude]),
			);
			min.lat = Math.min(min.lat, i.location.latitude);
			min.lng = Math.min(min.lng, i.location.longitude);
			max.lat = Math.max(max.lat, i.location.latitude);
			max.lng = Math.max(max.lng, i.location.longitude);
		}
	}

	map.SetMarkers(markers);

	const bounds = new maplibregl.LngLatBounds(
		new maplibregl.LngLat(min.lng - 0.01, min.lat - 0.01),
		new maplibregl.LngLat(max.lng + 0.01, max.lat + 0.01),
	);

	map.FitBounds(bounds, {
		padding: 50,
	});
}
function onFilterChange(filters) {
	console.log("Filters changed");
}
// Lifecycle hooks
onMounted(async () => {
	await loadFromAPI();

	// Subscribe to SSE events
	if (window.SSEManager) {
		window.SSEManager.subscribe("*", (e) => {
			if (e.resource.startsWith("rmo:")) {
				fetchCommunications();
			}
		});
	}

	// Setup map layer after next tick to ensure map is mounted
	await nextTick();

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
});
</script>
