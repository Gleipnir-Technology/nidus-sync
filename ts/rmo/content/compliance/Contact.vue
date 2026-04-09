<style scoped>
.benefit-box {
	background-color: #d1ecf1;
	border-left: 4px solid #0dcaf0;
	padding: 16px;
	border-radius: 4px;
}

.optional-badge {
	font-size: 0.85rem;
	color: #6c757d;
	font-weight: normal;
}
</style>
<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="6" />
		<!-- Main Content -->
		<main>
			<h2 class="h4 mb-3">Contact information</h2>

			<div class="benefit-box mb-4">
				<p class="mb-0">
					<i class="bi bi-info-circle"></i>
					<strong> Why share your contact information?</strong><br />
					<small>
						Providing your contact information helps the District review your
						response and coordinate with you if a visit is still needed. This
						can save time and prevent unnecessary follow-up actions.
					</small>
				</p>
			</div>

			<!-- Name -->
			<div class="mb-3">
				<label for="contact-name" class="form-label fw-semibold">
					Name
					<span class="optional-badge">(Optional)</span>
				</label>
				<input
					type="text"
					class="form-control"
					id="contact-name"
					name="name"
					placeholder="Enter your name"
					v-model="modelValue.contact.name"
				/>
			</div>

			<!-- Phone -->
			<div class="mb-3">
				<label for="contact-phone" class="form-label fw-semibold">
					Phone Number
					<span class="optional-badge">(Optional)</span>
				</label>
				<input
					type="tel"
					class="form-control"
					id="contact-phone"
					name="phone"
					placeholder="(555) 123-4567"
					v-model="modelValue.contact.phone"
				/>
			</div>

			<!-- Can we text? -->
			<div class="mb-3">
				<div class="form-check">
					<input
						class="form-check-input"
						type="checkbox"
						id="can-text"
						v-model="modelValue.contact.can_text"
					/>
					<label class="form-check-label" for="can-text">
						You may send text messages to this number
					</label>
				</div>
				<small class="text-muted ms-4 d-block mt-1">
					Text messages allow for faster communication and updates
				</small>
			</div>

			<!-- Email -->
			<div class="mb-4">
				<label for="contact-email" class="form-label fw-semibold">
					Email Address
					<span class="optional-badge">(Optional)</span>
				</label>
				<input
					type="email"
					class="form-control"
					id="contact-email"
					placeholder="your.email@example.com"
					v-model="modelValue.contact.email"
				/>
				<div class="form-text">
					We'll send you a confirmation and any updates about this request
				</div>
			</div>

			<div class="alert alert-light border" role="alert">
				<small class="text-muted">
					<i class="bi bi-shield-check"></i>
					Your contact information will only be used for this compliance matter
					and will be kept confidential.
				</small>
			</div>

			<!-- Navigation Buttons -->
			<div class="d-flex gap-2 mt-4">
				<RouterLink class="btn btn-outline-secondary" to="./permission">
					Back
				</RouterLink>
				<button class="btn btn-primary flex-grow-1" @click="doContinue()">
					Continue
				</button>
			</div>
		</main>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";

import { router } from "@/rmo/router";
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import type { Compliance } from "@/rmo/view/Compliance.vue";

export interface Contact {
	name: string;
	phone: string;
	can_text: boolean;
	email: string;
}
interface Emits {
	(e: "doContact"): void;
	(e: "update:modelValue", value: Compliance): void;
}
interface Props {
	district: District;
	modelValue: Compliance;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
function doContinue() {
	emit("update:modelValue", props.modelValue);
	emit("doContact");
	router.push("./process");
}
</script>
