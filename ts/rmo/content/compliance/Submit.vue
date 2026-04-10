<style scoped>
.summary-section {
	background-color: #fff;
	border-radius: 8px;
	padding: 16px;
	margin-bottom: 16px;
	border: 1px solid #dee2e6;
}

.summary-section h3 {
	font-size: 1rem;
	font-weight: 600;
	margin-bottom: 12px;
	color: #495057;
}

.summary-item {
	padding: 8px 0;
	border-bottom: 1px solid #f8f9fa;
}

.summary-item:last-child {
	border-bottom: none;
}

.summary-label {
	font-weight: 500;
	color: #6c757d;
	font-size: 0.875rem;
}

.summary-value {
	color: #212529;
	margin-top: 4px;
}

.photo-count {
	display: inline-flex;
	align-items: center;
	gap: 8px;
	background-color: #e7f1ff;
	padding: 4px 12px;
	border-radius: 16px;
	font-size: 0.875rem;
}

.status-badge {
	display: inline-block;
	padding: 4px 12px;
	border-radius: 16px;
	font-size: 0.875rem;
	font-weight: 500;
}

.status-provided {
	background-color: #d1e7dd;
	color: #0f5132;
}

.status-not-provided {
	background-color: #f8d7da;
	color: #842029;
}

.encouragement-box {
	background-color: #d1ecf1;
	border-left: 4px solid #0dcaf0;
	padding: 16px;
	border-radius: 4px;
}

.validation-alert {
	background-color: #fff3cd;
	border-left: 4px solid #ffc107;
	padding: 16px;
	border-radius: 4px;
}

.submit-btn {
	font-size: 1.125rem;
	padding: 14px;
}
</style>

<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="8" />

		<main>
			<h2 class="h4 mb-3">Review and submit your response</h2>

			<div class="encouragement-box mb-4">
				<p class="mb-2">
					<strong><i class="bi bi-lightbulb"></i> Before you submit</strong>
				</p>
				<p class="mb-0 small">
					Providing photos, access permissions, and contact information gives
					the District the best opportunity to review your response and
					potentially close this matter without further action. The more detail
					you provide, the better we can assess the situation.
				</p>
			</div>

			<!-- Response Summary -->
			<div class="mb-4">
				<h3 class="h6 mb-3 text-muted">Your Response Summary</h3>

				<!-- Address -->
				<div class="summary-section">
					<h3><i class="bi bi-geo-alt"></i> Property Address</h3>
					<div class="summary-item">
						<div class="summary-value" v-if="modelValue.address.raw">
							{{ modelValue.address.raw }}
							<span class="status-badge status-provided ms-2">
								<i class="bi bi-check-circle"></i> Provided
							</span>
						</div>
						<div class="summary-value" v-else>
							<span class="status-badge status-not-provided ms-2">
								<i class="bi bi-x-circle"></i> Not Provided
							</span>
						</div>
					</div>
				</div>

				<!-- Photos -->
				<div class="summary-section">
					<h3><i class="bi bi-camera"></i> Photos</h3>
					<div class="summary-item">
						<div class="summary-value">
							<span class="photo-count" v-if="modelValue.images.length > 0">
								<i class="bi bi-images"></i>
								{{ modelValue.images.length }} photo{{
									modelValue.images.length > 1 ? "s" : ""
								}}
								uploaded
							</span>
							<span class="photo-count status-badge status-not-provided" v-else>
								<i class="bi bi-x-circle"></i> Not Provided
							</span>
						</div>
						<div class="summary-value mt-2" v-if="modelValue.comments">
							<div class="summary-label">Comments:</div>
							<small class="text-muted">{{ modelValue.comments }}</small>
						</div>
					</div>
				</div>

				<!-- Access Permission -->
				<div class="summary-section">
					<h3><i class="bi bi-door-open"></i> Property Access</h3>
					<div class="summary-item">
						<div class="summary-value">
							<span
								class="status-badge status-provided"
								v-if="modelValue.access == PermissionAccess.GRANTED"
							>
								<i class="bi bi-check-circle"></i> Entry permitted without owner
								present
							</span>
							<span
								class="status-badge status-provided"
								v-else-if="modelValue.access == PermissionAccess.WITH_OWNER"
							>
								<i class="bi bi-check-circle"></i> Entry permitted with owner
								present
							</span>
							<span
								class="status-badge status-not-provided"
								v-else-if="modelValue.access == PermissionAccess.DENIED"
							>
								<i class="bi bi-x-circle"></i> Entry denied
							</span>
							<span class="status-badge status-not-provided" v-else>
								<i class="bi bi-x-circle"></i> Not provided
							</span>
						</div>
					</div>
				</div>

				<!-- Contact Information -->
				<div class="summary-section">
					<h3><i class="bi bi-person"></i> Contact Information</h3>
					<div class="summary-item">
						<div class="summary-label">Name</div>
						<div class="summary-value" v-if="modelValue.reporter?.name">
							{{ modelValue.reporter.name }}
						</div>
						<div class="summary-value status-badge status-not-provided" v-else>
							<i class="bi bi-x-circle"></i> Not provided
						</div>
					</div>
					<div class="summary-item">
						<div class="summary-label">Phone</div>
						<div class="summary-value" v-if="modelValue.reporter?.phone">
							{{ modelValue.reporter.phone }}
							<small class="text-muted" v-if="modelValue.reporter?.can_text"
								>(texting OK)</small
							>
						</div>
						<div class="summary-value status-badge status-not-provided" v-else>
							<i class="bi bi-x-circle"></i> Not provided
						</div>
					</div>
					<div class="summary-item">
						<div class="summary-label">Email</div>
						<div class="summary-value" v-if="modelValue.reporter?.email">
							{{ modelValue.reporter?.email }}
						</div>
						<div class="summary-value status-badge status-not-provided" v-else>
							<i class="bi bi-x-circle"></i> Not provided
						</div>
					</div>
				</div>
			</div>

			<!-- Submit Button -->
			<div class="d-flex gap-2 mt-4">
				<RouterLink class="btn btn-outline-secondary" to="process">
					Back
				</RouterLink>
				<button class="btn btn-primary flex-grow-1" @click="doContinue()">
					<i class="bi bi-check-circle"></i>
					Submit Response
				</button>
			</div>

			<div class="text-center mt-3">
				<small class="text-muted">
					By submitting, you confirm the information provided is accurate to the
					best of your knowledge.
				</small>
			</div>
		</main>
	</div>
</template>
<script setup lang="ts">
import { router } from "@/rmo/router";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import {
	type District,
	PermissionAccess,
	PublicReportCompliance,
} from "@/type/api";

interface Emits {
	(e: "doSubmit"): void;
}
interface Props {
	modelValue: PublicReportCompliance;
	district: District;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
function doContinue() {
	emit("doSubmit");
	router.push("./complete");
}
</script>
