import { defineStore } from "pinia";
import { ref } from "vue";
import { SSEManager } from "@/SSEManager";
import {
	Organization,
	URLs,
	User,
	UserNotificationCounts,
	UserResponse,
} from "@/types";

export const useSessionStore = defineStore("session", () => {
	// State
	const error = ref<string | null>(null);
	const loading = ref(false);
	const user = ref<User | null>(null);
	const urls = ref<URLs | null>(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.type !== "heartbeat") {
			fetchSession();
		}
	});

	// Actions
	async function fetchSession(): Promise<UserResponse> {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/user/self");
			if (!response.ok) throw new Error("Failed to fetch user");

			const data: UserResponse = await response.json();
			user.value = data.self;
			urls.value = data.urls;
			return data;
		} catch (e) {
			error.value = e instanceof Error ? e.message : "an error ocurred";
			console.error("Error fetching user:", e);
			throw new Error(error.value);
		} finally {
			loading.value = false;
		}
	}

	async function isAuthenticated(): Promise<boolean> {
		console.log("pretend check user auth");
		return true;
	}
	return {
		// State
		error,
		loading,
		user,
		urls,
		// Actions
		fetchSession,
		isAuthenticated,
	};
});
