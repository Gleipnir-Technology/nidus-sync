<style scoped>
.map-placeholder {
	width: 100%;
	height: 250px;
	background-color: #e9ecef;
	border: 2px dashed #6c757d;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #6c757d;
	font-size: 14px;
}

#address-input {
	font-size: 16px;
}
</style>
<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="2" />
		<main>
			<h2 class="h4 mb-3">Confirm the property address</h2>

			<p class="text-muted mb-4">
				Please enter the address so we can match your response with our records.
			</p>

			<form id="address-form" method="POST" action="/compliance/address">
				<div class="mb-4">
					<label for="address-input" class="form-label fw-semibold">
						Property Address
					</label>
					<input
						type="text"
						class="form-control form-control-lg"
						id="address-input"
						name="address"
						placeholder="Start typing your address..."
						autocomplete="off"
						data-autocomplete="true"
						required
					/>
					<div class="form-text">
						Begin typing and select your address from the suggestions
					</div>
				</div>

				<!-- Map Placeholder -->
				<div class="mb-4">
					<label class="form-label fw-semibold">Location Preview</label>
					<div class="map-placeholder" id="map-container">
						<span>Map will appear here after address is selected</span>
					</div>
				</div>

				<!-- Navigation Buttons -->
				<div class="d-flex gap-2 mt-4">
					<RouterLink class="btn btn-outline-secondary" to="../compliance">
						Back
					</RouterLink>
					<RouterLink class="btn btn-primary flex-grow-1" to="./concern">
						Continue
					</RouterLink>
				</div>
			</form>
		</main>
	</div>
</template>
<script setup lang="ts">
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
interface Props {
	district: District;
}
const props = defineProps<Props>();
</script>
