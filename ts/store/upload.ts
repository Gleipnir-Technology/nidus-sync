import { defineStore } from "pinia";
import { ref } from "vue";
import { Upload } from "@/type/api";
import { SSEManager, type SSEMessageResource } from "@/SSEManager";
import { useSessionStore } from "@/store/session";

export const useUploadStore = defineStore("upload", () => {
	// State
	const _byID = ref<Map<number, Upload>>(new Map());
	const all = ref<Upload[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.resource.startsWith("sync:upload")) {
			fetchAll();
		}
	});
	// Actions
	function byID(id: number) {
		return _byID.value.get(id);
	}
	async function fetchAll() {
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

			const response = await fetch(`${session.urls.api.upload}?${params}`);

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
	async function fetchOne(id: number) {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const response = await fetch(`${session.urls.api.upload}/${id}`);

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
