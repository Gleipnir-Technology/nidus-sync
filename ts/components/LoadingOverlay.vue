<template>
	<div class="loading-wrapper position-relative">
		<!-- Child content slot -->
		<slot></slot>

		<!-- Loading overlay -->
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

<script>
export default {
	name: "LoadingOverlay",
	props: {
		isLoading: {
			type: Boolean,
			default: false,
		},
		loadingText: {
			type: String,
			default: "",
		},
		spinnerSize: {
			type: String,
			default: "",
			validator: (value) => ["", "spinner-border-sm"].includes(value),
		},
		overlayOpacity: {
			type: String,
			default: "bg-light bg-opacity-75",
		},
	},
	computed: {
		overlayClass() {
			return this.overlayOpacity;
		},
	},
};
</script>

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
