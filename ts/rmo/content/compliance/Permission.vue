<style scoped>
.access-option {
	background-color: #fff;
	border: 2px solid #dee2e6;
	border-radius: 8px;
	padding: 16px;
	margin-bottom: 12px;
	cursor: pointer;
	transition: all 0.2s ease;
}

.access-option:hover {
	border-color: #0d6efd;
	background-color: #f8f9fa;
}

.access-option.selected {
	border-color: #0d6efd;
	background-color: #e7f1ff;
}

.access-option input[type="radio"] {
	width: 20px;
	height: 20px;
	margin-right: 12px;
	cursor: pointer;
}

.access-option label {
	cursor: pointer;
	margin-bottom: 0;
	flex-grow: 1;
}

.conditional-section {
	display: block;
	margin-top: 16px;
	padding: 16px;
	background-color: #fff;
	border-radius: 8px;
	border-left: 4px solid #0d6efd;
}

.dog-warning {
	background-color: #fff3cd;
	border-left: 4px solid #ffc107;
	padding: 12px;
	border-radius: 4px;
	margin-top: 12px;
}

.encouragement-box {
	background-color: #d1ecf1;
	border-left: 4px solid #0dcaf0;
	padding: 16px;
	border-radius: 4px;
}
</style>
<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="5" />

		<main>
			<h2 class="h4 mb-3">Property access permission</h2>

			<p class="text-muted mb-4">
				Granting access allows our technicians to inspect and potentially treat
				mosquito sources more quickly, helping protect you and your neighbors.
			</p>

			<form id="access-form" method="POST" action="/compliance/access">
				<!-- Access Options -->
				<div class="mb-4">
					<label class="form-label fw-semibold mb-3"
						>Please select an option:</label
					>

					<!-- Option 1: Enter without owner present -->
					<div class="access-option">
						<div class="d-flex align-items-start">
							<input
								type="radio"
								id="access-allowed"
								:value="Permission.GRANTED"
								v-model="selectedPermission"
								class="mt-1"
							/>
							<label for="access-allowed">
								<div class="fw-semibold">
									A technician may enter even if I am not home
								</div>
								<small class="text-muted">Fastest resolution</small>
							</label>
						</div>
					</div>

					<!-- Conditional fields for Option 1 -->
					<div
						v-if="
							selectedPermission && selectedPermission == Permission.GRANTED
						"
						class="conditional-section"
					>
						<div class="mb-3">
							<label for="access-instructions" class="form-label">
								Access Instructions
								<span class="text-muted">(Optional)</span>
							</label>
							<textarea
								class="form-control"
								id="access-instructions"
								name="access_instructions"
								rows="3"
								placeholder="Example: Side gate on left, backyard near shed..."
							></textarea>
						</div>

						<div class="mb-3">
							<label for="gate-code" class="form-label">
								Gate Code
								<span class="text-muted">(Optional)</span>
							</label>
							<input
								type="text"
								class="form-control"
								id="gate-code"
								name="gate_code"
								placeholder="Enter code if applicable"
							/>
						</div>

						<div class="mb-3">
							<label class="form-label">Dog on Property?</label>
							<div class="form-check">
								<input
									class="form-check-input"
									type="radio"
									name="has_dog"
									id="dog-no"
									value="no"
									checked
									onchange="toggleDogWarning()"
								/>
								<label class="form-check-label" for="dog-no"> No </label>
							</div>
							<div class="form-check">
								<input
									class="form-check-input"
									id="dog-yes"
									type="radio"
									name="has_dog"
									v-model="hasDog"
									onchange="toggleDogWarning()"
								/>
								<label class="form-check-label" for="dog-yes"> Yes </label>
							</div>
						</div>

						<div class="dog-warning" v-if="hasDog">
							<small>
								<i class="bi bi-exclamation-triangle"></i>
								<strong>Important:</strong> Our staff will only enter if the dog
								is secured indoors. Please ensure your pet is safely inside
								before a technician arrives.
							</small>
						</div>
					</div>

					<!-- Option 2: Enter with owner present -->
					<div class="access-option">
						<div class="d-flex align-items-start">
							<input
								type="radio"
								id="access-with-owner"
								:value="Permission.WITH_OWNER"
								v-model="selectedPermission"
								class="mt-1"
							/>
							<label for="access-with-owner">
								<div class="fw-semibold">
									A technician may enter, but I want to be present
								</div>
								<small class="text-muted">Requires scheduling</small>
							</label>
						</div>
					</div>

					<!-- Conditional fields for Option 2 -->
					<div
						class="conditional-section"
						v-if="
							selectedPermission && selectedPermission == Permission.WITH_OWNER
						"
					>
						<div class="form-check mb-3">
							<input
								class="form-check-input"
								type="checkbox"
								id="request-scheduled"
								name="request_scheduled"
								value="yes"
							/>
							<label class="form-check-label" for="request-scheduled">
								I would like to request a scheduled visit
							</label>
						</div>

						<div class="mb-3">
							<label for="availability-notes" class="form-label">
								Availability / Access Notes
								<span class="text-muted">(Optional)</span>
							</label>
							<textarea
								class="form-control"
								id="availability-notes"
								name="availability_notes"
								rows="3"
								placeholder="Example: Available weekday mornings, please call before visiting..."
							></textarea>
						</div>
					</div>

					<!-- Option 3: Not granting entry -->
					<div class="access-option">
						<div class="d-flex align-items-start">
							<input
								type="radio"
								id="access-denied"
								:value="Permission.DENIED"
								v-model="selectedPermission"
								class="mt-1"
							/>
							<label for="access-denied">
								<div class="fw-semibold">
									I am not granting entry at this time
								</div>
								<small class="text-muted">May require follow-up</small>
							</label>
						</div>
					</div>

					<!-- Conditional message for Option 3 -->
					<div
						class="conditional-section"
						v-if="selectedPermission && selectedPermission == Permission.DENIED"
					>
						<div class="encouragement-box">
							<p class="mb-2">
								<strong>We understand.</strong> Your cooperation is voluntary,
								but mosquito breeding sources can affect the health and comfort
								of the entire community.
							</p>
							<p class="mb-2">
								To help us review this situation and avoid unnecessary
								escalation, we strongly encourage you to:
							</p>
							<ul class="mb-2">
								<li>Provide detailed photos of the area</li>
								<li>Share your contact information</li>
								<li>Include any context that may be helpful</li>
							</ul>
							<p class="mb-0">
								<small class="text-muted">
									This allows our team to assess whether the concern has been
									addressed or if additional steps may be necessary.
								</small>
							</p>
						</div>
					</div>
				</div>

				<!-- Navigation Buttons -->
				<div class="d-flex gap-2 mt-4">
					<RouterLink class="btn btn-outline-secondary" to="./evidence">
						Back
					</RouterLink>
					<RouterLink class="btn btn-primary flex-grow-1" to="contact">
						Continue
					</RouterLink>
				</div>
			</form>
		</main>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
interface Props {
	district: District;
}
enum Permission {
	GRANTED = "access-granted",
	WITH_OWNER = "access-with-owner",
	DENIED = "access-denied",
}
const props = defineProps<Props>();
const hasDog = ref<boolean>(false);
const selectedPermission = ref<Permission | undefined>();

function toggleDogWarning() {
	/*
			const dogYes = document.getElementById("dog-yes");
			const dogWarning = document.getElementById("dog-warning");

			if (dogYes.checked) {
				dogWarning.style.display = "block";
			} else {
				dogWarning.style.display = "none";
			}
*/
}
</script>
