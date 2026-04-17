<template>
	<h5 class="mb-4">Actions</h5>
	<template v-if="!selectedSite">
		<p>select a site to see actions</p>
	</template>
	<template v-if="selectedSite">
		<ButtonLoading
			@click="emit('doRequestComplianceMailer', selectedSite?.id ?? 0)"
			:disabled="!selectedSite"
			icon="bi-check-circle"
			:loading="submitting"
			text="Send Compliance Mailer"
			variant="success"
		/>
	</template>
</template>
<script setup lang="ts">
import { Site } from "@/type/api";
import ButtonLoading from "@/components/common/ButtonLoading.vue";

interface Emits {
	(e: "doRequestComplianceMailer", id: number): void;
}
interface Props {
	selectedSite: Site | undefined;
	submitting?: boolean;
}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
	submitting: false,
});
</script>
