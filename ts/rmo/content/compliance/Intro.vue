<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<div class="mb-4">
			<div class="d-flex justify-content-between align-items-center mb-2">
				<span class="small text-muted">Step 1 of 8</span>
			</div>
			<div class="progress" style="height: 8px">
				<div
					class="progress-bar"
					role="progressbar"
					style="width: 12%"
					aria-valuenow="12"
					aria-valuemin="0"
					aria-valuemax="100"
				></div>
			</div>
		</div>

		<!-- Main Content -->
		<main>
			<h2 class="h4 mb-3">Action requested for this property</h2>

			<div class="alert alert-warning" role="alert">
				<p class="mb-2">
					The {{ district.name }} has identified a possible mosquito breeding
					source at this property during a recent inspection.
				</p>
				<p class="mb-0">
					Please confirm the property address and provide any photos, access
					details, and contact information to help us review and resolve this
					issue quickly.
				</p>
			</div>

			<div class="d-grid mt-4">
				<RouterLink class="btn btn-primary btn-lg" to="./address">
					Get Started</RouterLink
				>
			</div>
		</main>
	</div>
</template>
<script setup lang="ts">
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
interface Props {
	district: District;
}
const props = defineProps<Props>();
</script>
