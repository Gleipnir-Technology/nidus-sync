<style scoped>
.completion-container {
	max-width: 500px;
	margin: 0 auto;
	padding: 32px 16px;
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.district-branding {
	text-align: center;
	margin-bottom: 32px;
}

.district-branding img {
	height: 80px;
	width: auto;
	margin-bottom: 16px;
}

.district-branding h1 {
	font-size: 1.25rem;
	color: #495057;
	margin-bottom: 8px;
}

.district-branding .phone {
	color: #6c757d;
	font-size: 0.9rem;
}

.success-icon {
	width: 80px;
	height: 80px;
	background-color: #d1e7dd;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 0 auto 24px;
}

.success-icon i {
	font-size: 2.5rem;
	color: #0f5132;
}

.warning-icon {
	width: 80px;
	height: 80px;
	background-color: #fff3cd;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 0 auto 24px;
}

.warning-icon i {
	font-size: 2.5rem;
	color: #997404;
}

.alert-icon {
	width: 80px;
	height: 80px;
	background-color: #f8d7da;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 0 auto 24px;
}

.alert-icon i {
	font-size: 2.5rem;
	color: #842029;
}

.message-card {
	background-color: #fff;
	border-radius: 12px;
	padding: 24px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	width: 100%;
	margin-bottom: 24px;
}

.message-card h2 {
	font-size: 1.5rem;
	margin-bottom: 16px;
	text-align: center;
}

.message-card p {
	color: #495057;
	line-height: 1.6;
	margin-bottom: 12px;
}

.message-card p:last-child {
	margin-bottom: 0;
}

.info-box {
	background-color: #f8f9fa;
	border-left: 4px solid #6c757d;
	padding: 16px;
	border-radius: 4px;
	margin-top: 16px;
}

.warning-box {
	background-color: #fff3cd;
	border-left: 4px solid #ffc107;
	padding: 16px;
	border-radius: 4px;
	margin-top: 16px;
}

.reference-number {
	text-align: center;
	color: #6c757d;
	font-size: 0.9rem;
	margin-top: 24px;
}

.action-buttons {
	width: 100%;
}
</style>
<template>
	<div class="completion-container">
		<!-- District Branding -->
		<div class="district-branding">
			<img :src="district.url_logo" :alt="district.name + 'logo'" />
			<h1>{{ district.name }}</h1>
			<div class="phone" v-if="district.phone_office">
				<i class="bi bi-telephone"></i> {{ district.phone_office }}
			</div>
		</div>

		<template v-if="hasCompleteResponse">
			<!-- Mode 1: Complete response with contact info -->
			<div class="success-icon">
				<i class="bi bi-check-circle-fill"></i>
			</div>

			<div class="message-card">
				<h2>Thank you</h2>
				<p>
					Your response has been submitted successfully. We will review your
					submission and contact you if further action is needed.
				</p>
				<div class="info-box">
					<p class="mb-2">
						<strong>What you can expect:</strong>
					</p>
					<ul class="mb-0 small">
						<li>Our team will review your photos and information</li>
						<li>If we need to schedule a visit, we'll contact you first</li>
						<li>
							You'll receive updates at the contact information you provided
						</li>
					</ul>
				</div>
			</div>
		</template>
		<template v-else-if="hasUsefulInfo">
			<!-- Mode 2: Useful info but no contact -->
			<div class="warning-icon">
				<i class="bi bi-exclamation-circle-fill"></i>
			</div>

			<div class="message-card">
				<h2>Response Received</h2>
				<p>
					Thank you for your submission. We will review your information and let
					you know if further action is needed.
				</p>
				<div class="warning-box">
					<p class="mb-2">
						<strong><i class="bi bi-info-circle"></i> Important notice:</strong>
					</p>
					<p class="mb-0 small">
						You did not provide contact information. If further action is
						needed, the District may need to use warrant authority to enter the
						property. We prefer to coordinate access directly, and contact
						information makes that much easier.
					</p>
				</div>
				<p class="text-center mt-3 mb-0">
					<small class="text-muted">
						If you'd like to add contact information, please call our office.
					</small>
				</p>
			</div>
		</template>
		<template v-else>
			<!-- Mode 3: Insufficient information -->
			<div class="alert-icon">
				<i class="bi bi-exclamation-triangle-fill"></i>
			</div>

			<div class="message-card">
				<h2>Response Received</h2>
				<p>
					Your response has been recorded, but it does not contain enough
					information for us to resolve this matter.
				</p>
				<div class="warning-box">
					<p class="mb-2">
						<strong
							><i class="bi bi-exclamation-triangle"></i> Important:</strong
						>
					</p>
					<p class="mb-2 small">
						This response is not likely to resolve the issue and may require
						warrant entry on the property. If you want to help avoid that,
						please provide contact information or other evidence.
					</p>
					<p class="mb-0 small">
						<strong>You can still:</strong>
					</p>
					<ul class="mb-0 small">
						<li>Call our office to provide additional information</li>
						<li>Email photos showing current conditions</li>
						<li>Schedule a time for inspection</li>
					</ul>
				</div>
			</div>
		</template>

		<!-- Reference Number -->
		<div class="reference-number">
			<small>
				Reference number: <strong>{{ referenceNumber }}</strong>
			</small>
		</div>

		<!-- Action Buttons -->
		<div class="action-buttons mt-4">
			<div class="d-grid gap-2">
				<template v-if="hasCompleteResponse">
					<button class="btn btn-outline-primary" @click="changeResponse()">
						<i class="bi bi-dash"></i>Remove useful and complete
					</button>
				</template>
				<template v-else-if="hasUsefulInfo">
					<button class="btn btn-outline-primary" @click="changeResponse()">
						<i class="bi bi-plus"></i>Add complete response
					</button>
				</template>
				<template v-else>
					<button class="btn btn-outline-primary" @click="changeResponse()">
						<i class="bi bi-plus"></i>Add useful response
					</button>
				</template>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
interface Props {
	district: District;
}
const props = defineProps<Props>();
const referenceNumber = "abc-123";
const hasCompleteResponse = ref<boolean>(false);
const hasUsefulInfo = ref<boolean>(false);
function changeResponse() {
	if (hasCompleteResponse.value && hasUsefulInfo.value) {
		hasCompleteResponse.value = false;
		hasUsefulInfo.value = false;
	} else if (hasUsefulInfo.value) {
		hasCompleteResponse.value = true;
	} else {
		hasUsefulInfo.value = true;
	}
}
</script>
