<style scoped>
.concern-image {
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

.concern-image:hover {
	border-color: #0d6efd;
	box-shadow: 0 2px 8px rgba(13, 110, 253, 0.2);
}

.concern-image .overlay {
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
				The following images show areas we are concerned about related to your
				property. Tap any image to view details.
			</p>

			<!-- Observation Images -->
			<div class="mb-4">
				<label class="form-label fw-semibold mb-3">Site Photos</label>

				<div class="row g-3 mb-3">
					<div class="card">
						<div class="card-body">
							<img
								class="concern-image"
								@click="openImageViewer(index)"
								:key="index"
								:src="concern.url"
								v-for="(concern, index) in modelValue.concerns"
							/>
						</div>
					</div>
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
				<RouterLink
					class="btn btn-outline-secondary"
					:to="routes.ComplianceAddress(props.publicID)"
				>
					Back
				</RouterLink>
				<ButtonLoading
					class="flex-grow-1"
					@click="doContinue"
					icon="bi-caret-right-fill"
					:loading="isUploading"
					text="Continue"
				/>
			</div>
		</main>
		<ImageViewerModal
			@close="showImageModal = false"
			:image="currentImage"
			:show="showImageModal"
		/>
	</div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import ButtonLoading from "@/components/common/ButtonLoading.vue";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ImageViewerModal, { Image } from "@/rmo/components/ImageViewerModal.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import { router } from "@/rmo/route/config";
import { useRoutes } from "@/rmo/route/use";
import type { District, PublicReportCompliance } from "@/type/api";

interface Emits {
	(e: "doConcern"): void;
	(e: "update:modelValue", value: PublicReportCompliance): void;
}
interface Props {
	district: District;
	isUploading: boolean;
	modelValue: PublicReportCompliance;
	publicID: string;
}
const currentImageIndex = ref<number>(0);
const emit = defineEmits<Emits>();
const isLoading = ref<boolean>(true);
const props = defineProps<Props>();
const routes = useRoutes();
const showImageModal = ref(false);

const currentImage = computed((): Image | undefined => {
	const i = props.modelValue.concerns[currentImageIndex.value] ?? null;
	if (!i) {
		return undefined;
	}
	return {
		src: i.url,
	};
});
function doContinue() {
	emit("update:modelValue", props.modelValue);
	emit("doConcern");
	router.push(routes.ComplianceEvidence(props.publicID));
}
async function doMounted() {
	if (props.modelValue.concerns.length == 0) {
		router.push(routes.ComplianceEvidence(props.publicID));
	}
}
function openImageViewer(index: number) {
	currentImageIndex.value = index;
	showImageModal.value = true;
}
onMounted(() => {
	doMounted();
});
</script>
