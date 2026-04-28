<style scoped>
.slide-fade-enter-active {
	transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
	transition: all 0.3s ease-in;
}

.slide-fade-enter-from {
	transform: translateX(-50%) translateY(20px);
	opacity: 0;
}

.slide-fade-leave-to {
	transform: translateX(-50%) translateY(20px);
	opacity: 0;
}
</style>

<template>
	<Transition name="slide-fade">
		<div
			v-if="updateDate"
			class="alert alert-info alert-dismissible fade show position-fixed bottom-0 start-50 translate-middle-x mb-3 shadow-sm"
			style="z-index: 1050; max-width: 400px"
			role="alert"
		>
			<div class="d-flex align-items-center gap-2">
				<div class="flex-grow-1">
					<div class="fw-semibold">App updates available</div>
					<small class="text-muted">{{ timeSinceUpdate }}</small>
				</div>

				<i
					ref="infoIcon"
					class="bi bi-info-circle text-primary"
					style="cursor: pointer; font-size: 1.2rem"
					data-bs-toggle="popover"
					data-bs-trigger="hover focus"
					data-bs-placement="top"
					data-bs-content="Please save your work before refreshing the page."
					tabindex="0"
				></i>

				<button class="btn btn-sm btn-primary" @click="refreshPage">
					Refresh
				</button>
			</div>
		</div>
	</Transition>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { Popover } from "bootstrap";

interface Props {
	updateDate: Date | null;
}

const props = defineProps<Props>();

const infoIcon = ref<HTMLElement | null>(null);
let popoverInstance: Popover | null = null;
let intervalId: number | null = null;

// Force reactivity update every minute
const tick = ref(0);

const timeSinceUpdate = computed(() => {
	if (!props.updateDate) return "";

	// Access tick to ensure reactivity
	tick.value;

	const now = new Date();
	const diff = now.getTime() - props.updateDate.getTime();
	const minutes = Math.floor(diff / 60000);
	const hours = Math.floor(minutes / 60);
	const days = Math.floor(hours / 24);

	if (days > 0) return `${days} day${days > 1 ? "s" : ""} ago`;
	if (hours > 0) return `${hours} hour${hours > 1 ? "s" : ""} ago`;
	if (minutes > 0) return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
	return "just now";
});

const refreshPage = () => {
	window.location.reload();
};

// Initialize popover when icon is mounted
watch(infoIcon, (element) => {
	if (element && !popoverInstance) {
		popoverInstance = new Popover(element);
	}
});

onMounted(() => {
	// Update time display every minute
	intervalId = window.setInterval(() => {
		tick.value++;
	}, 60000);
});

onUnmounted(() => {
	if (popoverInstance) {
		popoverInstance.dispose();
	}
	if (intervalId !== null) {
		clearInterval(intervalId);
	}
});
</script>
