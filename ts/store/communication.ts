import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import { SSEManager, SSEMessageResource } from "@/SSEManager";
import { useSessionStore } from "@/store/session";
import { Communication, CommunicationDTO } from "@/type/api";

export const useCommunicationStore = defineStore("communication", () => {
	// State
	const all = ref<Communication[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (
			msg.resource.startsWith("sync:communication") &&
			msg.type == "updated"
		) {
			fetchOne(msg.uri);
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
			const data = (await apiClient.JSONGet(url)) as CommunicationDTO[];

			all.value = data.map((c: CommunicationDTO) => Communication.fromJSON(c));
			return all.value;
		} catch (err) {
			console.error("Error loading communications:", err);
			throw err;
		} finally {
			loading.value = false;
		}
	}
	async function fetchOne(uri: string) {
		const data = (await apiClient.JSONGet(uri)) as CommunicationDTO;
		if (!all.value) {
			return;
		}
		for (var i = 0; i < all.value.length; i++) {
			const c = all.value[i];
			if (c.uri == data.uri) {
				all.value[i] = Communication.fromJSON(data);
			}
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
