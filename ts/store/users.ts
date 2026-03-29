import { defineStore } from "pinia";
import { ref } from "vue";
import { User } from "../types";
import { SSEManager } from "../SSEManager";
import { useUserStore } from "./user";

export const useUsersStore = defineStore("users", () => {
	// State
	const _byID = ref<Map<int, User>>(new Map());
	const all = ref<User[] | null>(null);
	const loading = ref(false);
	const error = ref(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.resource.startsWith("users")) {
			fetchAll();
		}
	});
	// Actions
	function byID(id: int) {
		const result = _byID.value.get(id);
		console.log("user", id, result);
		return result;
	}
	async function fetchAll(): Promise<User[]> {
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

			const response = await fetch(`${userStore.urls.api.users}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			all.value = data.users;
			for (const u of data.users) {
				_byID.value.set(u.id, u);
			}
			return data.users;
		} catch (err) {
			console.error("Error loading users:", err);
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
			const response = await fetch(`${userStore.urls.api.user}/${id}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			_byID.value.set(data.id, data);
			return data;
		} catch (err) {
			console.error("Error loading users:", err);
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
