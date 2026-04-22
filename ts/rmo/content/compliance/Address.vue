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

			<AddressAndMapLocator
				:initialCamera="initialCamera"
				v-model="modelValue.address"
			/>

			<div class="d-flex gap-2 mt-4">
				<RouterLink
					class="btn btn-outline-secondary"
					:to="routes.ComplianceIntro(props.publicID)"
				>
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
import { computed, ref } from "vue";

import { router } from "@/rmo/route/config";
import type { District, PublicReportCompliance } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import AddressAndMapLocator from "@/rmo/components/AddressAndMapLocator.vue";
import { Camera } from "@/type/map";
import { useRoutes } from "@/rmo/route/use";

interface Emits {
	(e: "doAddress"): void;
	(e: "update:modelValue", value: PublicReportCompliance): void;
}
interface Props {
	district: District;
	modelValue: PublicReportCompliance;
	publicID: string;
}
const emit = defineEmits<Emits>();
const error = ref<string>("");
const props = defineProps<Props>();
const initialCamera = computed((): Camera | undefined => {
	if (props.modelValue.location) {
		return {
			location: props.modelValue.location,
			zoom: 15,
		};
	}
	return undefined;
});
const routes = useRoutes();
function doContinue() {
	emit("update:modelValue", props.modelValue);
	emit("doAddress");
	if (props.modelValue.concerns.length > 0) {
		router.push(routes.ComplianceConcern(props.publicID));
	} else {
		router.push(routes.ComplianceEvidence(props.publicID));
	}
}
</script>
