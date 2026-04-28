import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import { SSEManager, SSEMessage } from "@/SSEManager";
import { useSessionStore } from "@/store/session";
import { Communication, CommunicationDTO } from "@/type/api";

export const useCommunicationStore = defineStore("communication", () => {
	// State
	const all = ref<Communication[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.resource.startsWith("rmo:")) {
			fetchAll();
		}
	});
	// Actions
	async function fetchAll(): Promise<Communication[]> {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();
			params.append("sort", "-created");
			//if (typeFilter.value) params.append("type", typeFilter.value);

			const url = `${session.urls.api.communication}?${params}`;
			const data = await apiClient.JSONGet(url);

			all.value = data.communications.map((c: CommunicationDTO) =>
				Communication.fromJSON(c),
			);
			return data.communications;
		} catch (err) {
			console.error("Error loading communications:", err);
			throw err;
		} finally {
			loading.value = false;
		}
	}
	return {
		// State
		all,
		loading,
		// Actions
		fetchAll,
	};
});
