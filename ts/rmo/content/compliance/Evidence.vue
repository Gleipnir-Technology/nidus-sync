<style scoped>
.upload-area {
	border: 2px dashed #0d6efd;
	border-radius: 8px;
	padding: 32px 16px;
	text-align: center;
	background-color: #fff;
	cursor: pointer;
	transition: all 0.2s ease;
}

.upload-area:hover {
	background-color: #f8f9fa;
	border-color: #0a58ca;
}

.upload-area.dragover {
	background-color: #e7f1ff;
	border-color: #0a58ca;
}

#file-input {
	display: none;
}

.photo-preview {
	position: relative;
	width: 100%;
	height: 120px;
	border-radius: 8px;
	overflow: hidden;
	background-color: #e9ecef;
	border: 1px solid #dee2e6;
}

.photo-preview img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.photo-preview .remove-btn {
	position: absolute;
	top: 4px;
	right: 4px;
	background-color: rgba(220, 53, 69, 0.9);
	color: white;
	border: none;
	border-radius: 50%;
	width: 28px;
	height: 28px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s ease;
}

.photo-preview .remove-btn:hover {
	background-color: rgba(220, 53, 69, 1);
	transform: scale(1.1);
}

.guidance-list {
	list-style: none;
	padding-left: 0;
}

.guidance-list li {
	padding: 8px 0;
	padding-left: 28px;
	position: relative;
}

.guidance-list li::before {
	content: "✓";
	position: absolute;
	left: 0;
	color: #198754;
	font-weight: bold;
}
</style>
<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="4" />
		<!-- Main Content -->
		<main>
			<h2 class="h4 mb-3">Upload photos of the area</h2>

			<p class="text-muted mb-3">
				Please provide current photos to help us assess the situation.
			</p>

			<!-- Photo Guidance -->
			<div class="card mb-4">
				<div class="card-body">
					<h6 class="card-subtitle mb-3 text-muted">
						<i class="bi bi-lightbulb"></i> Helpful photos are:
					</h6>
					<ul class="guidance-list mb-0">
						<li>Recent (taken within the last 24 hours)</li>
						<li>Showing the specific area of concern</li>
						<li>Making water conditions clearly visible</li>
					</ul>
				</div>
			</div>

			<!-- Upload Area -->
			<div class="mb-4">
				<label class="form-label fw-semibold">Photos</label>
				<ImageUpload v-model="images" />
			</div>

			<div class="mb-4" v-if="modelValue.images.length > 0">
				<div class="alert alert-primary" role="alert">
					You've already added {{ modelValue.images.length }} image{{
						modelValue.images.length == 1 ? "" : "s"
					}}
					to this report. If you add images below, they will replace the image{{
						modelValue.images.length == 1 ? "" : "s"
					}}
					already on this report.
				</div>
			</div>

			<!-- Additional Comments -->
			<div class="mb-4">
				<label for="comments" class="form-label fw-semibold">
					Additional Comments
					<span class="text-muted fw-normal">(Optional)</span>
				</label>
				<textarea
					class="form-control"
					id="comments"
					name="comments"
					placeholder="Provide any additional information that may be helpful..."
					rows="4"
					v-model="modelValue.comments"
				></textarea>
				<div class="form-text">
					Example: "This standing water appeared after recent rain" or "I've
					already taken steps to address this issue"
				</div>
			</div>

			<!-- Navigation Buttons -->
			<div class="d-flex gap-2 mt-4">
				<RouterLink to="./address" class="btn btn-outline-secondary">
					Back
				</RouterLink>
				<button class="btn btn-primary flex-grow-1" @click="doContinue">
					Continue
				</button>
			</div>
		</main>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";

import { router } from "@/rmo/router";
import type { District, PublicReportCompliance } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ImageUpload, { Image } from "@/components/ImageUpload.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";

interface Emits {
	(e: "update:modelValue", value: PublicReportCompliance): void;
	(e: "doEvidence", images: Image[]): void;
}
interface Props {
	district: District;
	modelValue: PublicReportCompliance;
}
const emit = defineEmits<Emits>();
const images = ref<Image[]>([]);
const props = defineProps<Props>();
function doContinue() {
	emit("update:modelValue", props.modelValue);
	emit("doEvidence", images.value);
	router.push("./permission");
}
</script>
