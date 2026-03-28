import { defineStore } from "pinia";
import { ref } from "vue";
import { ReviewTask } from "../types";
import { SSEManager } from "../SSEManager";
import { useUserStore } from "./user";

export const useReviewTaskStore = defineStore("review-task", () => {
	// State
	const _byID = ref<Map<int, ReviewTask>>(new Map());
	const all = ref<ReviewTask[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.resource.startsWith("review-task")) {
			fetchAll();
		}
	});
	// Actions
	function byID(id: int) {
		return _byID.value.get(id);
	}
	async function fetchAll(): Promise<void> {
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

			const response = await fetch(
				`${userStore.urls.api.review_task}?${params}`,
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			all.value = data.tasks;
			for (const t of data.tasks) {
				_byID.value.set(t.id, t);
			}
		} catch (err) {
			error.value = err instanceof Error ? err.message : "Unknown error";
			console.error("Error loading tasks:", err);
			throw err;
		} finally {
			loading.value = false;
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
			const response = await fetch(`${userStore.urls.api.review_task}/${id}`);

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

	return {
		// State
		all,
		// Actions
		byID,
		fetchAll,
		fetchOne,
	};
});
