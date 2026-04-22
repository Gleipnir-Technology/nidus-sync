<template>
	<router-view />
</template>

<script setup lang="ts">
import * as Sentry from "@sentry/vue";
import { onMounted } from "vue";
import { apiClient } from "@/client";
import router from "@/route/config";

import { SSEManager, type SSEMessage } from "@/SSEManager";

async function sentryInit() {
	const config = await apiClient.JSONGet("/api");
	Sentry.init({
		dsn: config.DSN,
		integrations: [Sentry.browserTracingIntegration({ router })],
		environment: config.ENVIRONMENT,
		release: config.RELEASE,
		tracesSampleRate: 0.01,
	});
	console.log("sentry initialized");
}
onMounted(() => {
	SSEManager.connect("/api/events");
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.type != "heartbeat") {
			console.log("SSE", msg);
		}
	});
	sentryInit();
});
</script>
