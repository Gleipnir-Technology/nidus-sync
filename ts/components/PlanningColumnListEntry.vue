<template>
	<div class="small fw-semibold">{{ signal.title }}</div>
	<div class="signal-address">
		{{ shortAddress(signal.address) }}
	</div>
	<div class="text-muted small">{{ signal.description }}</div>
	<span v-if="signal.badge" class="badge bg-secondary mt-1">
		{{ signal.badge }}
	</span>
</template>
<script lang="ts">
interface Props {
	signal: Signal;
}
const props = defineProps<Props>();
</script>
