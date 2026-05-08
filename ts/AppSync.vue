<style>
.global-error-toast {
	position: fixed;
	top: 20px;
	right: 20px;
	background: #c00;
	color: white;
	padding: 16px 20px;
	border-radius: 8px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	z-index: 9999;
	max-width: 400px;
	animation: slideIn 0.3s ease;
}

@keyframes slideIn {
	from {
		transform: translateX(400px);
		opacity: 0;
	}
	to {
		transform: translateX(0);
		opacity: 1;
	}
}
</style>

<template>
	<div id="app">
		<div v-if="error" class="global-error-toast">
			{{ error.message }}
			<button @click="errorClear">x</button>
		</div>
		<router-view />
	</div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { apiClient } from "@/client";
import router from "@/route/config";
import { useErrorHandler } from "@/composable/error-handler";

import { SSEManager, type SSEMessageResource } from "@/SSEManager";

const { error, errorClear } = useErrorHandler();

onMounted(() => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
});
</script>
