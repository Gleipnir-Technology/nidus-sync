import { defineStore } from "pinia";
import { ref } from "vue";
import { SSEManager, SSEMessageResource } from "@/SSEManager";
import { ReviewTask, ReviewTaskListResponse } from "@/type/api";
import { useSessionStore } from "@/store/session";

export const useStoreReviewTask = defineStore("review-task", () => {
	// State
	const _byID = ref<Map<number, ReviewTask>>(new Map());
	const loading = ref<boolean>(false);
	const error = ref<string | null>(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.resource.startsWith("sync:review-task")) {
			fetchAll();
		}
	});
	// Actions
	function all(): ReviewTask[] {
		return Array.from(_byID.value.values());
	}
	function byID(id: number): ReviewTask | undefined {
		return _byID.value.get(id);
	}
	async function fetchAll(): Promise<ReviewTask[]> {
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

			const response = await fetch(`${session.urls.api.review_task}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data: ReviewTaskListResponse = await response.json();
			_byID.value = new Map();
			for (const t of data.tasks) {
				_byID.value.set(t.id, t);
			}
			return data.tasks;
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
			const response = await fetch(`${session.urls.api.review_task}/${id}`);

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
