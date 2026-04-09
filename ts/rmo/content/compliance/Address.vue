<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<!-- Progress Bar -->
		<ProgressBarCompliance :step="2" />
		<main>
			<h2 class="h4 mb-3">Confirm the property address</h2>

			<p class="text-muted mb-4">
				Please enter the address so we can match your response with our records.
			</p>

			<AddressAndMapLocator v-model="locator" />

			<div class="d-flex gap-2 mt-4">
				<RouterLink class="btn btn-outline-secondary" to="../compliance">
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

import type { District } from "@/type/api";
import { router } from "@/rmo/router";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import AddressAndMapLocator from "@/rmo/components/AddressAndMapLocator.vue";
import { Locator } from "@/type/map";
interface Emits {
	(e: "doLocator", locator: Locator | null): void;
}
interface Props {
	district: District;
}
const emit = defineEmits<Emits>();
const error = ref<string>("");
const props = defineProps<Props>();
const locator = ref<Locator | null>(null);
function doContinue() {
	emit("doLocator", locator.value);
	router.push("./concern");
}
</script>
