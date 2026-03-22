import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { Signal } from "../types";
import { SSEManager } from "../SSEManager";
import { useUserStore } from "./user";

export const useSignalStore = defineStore("signal", () => {
	// State
	const all = ref<Signal[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.resource.startsWith("signal")) {
			fetchAll();
		}
	});
	// Actions
	async function fetchAll() {
		const userStore = useUserStore();
		if (userStore.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();
			params.append("sort", "-created");
			//if (typeFilter.value) params.append("type", typeFilter.value);

			const response = await fetch(`${userStore.urls.api.signal}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			all.value = data.signals;
		} catch (err) {
			console.error("Error loading signals:", err);
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
