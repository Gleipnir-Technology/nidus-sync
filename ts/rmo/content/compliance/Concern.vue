<style scoped>
.observation-image {
	width: 100%;
	height: 200px;
	background-color: #e9ecef;
	border: 1px solid #dee2e6;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #6c757d;
	cursor: pointer;
	transition: all 0.2s ease;
	position: relative;
	overflow: hidden;
}

.observation-image:hover {
	border-color: #0d6efd;
	box-shadow: 0 2px 8px rgba(13, 110, 253, 0.2);
}

.observation-image .overlay {
	position: absolute;
	bottom: 0;
	left: 0;
	right: 0;
	background: rgba(0, 0, 0, 0.7);
	color: white;
	padding: 8px;
	font-size: 12px;
	text-align: center;
}

.inspector-notes {
	background-color: #fff;
	border-left: 4px solid #0d6efd;
	padding: 16px;
	border-radius: 4px;
}
</style>
<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="3" />

		<!-- Main Content -->
		<main>
			<h2 class="h4 mb-3">
				District observations indicate a possible mosquito problem at this
				property
			</h2>

			<p class="text-muted mb-4">
				Our inspector documented the following conditions during their visit.
				Tap any image to view details.
			</p>

			<!-- Observation Images -->
			<div class="mb-4">
				<label class="form-label fw-semibold mb-3">Inspection Photos</label>

				<div class="row g-3 mb-3">
					<div class="col-6">
						<div
							class="observation-image"
							data-bs-toggle="modal"
							data-bs-target="#imageModal1"
						>
							<span class="text-center">
								<i class="bi bi-camera" style="font-size: 2rem"></i>
								<br />Photo 1
							</span>
							<div class="overlay">Tap to view</div>
						</div>
					</div>
					<div class="col-6">
						<div
							class="observation-image"
							data-bs-toggle="modal"
							data-bs-target="#imageModal2"
						>
							<span class="text-center">
								<i class="bi bi-camera" style="font-size: 2rem"></i>
								<br />Photo 2
							</span>
							<div class="overlay">Tap to view</div>
						</div>
					</div>
				</div>

				<div class="row g-3">
					<div class="col-6">
						<div
							class="observation-image"
							data-bs-toggle="modal"
							data-bs-target="#imageModal3"
						>
							<span class="text-center">
								<i class="bi bi-camera" style="font-size: 2rem"></i>
								<br />Photo 3
							</span>
							<div class="overlay">Tap to view</div>
						</div>
					</div>
					<div class="col-6">
						<div
							class="observation-image"
							data-bs-toggle="modal"
							data-bs-target="#imageModal4"
						>
							<span class="text-center">
								<i class="bi bi-camera" style="font-size: 2rem"></i>
								<br />Photo 4
							</span>
							<div class="overlay">Tap to view</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Inspector Comments -->
			<div class="mb-4">
				<label class="form-label fw-semibold mb-3">Inspector Comments</label>
				<div class="inspector-notes">
					<div class="mb-2">
						<small class="text-muted">
							<i class="bi bi-person-badge"></i>
							Inspector #1 | 2 days ago
						</small>
					</div>
					<p class="mb-0">some fake comments here</p>
				</div>
			</div>

			<div class="alert alert-info" role="alert">
				<small>
					<i class="bi bi-info-circle"></i>
					These observations were made from outside the property. In the next
					steps, you'll have an opportunity to provide your response and any
					additional information.
				</small>
			</div>

			<!-- Navigation Buttons -->
			<div class="d-flex gap-2 mt-4">
				<RouterLink to="./address" class="btn btn-outline-secondary">
					Back
				</RouterLink>
				<RouterLink to="evidence" class="btn btn-primary flex-grow-1">
					Continue
				</RouterLink>
			</div>
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
