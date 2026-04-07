<template>
	<div class="mb-4">
		<div class="d-flex justify-content-between align-items-center mb-2">
			<span class="small text-muted">Step {{ step }} of {{ maxSteps }}</span>
		</div>
		<div class="progress" style="height: 8px">
			<div
				class="progress-bar"
				role="progressbar"
				:style="{ width: (step / maxSteps) * 100 + '%' }"
				aria-valuenow="12"
				aria-valuemin="0"
				aria-valuemax="100"
			></div>
		</div>
	</div>
</template>
<script setup lang="ts">
interface Props {
	step: number;
}
const props = defineProps<Props>();
const maxSteps = 8;
</script>
