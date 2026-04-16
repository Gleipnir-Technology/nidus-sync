<template>
	<h5 class="mb-4">Actions</h5>
	<template v-if="!selectedSite">
		<p>select a site to see actions</p>
	</template>
	<template v-if="selectedSite">
		<button
			class="btn btn-success action-btn"
			@click="emit('doRequestComplianceMailer', selectedSite?.id ?? 0)"
			:disabled="!selectedSite || submitting"
		>
			<span v-if="!submitting">
				<i class="bi bi-check-circle"></i> Send Compliance Mailer
			</span>
			<span v-else>
				<span class="spinner-border spinner-border-sm" role="status"></span>
				Submitting...
			</span>
		</button>
	</template>
</template>
<script setup lang="ts">
import { Site } from "@/type/api";

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
