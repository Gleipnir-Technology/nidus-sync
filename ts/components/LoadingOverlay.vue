<style scoped>
.loading-wrapper {
	min-height: 50px; /* Ensure minimum height for overlay positioning */
}

.loading-overlay {
	z-index: 1000;
	backdrop-filter: blur(1px);
}

/* Ensure child content is not interactive when loading */
.loading-wrapper:has(.loading-overlay) > *:not(.loading-overlay) {
	pointer-events: none;
	user-select: none;
}
</style>
<template>
	<div class="loading-wrapper position-relative">
		<slot></slot>

		<div
			v-if="isLoading"
			class="loading-overlay position-absolute top-0 start-0 w-100 h-100 d-flex justify-content-center align-items-center"
			:class="overlayClass"
		>
			<div class="text-center">
				<!-- Bootstrap spinner -->
				<div
					class="spinner-border text-primary"
					:class="spinnerSize"
					role="status"
				>
					<span class="visually-hidden">Loading...</span>
				</div>
				<!-- Optional loading text -->
				<div v-if="loadingText" class="mt-2">
					<small class="text-muted">{{ loadingText }}</small>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from "vue";

type SpinnerSize = "" | "spinner-border-sm";

interface Props {
	isLoading?: boolean;
	loadingText?: string;
	spinnerSize?: SpinnerSize;
	overlayOpacity?: string;
}

const props = withDefaults(defineProps<Props>(), {
	isLoading: false,
	loadingText: "",
	spinnerSize: "",
	overlayOpacity: "bg-light bg-opacity-75",
});

const overlayClass = computed((): string => {
	return props.overlayOpacity;
});
</script>
