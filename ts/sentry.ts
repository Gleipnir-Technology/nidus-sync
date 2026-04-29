import { type App } from "vue";
import { type Pinia } from "pinia";
import { type Router } from "vue-router";
import * as Sentry from "@sentry/vue";
import { apiClient } from "@/client";
import { APIProperties } from "@/type/api";

export async function Init(app: App, pinia: Pinia) {
	const api_info: APIProperties = await apiClient.JSONGet("/api");
	Sentry.init({
		app,
		dsn: api_info.sentry_dsn,
		//integrations: [Sentry.browserTracingIntegration({ router })],
		environment: api_info.environment,
		release:
			api_info.version.revision +
			(api_info.version.is_modified ? "-dirty" : ""),
		tracesSampleRate: 0.01,
	});
	pinia.use(Sentry.createSentryPiniaPlugin());
	console.log("sentry initialized", api_info.sentry_dsn, api_info.environment);
}
