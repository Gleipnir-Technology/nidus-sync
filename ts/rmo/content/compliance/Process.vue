<style scoped>
.process-step {
	background-color: #fff;
	border-radius: 8px;
	padding: 16px;
	margin-bottom: 16px;
	border-left: 4px solid #0d6efd;
	display: flex;
	gap: 16px;
}

.step-number {
	background-color: #0d6efd;
	color: white;
	width: 36px;
	height: 36px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: bold;
	flex-shrink: 0;
}

.step-content {
	flex-grow: 1;
}

.encouragement-box {
	background-color: #d1ecf1;
	border-left: 4px solid #0dcaf0;
	padding: 16px;
	border-radius: 4px;
}

.reminder-box {
	background-color: #fff3cd;
	border-left: 4px solid #ffc107;
	padding: 16px;
	border-radius: 4px;
}
</style>

<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="7" />
		<main>
			<h2 class="h4 mb-3">What happens next</h2>

			<p class="text-muted mb-4">
				Understanding the review process helps you know what to expect.
			</p>

			<!-- Process Steps -->
			<div class="mb-4">
				<!-- Step 1 -->
				<div class="process-step">
					<div class="step-number">1</div>
					<div class="step-content">
						<h3 class="h6 mb-2">We review your response</h3>
						<p class="mb-0 small text-muted">
							Our team will evaluate the photos, comments, and information
							you've provided to assess the current conditions at the property.
						</p>
					</div>
				</div>

				<!-- Step 2 -->
				<div class="process-step">
					<div class="step-number">2</div>
					<div class="step-content">
						<h3 class="h6 mb-2">Your response may reduce or close follow-up</h3>
						<p class="mb-0 small text-muted">
							If your photos and information clearly show that the concern has
							been addressed, we may be able to reduce or close this matter
							without additional inspection.
						</p>
					</div>
				</div>

				<!-- Step 3 -->
				<div class="process-step">
					<div class="step-number">3</div>
					<div class="step-content">
						<h3 class="h6 mb-2">We'll contact you if inspection is needed</h3>
						<p class="mb-0 small text-muted">
							If an inspection is still necessary and you provided contact
							information, we will try to reach you in advance to coordinate
							timing and access.
						</p>
					</div>
				</div>

				<!-- Step 4 -->
				<div class="process-step">
					<div class="step-number">4</div>
					<div class="step-content">
						<h3 class="h6 mb-2">
							The compliance process continues if unresolved
						</h3>
						<p class="mb-0 small text-muted">
							If the mosquito breeding source remains unaddressed, the District
							will continue with standard compliance procedures to protect
							public health.
						</p>
					</div>
				</div>
			</div>

			<!-- Encouragement Box -->
			<div class="encouragement-box mb-3">
				<p class="mb-2">
					<strong
						><i class="bi bi-check-circle"></i> Your response helps</strong
					>
				</p>
				<p class="mb-0 small">
					By providing photos, access information, and contact details, you give
					our team the ability to review the situation thoroughly before taking
					further action. This can save everyone time and help resolve the
					matter more efficiently.
				</p>
			</div>

			<!-- Reminder Box -->
			<div class="reminder-box">
				<p class="mb-0 small">
					<strong><i class="bi bi-exclamation-circle"></i> Important:</strong>
					Submitting this form does not automatically close this compliance
					request. The District must verify that mosquito breeding conditions
					have been corrected to protect community health.
				</p>
			</div>

			<!-- Navigation Buttons -->
			<div class="d-flex gap-2 mt-4">
				<RouterLink
					class="btn btn-outline-secondary"
					:to="routes.ComplianceContact(publicID)"
				>
					Back
				</RouterLink>
				<RouterLink
					class="btn btn-primary flex-grow-1"
					:to="routes.ComplianceSubmit(publicID)"
					><i class="bi bi-caret-right-fill"></i>
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
import { useRoutes } from "@/rmo/route/use";
interface Props {
	district: District;
	publicID: string;
}
const props = defineProps<Props>();
const routes = useRoutes();
</script>
