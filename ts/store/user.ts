import { defineStore } from "pinia";
import { ref } from "vue";
import { Session, User } from "@/types";
import { SSEManager, type SSEMessage } from "@/SSEManager";
import { useSessionStore } from "@/store/session";

export const useUserStore = defineStore("users", () => {
	// State
	const _all = ref<User[] | null>(null);
	const _byID = ref<Map<number, User>>(new Map());
	const error = ref(null);
	const loading = ref(false);
	const ongoingFetch = ref<Promise<User[]> | null>(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.resource.startsWith("sync:user")) {
			fetchAll();
		}
	});
	// Actions
	function byID(id: number): User | null {
		const result = _byID.value.get(id);
		if (!result) {
			return null;
		}
		console.log("user", id, result);
		return result;
	}
	async function byURI(uri: string): Promise<User | null> {
		const all = await withAll();
		const result = all.find((u: User) => u.uri == uri);
		return result || null;
	}
	async function fetchAll(): Promise<User[]> {
		const sessionStore = useSessionStore();
		const session = await sessionStore.get();
		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();
			params.append("sort", "-created");
			//if (typeFilter.value) params.append("type", typeFilter.value);

			const response = await fetch(`${session.urls.api.user}?${params}`);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const users = await response.json();
			_all.value = users;
			for (const u of users) {
				_byID.value.set(u.id, u);
			}
			return users;
		} catch (err) {
			console.error("Error loading users:", err);
			throw err;
		}
	}
	async function withAll(): Promise<User[]> {
		if (_all.value != null) {
			return _all.value;
		}

		if (ongoingFetch.value !== null) {
			return ongoingFetch.value;
		}

		ongoingFetch.value = fetchAll().finally(() => {
			ongoingFetch.value = null;
		});
		return ongoingFetch.value;
	}
	async function fetchOne(id: number) {
		const session = useSessionStore();
		if (session.urls == null) {
			throw new Error("can't fetch without user URL data");
		}

		loading.value = true;
		error.value = null;
		try {
			const response = await fetch(`${session.urls.api.user}/${id}`);

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
		// Actions
		byID,
		byURI,
		fetchAll,
		fetchOne,
		withAll,
	};
});
