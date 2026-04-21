<template>
	<router-view />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { SSEManager, type SSEMessage } from "@/SSEManager";

onMounted(() => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
});
</script>
