<template>
	<h5 class="mb-4">Actions</h5>
	<template v-if="!selectedSite">
		<p>select a site to see actions</p>
	</template>
	<template v-if="selectedSite">
		<template v-if="session.organization?.lob_address_id">
			<ButtonLoading
				@click="emit('doRequestComplianceMailer', selectedSite?.id ?? 0)"
				:disabled="!selectedSite"
				icon="bi-check-circle"
				:loading="submitting"
				text="Send Compliance Mailer"
				variant="success"
			/>
		</template>
		<template v-else>
			<p>Set Lob Address ID</p>
		</template>
	</template>
</template>
<script setup lang="ts">
import ButtonLoading from "@/components/common/ButtonLoading.vue";
import { useSessionStore } from "@/store/session";
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
const session = useSessionStore();
</script>
