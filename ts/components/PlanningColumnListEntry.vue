<template>
	<div class="small fw-semibold">{{ title(signal) }}</div>
	<div class="signal-address">
		{{ shortAddress(signal.address) }}
	</div>
	<div class="text-muted small">{{ signal.description }}</div>
	<span v-if="signal.badge" class="badge bg-secondary mt-1">
		{{ signal.badge }}
	</span>
</template>

<script setup lang="ts">
import { shortAddress } from "../format";
interface Props {
	signal: Signal;
};
const props = defineProps<Props>();
function title(signal: Signal): string {
	if (signal.type == "flyover pool") {
		return "Green pool";
	} else {
		return "Unknown signal";
	}
}
</script>
