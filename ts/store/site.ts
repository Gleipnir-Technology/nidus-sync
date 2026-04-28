import { defineStore } from "pinia";
import { ref } from "vue";
import { SSEManager, SSEMessageResource } from "@/SSEManager";
import { Site, SiteListResponse } from "@/type/api";
import { useSessionStore } from "@/store/session";

export const useStoreSite = defineStore("site", () => {
	// State
	const _byID = ref<Map<number, Site>>(new Map());
	const loading = ref<boolean>(false);
	const error = ref<string | null>(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.resource.startsWith("sync:site")) {
			fetchAll();
		}
	});
	// Actions
	function all(): Site[] {
		return Array.from(_byID.value.values());
	}
	function byID(id: number): Site | undefined {
		return _byID.value.get(id);
	}
	async function fetchAll(): Promise<Site[]> {
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

			const response = await fetch(`${session.urls.api.site}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data: Site[] = await response.json();
			_byID.value = new Map();
			for (const t of data) {
				_byID.value.set(t.id, t);
			}
			return data;
		} catch (err) {
			error.value = err instanceof Error ? err.message : "Unknown error";
			console.error("Error loading tasks:", err);
			throw err;
		} finally {
			loading.value = false;
		}
	}
	async function fetchOne(id: number) {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const response = await fetch(`${session.urls.api.site}/${id}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			_byID.value.set(data.id, data);
			return data;
		} catch (err) {
			console.error("Error loading tasks:", err);
			throw err;
		}
	}
	function remove(id: number) {
		_byID.value.delete(id);
	}

	return {
		// State
		all,
		// Actions
		byID,
		fetchAll,
		fetchOne,
		remove,
	};
});
