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
					<a class="btn btn-outline-secondary" href="../compliance"> Back </a>
					<a class="btn btn-primary flex-grow-1" href="concern"> Continue </a>
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
