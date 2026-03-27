import { defineStore } from "pinia";
import { ref } from "vue";
import { Upload } from "../types";
import { SSEManager } from "../SSEManager";
import { useUserStore } from "./user";

export const useUploadStore = defineStore("upload", () => {
	// State
	const _byID = ref<Map<int, Upload>>(new Map());
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
	function byID(id: int) {
		return _byID.value.get(id);
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
			for (const u of data.uploads) {
				_byID.value.set(u.id, u);
			}
		} catch (err) {
			console.error("Error loading uploads:", err);
			throw err;
		}
	}
	async function fetchOne(id: int) {
		const userStore = useUserStore();
		if (userStore.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const response = await fetch(`${userStore.urls.api.upload}/${id}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			_byID.value.set(data.id, data);
			return data;
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
		fetchOne,
	};
});
