<style scoped>
.modal.show {
	background-color: rgba(0, 0, 0, 0.5);
}
</style>
<template>
	<div
		class="modal fade"
		:class="{ 'show d-block': show }"
		tabindex="-1"
		v-show="show"
		@click.self="emit('close')"
	>
		<div class="modal-dialog modal-lg modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">
						Photo {{ currentPhotoIndex + 1 }} of
						{{ images.length || 0 }}
					</h5>
					<button
						type="button"
						class="btn-close"
						@click="emit('close')"
					></button>
				</div>
				<div class="modal-body text-center">
					<div v-if="images && show">
						<img
							:src="images[currentPhotoIndex].url_content"
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
										{{ images[currentPhotoIndex].exif?.created || "N/A" }}
									</span>
								</div>
								<div class="col-md-4">
									<small class="text-muted d-block">Camera</small>
									<span>
										{{
											(images[currentPhotoIndex].exif?.make || "") +
												" " +
												(images[currentPhotoIndex].exif?.model || "") || "N/A"
										}}
									</span>
								</div>
								<div class="col-md-4">
									<small class="text-muted d-block"
										>Distance from Reporter</small
									>
									<span v-if="images[currentPhotoIndex].location != null">
										{{
											formatDistance(
												images[currentPhotoIndex].distance_from_reporter_meters,
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
						@click="emit('imagePrevious')"
						:disabled="currentPhotoIndex === 0"
					>
						<i class="bi bi-chevron-left"></i> Previous
					</button>
					<button
						class="btn btn-outline-secondary"
						@click="emit('imageNext')"
						:disabled="currentPhotoIndex >= (images?.length || 1) - 1"
					>
						Next <i class="bi bi-chevron-right"></i>
					</button>
				</div>
			</div>
		</div>
	</div>
	<div
		class="modal-backdrop fade show"
		v-show="show"
		@click="emit('close')"
	></div>
</template>

<script setup lang="ts">
interface Emits {
	(e: "close"): void;
	(e: "imageNext"): void;
	(e: "imagePrevious"): void;
}
interface Props {
	currentPhotoIndex: int | null;
	images: Photo[] | null;
	show: boolean;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
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
</script>
