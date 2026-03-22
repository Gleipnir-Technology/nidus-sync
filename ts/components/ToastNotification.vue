<template>
	<div class="toast-container position-fixed bottom-0 end-0 p-3">
		<div class="toast" :class="{ show: show }" role="alert">
			<div class="toast-header">
				<i class="bi bi-check-circle text-success me-2"></i>
				<strong class="me-auto">{{ title }}</strong>
				<button type="button" class="btn-close" @click="show = false"></button>
			</div>
			<div class="toast-body">{{ message }}</div>
		</div>
	</div>
</template>

<script setup lang="ts">
interface Props {
	message: string;
	show: boolean;
	title: string;
}
const props = defineProps<Props>();
</script>
