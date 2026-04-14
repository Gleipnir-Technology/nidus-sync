import { defineStore } from "pinia";
import { ref } from "vue";
import { ServiceRequest } from "@/type/api";
import { SSEManager, SSEMessage } from "@/SSEManager";
import { useSessionStore } from "@/store/session";

export const useStoreServiceRequest = defineStore("service-request", () => {
	// State
	const all = ref<ServiceRequest[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.resource.startsWith("sync:service-request")) {
			fetchAll();
		}
	});
	// Actions
	async function fetchAll(): Promise<ServiceRequest[]> {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();

			const response = await fetch(
				`${session.urls.api.service_request}?${params}`,
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = (await response.json()) as ServiceRequest[];
			all.value = data;
			return data;
		} catch (err) {
			console.error("Error loading communications:", err);
			throw err;
		}
	}
	return {
		// State
		all,
		// Actions
		fetchAll,
	};
});
