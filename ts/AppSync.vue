<template>
	<router-view />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { apiClient } from "@/client";
import router from "@/route/config";

import { SSEManager, type SSEMessageResource } from "@/SSEManager";

onMounted(() => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
});
</script>
