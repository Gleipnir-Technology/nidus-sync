
import { defineStore } from "pinia";
import { ref } from "vue";
import { Upload } from "../types";
import { SSEManager } from "../SSEManager";
import { useUserStore } from "./user";

export const useUploadStore = defineStore("upload", () => {
	// State
	const all = ref<Upload[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.resource.startsWith("upload")) {
			fetchAll();
		}
	});
	// Actions
	function byID(id: int): Upload? {
		if (all.value == null) {
			return null
		}
		return all.value.find((upload) => upload.id == id);
	}
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

			const response = await fetch(`${userStore.urls.api.upload}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			all.value = data.uploads;
		} catch (err) {
			console.error("Error loading uploads:", err);
			throw err;
		}
	}

	return {
		// State
		all,
		// Actions
		byID,
		fetchAll,
	};
});
