import { defineStore } from "pinia";
import { ref } from "vue";
import { Sync } from "@/type/api";
import { SSEManager, SSEMessage } from "@/SSEManager";
import { useSessionStore } from "@/store/session";

export const useStoreSync = defineStore("sync", () => {
	// State
	const all = ref<Sync[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.resource.startsWith("sync:sync")) {
			fetchAll();
		}
	});
	// Actions
	async function fetchAll(): Promise<Sync[]> {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();

			const response = await fetch(`${session.urls.api.sync}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = (await response.json()) as Sync[];
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
