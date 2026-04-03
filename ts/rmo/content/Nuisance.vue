<style scoped>
.district-logo {
	max-height: 80px;
	width: auto;
}
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
.map-container {
	background-color: #e9ecef;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-top: 20px;
}
#map {
	width: 100%;
	height: 100%;
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
</style>

<template>
	<main>
		<slot name="header"></slot>
		<div class="container">
			<div class="row mb-4">
				<div class="col-12">
					<h2>Report Mosquito Nuisance</h2>
					<p class="lead">Help us identify mosquito activity in your area</p>
				</div>
			</div>

			<!-- Report Form -->
			<form
				id="mosquitoNuisanceForm"
				action="/api/rmo/nuisance"
				method="POST"
				enctype="multipart/form-data"
			>
				<!-- Location Section -->
				<div class="form-section">
					<div class="section-heading">
						<i class="bi bi-geo-alt"></i>
						<h3>Nuisance Location Information</h3>
					</div>
					<div class="col-md-12">
						<div class="alert alert-info" role="info">
							<p class="mb-0">
								You can select the location by address or by moving the marker
								on the map.
							</p>
						</div>
					</div>

					<div class="col-md-6">
						<div class="mb-3 position-relative">
							<AddressSuggestion
								v-model="selectedAddress"
								placeholder="Start typing an address (min 3 characters)"
								@address-selected="doAddressSelected"
							>
							</AddressSuggestion>
						</div>
					</div>
				</div>
				<p class="small text-muted mb-2">
					You can also click on the map to mark the location precisely
				</p>
				<div class="map-container">
					<map-locator id="map"></map-locator>
				</div>
				<input type="hidden" id="map-zoom" name="map-zoom" />
				<input type="hidden" id="address-country" name="address-country" />
				<input type="hidden" id="address-locality" name="address-locality" />
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

				<!-- Mosquito Activity Section -->
				<div class="form-section">
					<div class="section-heading">
						<i class="bi bi-mosquito-color"></i>
						<h3>Mosquito Activity Information</h3>
					</div>
					<p class="mb-4">
						The time when mosquitoes are active can help us identify the species
						and likely breeding sources.
					</p>

					<!-- Time of Day -->
					<div class="row mb-4">
						<div class="col-12">
							<label class="form-label"
								>When do you typically notice mosquitoes? (Select all that
								apply)</label
							>
							<div class="row">
								<div class="col-6 col-md-3">
									<input
										type="checkbox"
										class="btn-check"
										id="earlyMorning"
										name="tod-early"
										autocomplete="off"
									/>
									<label
										class="btn btn-outline-primary time-of-day-btn"
										for="earlyMorning"
									>
										<span class="time-of-day-icon"
											><i class="bi bi-sunrise"></i
										></span>
										<span class="time-label">Early Morning</span>
										<small class="text-muted">5am-8am</small>
									</label>
								</div>
								<div class="col-6 col-md-3">
									<input
										type="checkbox"
										class="btn-check"
										id="daytime"
										name="tod-day"
										autocomplete="off"
									/>
									<label
										class="btn btn-outline-primary time-of-day-btn"
										for="daytime"
									>
										<span class="time-of-day-icon"
											><i class="bi bi-sun"></i
										></span>
										<span class="time-label">Daytime</span>
										<small class="text-muted">8am-5pm</small>
									</label>
								</div>
								<div class="col-6 col-md-3">
									<input
										type="checkbox"
										class="btn-check"
										id="evening"
										name="tod-evening"
										autocomplete="off"
									/>
									<label
										class="btn btn-outline-primary time-of-day-btn"
										for="evening"
									>
										<span class="time-of-day-icon"
											><i class="bi bi-sunset"></i
										></span>
										<span class="time-label">Evening</span>
										<small class="text-muted">5pm-9pm</small>
									</label>
								</div>
								<div class="col-6 col-md-3">
									<input
										type="checkbox"
										class="btn-check"
										id="night"
										name="tod-night"
										autocomplete="off"
									/>
									<label
										class="btn btn-outline-primary time-of-day-btn"
										for="night"
									>
										<span class="time-of-day-icon"
											><i class="bi bi-moon-stars"></i
										></span>
										<span class="time-label">Night</span>
										<small class="text-muted">9pm-5am</small>
									</label>
								</div>
							</div>
						</div>
					</div>

					<!-- Duration -->
					<div class="row mb-4">
						<div class="col-md-6">
							<label for="duration" class="form-label"
								>How long have you been experiencing this mosquito
								problem?</label
							>
							<select class="form-select" name="duration">
								<option value="just-noticed">Just noticed recently</option>
								<option value="few-days">A few days</option>
								<option value="1-2-weeks">1-2 weeks</option>
								<option value="2-4-weeks">2-4 weeks</option>
								<option value="1-3-months">1-3 months</option>
								<option value="seasonal">All season (recurring issue)</option>
							</select>
						</div>
					</div>

					<!-- Location -->
					<div class="row">
						<div class="col-md-12">
							<label for="source-location" class="form-label"
								>Where on your property do you notice the most mosquito
								activity?</label
							>
							<select
								class="form-select tall"
								multiple="true"
								name="source-location"
							>
								<option value="frontyard">Front yard</option>
								<option value="backyard">Back yard</option>
								<option value="garden">Garden</option>
								<option value="pool-area">Pool area</option>
								<option value="other">Other area</option>
							</select>
						</div>
					</div>
				</div>

				<button
					id="toggle-additional"
					class="btn btn-warning"
					type="button"
					@click="toggleCollapse()"
				>
					Click here to answer a few more questions to better help us solve your
					mosquito problem
				</button>
				<div v-show="!isCollapsed" class="collapse">
					<!-- Potential Sources Section -->
					<div class="form-section">
						<div class="section-heading">
							<i class="bi bi-search"></i>
							<h3>Potential Mosquito Sources</h3>
						</div>
						<p class="mb-3">
							Have you noticed any of these common mosquito breeding sources in
							your area?
						</p>

						<div class="card card-highlight mb-4">
							<div class="card-body">
								<h5 class="card-title">Did you know?</h5>
								<p class="card-text">
									Mosquitoes can breed in as little as a bottle cap of water!
									Eliminating standing water is the most effective way to reduce
									mosquito populations.
								</p>
							</div>
						</div>

						<div class="row g-4 mb-4">
							<!-- Source 1 -->
							<div class="col-md-4">
								<div class="card source-card">
									<div class="card-body text-center">
										<div class="source-icon">
											<i class="bi bi-water"></i>
										</div>
										<h5 class="card-title">Stagnant Water</h5>
										<p class="card-text">
											Green pools, ponds, fountains, or birdbaths that aren't
											maintained
										</p>
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												name="source-stagnant"
												id="sourceStagnantWater"
											/>
											<label class="form-check-label" for="sourceStagnantWater">
												I've noticed this
											</label>
										</div>
									</div>
								</div>
							</div>

							<!-- Source 2 -->
							<div class="col-md-4">
								<div class="card source-card">
									<div class="card-body text-center">
										<div class="source-icon">
											<i class="bi bi-droplet"></i>
										</div>
										<h5 class="card-title">Containers</h5>
										<p class="card-text">
											Buckets, planters, toys, tires, or any items that collect
											rainwater
										</p>
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												name="source-container"
												id="sourceContainers"
											/>
											<label class="form-check-label" for="sourceContainers">
												I've noticed this
											</label>
										</div>
									</div>
								</div>
							</div>

							<!-- Source 3 -->
							<div class="col-md-4">
								<div class="card source-card">
									<div class="card-body text-center">
										<div class="source-icon">
											<i class="bi bi-house"></i>
										</div>
										<h5 class="card-title">Sprinklers & Gutters</h5>
										<p class="card-text">
											Clogged street gutters, yard drains, or AC units that
											collect water
										</p>
										<div class="form-check">
											<input
												class="form-check-input"
												type="checkbox"
												name="source-gutters"
												id="sourceGutters"
											/>
											<label class="form-check-label" for="sourceGutters">
												I've noticed this
											</label>
										</div>
									</div>
								</div>
							</div>
						</div>

						<div class="row">
							<div class="col-md-12">
								<label for="otherSources" class="form-label"
									>Have you noticed any other potential mosquito breeding
									sources?</label
								>
								<textarea
									class="form-control"
									id="otherSources"
									name="source-description"
									rows="2"
									placeholder="Describe any other potential breeding sites you've noticed..."
								></textarea>
								<photo-upload />
							</div>
						</div>
					</div>

					<div class="form-section">
						<div class="section-heading">
							<i class="bi bi-card-text"></i>
							<h3>Additional Information</h3>
						</div>

						<div class="row">
							<div class="col-md-12">
								<label for="additionalInfo" class="form-label"
									>Is there anything else you'd like us to know?</label
								>
								<textarea
									class="form-control"
									id="additionalInfo"
									maxlength="2048"
									name="additional-info"
									rows="4"
									placeholder="Additional details about the mosquito issue..."
								></textarea>
							</div>
						</div>
					</div>
				</div>
				<!-- Submit Section -->
				<div class="submit-container">
					<div class="row align-items-center">
						<div class="col-md-8">
							<p class="mb-0">
								<strong>Thank you for reporting this mosquito issue.</strong>
							</p>
							<p class="mb-0 small text-muted">
								After submission, you'll receive a confirmation with a report ID
								and further information.
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
</template>

<script setup lang="ts">
import { ref } from "vue";
import AddressSuggestion from "@/components/AddressSuggestion.vue";
import { Address } from "@/type/stadia";

const isCollapsed = ref<boolean>(true);
const toggleCollapse = () => {
	isCollapsed.value = !isCollapsed.value;
};
const selectedAddress = ref<Address | null>(null);
function doAddressSelected(address: Address) {
	console.log("Address selected", address);
}
</script>
