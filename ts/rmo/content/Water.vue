<style scoped>
.form-section {
	margin-bottom: 2.5rem;
	padding-bottom: 2rem;
	border-bottom: 1px solid #dee2e6;
}
.form-section:last-child {
	border-bottom: none;
	margin-bottom: 1rem;
	padding-bottom: 0;
}
.section-heading {
	margin-bottom: 1.5rem;
	display: flex;
	align-items: center;
}
.section-heading i {
	margin-right: 10px;
	font-size: 1.5rem;
	color: #0d6efd;
}
.submit-container {
	background-color: #f8f9fa;
	padding: 20px;
	border-radius: 5px;
	margin-top: 2rem;
}
.source-card {
	height: 100%;
	transition: transform 0.3s;
}
.source-card:hover {
	transform: translateY(-5px);
	box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}
.source-icon {
	font-size: 2rem;
	margin-bottom: 1rem;
	color: #0d6efd;
}
.time-of-day-btn {
	width: 100%;
	margin-bottom: 10px;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 15px 0;
}
.time-of-day-icon {
	font-size: 1.5rem;
	margin-bottom: 8px;
}
.time-label {
	font-size: 0.9rem;
}
select.tall {
	height: 160px;
}
.severity-item {
	text-align: center;
	padding: 10px;
}
.severity-scale {
	display: flex;
	justify-content: space-between;
	margin: 20px 0;
}
.btn-check:checked + .btn.time-of-day-btn {
	background-color: $info;
	color: white;
}
.inspection-type-card {
	cursor: pointer;
	border: 1px solid #dee2e6;
	padding: 20px;
	border-radius: 5px;
	height: 100%;
	transition: all 0.3s;
}
.inspection-type-card.selected {
	border-color: #0d6efd;
	background-color: rgba(13, 110, 253, 0.05);
}
.inspection-type-card:hover {
	border-color: #0d6efd;
}
.card-highlight {
	border-left: 4px solid #0d6efd;
	background-color: #f8f9fa;
}
.map-container {
	background-color: #e9ecef;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 20px;
	margin-top: 20px;
	/* Prevent touch scrolling issues */
	touch-action: pan-y pinch-zoom;
}
#map {
	width: 100%;
	height: 100%;
}

/* Mobile-specific adjustments */
@media (max-width: 768px) {
	.map-container {
		height: 400px;
		margin-bottom: 15px;
		margin-top: 15px;
	}
}

/* Extra small devices */
@media (max-width: 576px) {
	.map-container {
		height: 350px;
		border-radius: 5px;
	}
}
</style>
<template>
	<main class="py-5">
		<slot name="header"></slot>
		<div class="container">
			<!-- Page Title -->
			<div class="row mb-4">
				<div class="col-12">
					<h2>Report Standing Water</h2>
					<p class="lead">
						Help us locate and treat potential mosquito production sources in
						your area
					</p>
				</div>
			</div>

			<!-- Report Form -->
			<form
				enctype="multipart/form-data"
				ref="formElement"
				@submit.prevent="doSubmit"
			>
				<!-- Photo Upload Section -->
				<div class="form-section">
					<div class="section-heading">
						<i class="bi bi-camera"></i>
						<h3>Photos</h3>
					</div>
					<p class="mb-3">
						Photos help us identify the severity of the issue and may contain
						location data that can help us find the production source.
					</p>
					<div class="mb-4">
						<ImageUpload v-model="images" />
					</div>
				</div>

				<!-- Additional Information Section -->
				<div class="form-section">
					<div class="section-heading">
						<i class="bi bi-card-text"></i>
						<h3>Additional Information</h3>
					</div>
					<p class="mb-3">
						Please provide any other information that might help us address this
						mosquito production source.
					</p>

					<div class="row">
						<div class="col-md-12">
							<label for="comments" class="form-label"
								>Additional Details</label
							>
							<textarea
								class="form-control"
								id="comments"
								maxlength="2048"
								name="comments"
								rows="4"
								placeholder="Example: The house appears to be vacant. There is algae growth in the pool. I've noticed increased mosquito activity in the evenings."
							></textarea>
						</div>
					</div>
				</div>

				<!-- Location Section -->
				<div class="form-section">
					<div class="section-heading">
						<i class="bi bi-geo-alt"></i>
						<h3>Location</h3>
					</div>
					<p class="mb-3">
						Please provide the location of the potential mosquito production
						source. We may be able to extract this information from your photos
						if they contain location data.
					</p>
					<div class="col-md-12">
						<div class="alert alert-info" role="info">
							<p class="mb-0">
								You can select the location by address or by moving the marker
								on the map.
							</p>
						</div>
					</div>

					<div class="row mb-3">
						<!-- Hidden fields for location data -->
						<input type="hidden" id="address-country" name="address-country" />
						<input
							type="hidden"
							id="address-locality"
							name="address-locality"
						/>
						<input type="hidden" id="address-number" name="address-number" />
						<input
							type="hidden"
							id="address-postalcode"
							name="address-postalcode"
						/>
						<input type="hidden" id="address-region" name="address-region" />
						<input type="hidden" id="address-street" name="address-street" />
						<input type="hidden" id="latitude" name="latitude" />
						<input type="hidden" id="longitude" name="longitude" />
						<input
							type="hidden"
							id="latlng-accuracy-type"
							name="latlng-accuracy-type"
						/>
						<input
							type="hidden"
							id="latlng-accuracy-value"
							name="latlng-accuracy-value"
						/>
					</div>

					<p class="small text-muted mb-2">
						You can also click on the map to mark the location precisely
					</p>
					<AddressAndMapLocator v-model="address" />
				</div>

				<button
					class="btn btn-warning"
					@click="showMore = true"
					id="toggle-additional"
					type="button"
					v-if="!showMore"
				>
					Click here to answer a few more questions to better help us solve your
					mosquito problem
				</button>
				<div :class="{ collapse: !showMore }" id="collapse-additional-fields">
					<!-- Source Details Section -->
					<div class="form-section">
						<div class="section-heading">
							<i class="bi bi-water"></i>
							<h3>Source Details</h3>
						</div>

						<div class="row mb-4">
							<div class="col-md-6">
								<label for="duration" class="form-label"
									>How long has this production source been present?</label
								>
								<select
									class="form-select"
									id="duration"
									name="source-duration"
								>
									<option value="none">I don't know</option>
									<option value="less-than-week">Less than a week</option>
									<option value="1-2-weeks">1-2 weeks</option>
									<option value="2-4-weeks">2-4 weeks</option>
									<option value="1-3-months">1-3 months</option>
									<option value="more-than-3-months">More than 3 months</option>
								</select>
							</div>

							<div class="col-md-6">
								<label class="form-label d-block"
									>Have you observed any of the following?
									<a
										href="#"
										data-bs-toggle="modal"
										data-bs-target="#larvaeInfoModal"
										><i class="bi bi-question-circle small ms-1"></i></a
								></label>
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										id="larvae"
										name="has-larvae"
									/>
									<label class="form-check-label" for="larvae">
										Larvae (wigglers) in water
									</label>
								</div>
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										id="pupae"
										name="has-pupae"
									/>
									<label class="form-check-label" for="pupae">
										Pupae (tumblers) in water
									</label>
								</div>
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										id="adult"
										name="has-adult"
									/>
									<label class="form-check-label" for="adult">
										Adult mosquitoes near the source
									</label>
								</div>
							</div>
						</div>
					</div>

					<!-- Access Information Section -->
					<div class="form-section">
						<div class="section-heading">
							<i class="bi bi-unlock"></i>
							<h3>Access Information</h3>
						</div>
						<p class="mb-3">
							Please provide any details about how to access the mosquito
							source. This helps our technicians when they visit the site.
						</p>

						<div class="row mb-3">
							<div class="col-md-12">
								<label for="access-comments" class="form-label"
									>How can the source be accessed?</label
								>
								<textarea
									class="form-control"
									id="access-comments"
									maxlength="1024"
									name="access-comments"
									rows="3"
									placeholder="Example: The pool is in the backyard, which can be accessed through a side gate on the right side of the house."
								></textarea>
							</div>
						</div>

						<div class="row mb-3">
							<div class="col-md-12">
								<label class="form-label d-block"
									>Access obstacles (check all that apply):</label
								>
								<div class="row">
									<div class="col-md-4">
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												id="gate"
												name="access-gate"
											/>
											<label class="form-check-label" for="gate">Gate</label>
										</div>
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												id="fence"
												name="access-fence"
											/>
											<label class="form-check-label" for="fence">Fence</label>
										</div>
									</div>
									<div class="col-md-4">
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												id="locked"
												name="access-locked"
											/>
											<label class="form-check-label" for="locked"
												>Locked entrance</label
											>
										</div>
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												id="dogs"
												name="access-dog"
											/>
											<label class="form-check-label" for="dogs"
												>Dogs/pets</label
											>
										</div>
									</div>
									<div class="col-md-4">
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												id="access-other"
												name="access-other"
											/>
											<label class="form-check-label" for="access-other"
												>Other obstacle</label
											>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Contact Information Sections -->
					<div class="form-section">
						<div class="section-heading">
							<i class="bi bi-person-lines-fill"></i>
							<h3>Property Owner Information (if known)</h3>
						</div>

						<div class="row mb-4">
							<div class="col-md-4 mb-3">
								<label for="owner-name" class="form-label">Owner Name</label>
								<input
									type="text"
									class="form-control"
									id="owner-name"
									maxlength="100"
									name="owner-name"
								/>
							</div>
							<div class="col-md-4 mb-3">
								<label for="owner-phone" class="form-label">Owner Phone</label>
								<input
									type="tel"
									class="form-control"
									id="owner-phone"
									maxlength="50"
									name="owner-phone"
								/>
							</div>
							<div class="col-md-4 mb-3">
								<label for="owner-email" class="form-label">Owner Email</label>
								<input
									type="email"
									class="form-control"
									id="owner-email"
									maxlength="100"
									name="owner-email"
								/>
							</div>
						</div>
						<div class="row mb-4">
							<div class="col-md-6 mb-3 row">
								<div class="form-check mt-4">
									<input
										class="form-check-input"
										id="property-ownership"
										name="property-ownership"
										type="checkbox"
									/>
									<label class="form-check-label" for="property-ownership"
										>This is my property</label
									>
								</div>
								<div class="form-check mt-4">
									<input
										class="form-check-input"
										id="backyard-permission"
										name="backyard-permission"
										type="checkbox"
									/>
									<label class="form-check-label" for="backyard-permission"
										>I grant permission to enter the back yard of this
										property.</label
									>
								</div>
								<div class="form-check mt-4">
									<input
										class="form-check-input"
										id="reporter-confidential"
										name="reporter-confidential"
										type="checkbox"
									/>
									<label class="form-check-label" for="reporter-confidential">
										<Tooltip
											placement="top"
											title="We share your information with mosquito control districts so they can follow up with any questions they may have about your report. Check this box if you would like the district to be careful not to share your information outside of the district operations team."
										>
											<i class="bi bi-info-circle-fill text-primary ms-1"></i>
										</Tooltip>
										I would like my personal information kept
										confidential.</label
									>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- Submit Section -->
				<div class="submit-container">
					<div class="row align-items-center">
						<div class="col-md-8">
							<p class="mb-0">
								<strong
									>Thank you for helping us keep our community safe from
									mosquito-borne illnesses.</strong
								>
							</p>
							<p class="mb-0 small text-muted">
								After submission, you will receive a confirmation with a report
								ID for tracking purposes.
							</p>
						</div>
						<div class="col-md-4 text-md-end mt-3 mt-md-0">
							<button type="submit" class="btn btn-primary btn-lg">
								Submit Report
							</button>
						</div>
					</div>
				</div>
			</form>
		</div>
	</main>

	<!-- Larvae Info Modal -->
	<div
		class="modal fade"
		id="larvaeInfoModal"
		tabindex="-1"
		aria-labelledby="larvaeInfoModalLabel"
		aria-hidden="true"
	>
		<div class="modal-dialog modal-lg">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="larvaeInfoModalLabel">
						How to Identify Mosquito Larvae and Pupae
					</h5>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
					></button>
				</div>
				<div class="modal-body">
					<div class="row mb-4">
						<div class="col-md-6">
							<h6>Mosquito Larvae (Wigglers)</h6>
							<p>Mosquito larvae, often called "wigglers," are:</p>
							<ul>
								<li>Small, worm-like aquatic organisms</li>
								<li>Usually 1/4 to 1/2 inch long</li>
								<li>Move with a wiggling motion in water</li>
								<li>Hang upside-down at the water surface to breathe</li>
								<li>Visible to the naked eye in standing water</li>
							</ul>
						</div>
						<div class="col-md-6">
							<h6>Mosquito Pupae (Tumblers)</h6>
							<p>Mosquito pupae, often called "tumblers," are:</p>
							<ul>
								<li>Comma-shaped organisms</li>
								<li>Typically darker than larvae</li>
								<li>Move with a tumbling motion when disturbed</li>
								<li>Rest at the water surface</li>
								<li>The stage just before adult mosquitoes emerge</li>
							</ul>
						</div>
					</div>
					<p>
						When looking for mosquito larvae and pupae, check standing water
						sources like:
					</p>
					<ul>
						<li>Swimming pools</li>
						<li>Bird baths</li>
						<li>Buckets or containers</li>
						<li>Drainage ditches</li>
						<li>Plant saucers</li>
						<li>Rain gutters</li>
					</ul>
					<p>
						If you see small creatures moving in standing water, there's a good
						chance they're mosquito larvae or pupae.
					</p>
					<div class="text-center">
						<a href="#" class="btn btn-outline-primary"
							>View Detailed Identification Guide</a
						>
					</div>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						data-bs-dismiss="modal"
					>
						Close
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import ImageUpload, { Image } from "@/components/ImageUpload.vue";
import Tooltip from "@/components/Tooltip.vue";
import { useGeocodeStore } from "@/store/geocode";
import { useStoreLocal } from "@/store/local";
import { useStoreLocation } from "@/store/location";
import { useStorePublicReport } from "@/store/publicreport";
import type { Marker } from "@/types";
import {
	Address,
	type Geocode,
	type GeocodeSuggestion,
	type Location,
	type PublicReport,
} from "@/type/api";
import type { Camera } from "@/type/map";

const address = ref<Address>(new Address());
const currentCamera = ref<Camera | null>(null);
const currentLocation = ref<Location | null>(null);
const errorMessage = ref("");
const formElement = ref<HTMLFormElement | null>(null);
const geocode = useGeocodeStore();
const images = ref<Image[]>([]);
const isSubmitting = ref(false);
const marker = ref<Marker | null>(null);
const markers = computed((): Marker[] => {
	if (marker.value) {
		return [marker.value];
	} else {
		return [];
	}
});
const storeLocal = useStoreLocal();
const storeLocation = useStoreLocation();
const router = useRouter();
const showMore = ref<boolean>(false);
const storePublicReport = useStorePublicReport();
async function doSubmit() {
	if (!formElement.value) return;

	isSubmitting.value = true;
	errorMessage.value = "";
	try {
		const client_id = storeLocal.getClientID();
		const formData = new FormData(formElement.value);
		formData.append("client_id", client_id);
		if (address.value) {
			formData.append("address.gid", address.value.gid);
			formData.append("address.raw", address.value.raw);
			if (address.value.location) {
				formData.append(
					"address.location.latitude",
					address.value.location.latitude.toString(),
				);
				formData.append(
					"address.location.longitude",
					address.value.location.longitude.toString(),
				);
			}
		}
		formData.append(
			"location.accuracy",
			currentLocation.value?.accuracy?.toString() ?? "0",
		);
		formData.append(
			"location.latitude",
			currentLocation.value?.latitude.toString() ?? "0",
		);
		formData.append(
			"location.longitude",
			currentLocation.value?.longitude.toString() ?? "0",
		);
		images.value.forEach((image, index) => {
			formData.append(`image[${index}]`, image.file, image.name);
		});
		const resp = await fetch("/api/rmo/water", {
			method: "POST",
			body: formData,
			// Don't set Content-Type, the borwser should do it
		});
		const data: PublicReport = (await resp.json()) as PublicReport;
		storePublicReport.add(data);
		router.push("/submitted/" + data.public_id);
	} catch (error) {
		errorMessage.value =
			error instanceof Error ? error.message : "Upload failed";
	} finally {
		isSubmitting.value = false;
	}
}
onMounted(() => {
	storeLocation
		.get()
		.then((loc: GeolocationPosition) => {
			console.log("user geolocation", loc);
			const coords = loc.coords;
			currentLocation.value = coords;
			currentCamera.value = {
				location: coords,
				zoom: 15,
			};
		})
		.catch((e) => {
			console.log("failed to get location", e);
		});
});
</script>
