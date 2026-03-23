import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { SSEManager } from "../SSEManager";

// Define interfaces matching your Go structs
interface URLsAPI {
	communication: string;
	signal: string;
}

interface URLs {
	api: URLsAPI;
	tegola: string;
	tile: string;
}

interface User {
	display_name: string;
	initials: string;
	notification_counts: NotificationCounts;
	notifications: any[]; // Replace with proper type
	organization: string; // Replace with proper type
	role: string;
	username: string;
}

interface UserResponse {
	self: User;
	urls: URLs;
}

interface NotificationCounts {
	// Add the actual structure based on your API
	[key: string]: number;
}

export const useUserStore = defineStore("user", () => {
	// State
	const display_name = ref<string | null>(null);
	const error = ref<string | null>(null);
	const initials = ref<string | null>(null);
	const loading = ref(false);
	const notification_counts = ref<NotificationCounts | null>(null);
	const notifications = ref<any[] | null>(null);
	const organization = ref<string | null>(null);
	const role = ref<string | null>(null);
	const urls = ref<URLs | null>(null);
	const username = ref<string | null>(null);

	// Subscription
	SSEManager.subscribe("*", (e) => {
		if (e.type !== "heartbeat") {
			fetchUser();
		}
	});

	// Actions
	async function fetchUser() {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/user/self");
			if (!response.ok) throw new Error("Failed to fetch user");

			const data: UserResponse = await response.json();
			display_name.value = data.self.display_name;
			initials.value = data.self.initials;
			notification_counts.value = data.self.notification_counts;
			notifications.value = data.self.notifications;
			organization.value = data.self.organization;
			role.value = data.self.role;
			urls.value = data.urls;
			username.value = data.self.username;
			console.log("loaded user data", data);
		} catch (e) {
			error.value = e instanceof Error ? e.message : "an error ocurred";
			console.error("Error fetching user:", e);
		} finally {
			loading.value = false;
		}
	}

	async function isAuthenticated(): boolean {
		console.log("pretend check user auth");
		return true;
	}
	return {
		// State
		display_name,
		error,
		initials,
		loading,
		notification_counts,
		notifications,
		organization,
		role,
		urls,
		username,
		// Actions
		fetchUser,
		isAuthenticated,
	};
});
