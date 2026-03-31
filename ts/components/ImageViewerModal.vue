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
						Image {{ currentImageIndex + 1 }} of
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
							:src="images[currentImageIndex].url_content"
							class="img-fluid rounded"
							style="max-height: 60vh"
						/>

						<!-- EXIF Data Section -->
						<div class="mt-4 pt-3 border-top text-start">
							<h6 class="text-muted mb-3">Image Information</h6>
							<div class="row g-3">
								<div class="col-md-4">
									<small class="text-muted d-block">Date Taken</small>
									<span>
										{{ images[currentImageIndex].exif?.created || "N/A" }}
									</span>
								</div>
								<div class="col-md-4">
									<small class="text-muted d-block">Camera</small>
									<span>
										{{
											(images[currentImageIndex].exif?.make || "") +
												" " +
												(images[currentImageIndex].exif?.model || "") || "N/A"
										}}
									</span>
								</div>
								<div class="col-md-4">
									<small class="text-muted d-block"
										>Distance from Reporter</small
									>
									<span v-if="images[currentImageIndex].location != null">
										{{
											formatDistance(
												images[currentImageIndex].distance_from_reporter_meters,
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
						:disabled="currentImageIndex === 0"
					>
						<i class="bi bi-chevron-left"></i> Previous
					</button>
					<button
						class="btn btn-outline-secondary"
						@click="emit('imageNext')"
						:disabled="currentImageIndex >= (images?.length || 1) - 1"
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
import { formatDistance } from "@/format";
import { Image } from "@/types";

interface Emits {
	(e: "close"): void;
	(e: "imageNext"): void;
	(e: "imagePrevious"): void;
}
interface Props {
	currentImageIndex: number;
	images: Image[];
	show: boolean;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
</script>
