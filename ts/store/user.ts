import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useUserStore = defineStore("user", () => {
	// State
	const display_name = ref(null);
	const error = ref(null);
	const initials = ref(null);
	const loading = ref(false);
	const notification_counts = ref(null);
	const notifications = ref(null);
	const organization = ref(null);
	const role = ref(null);
	const urls = ref(null);
	const username = ref(null);

	// Getters
	const isAuthenticated = computed(() => user.value !== null);
	const userName = computed(() => user.value?.name ?? "");

	// Actions
	async function fetchUser() {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/user/self");
			if (!response.ok) throw new Error("Failed to fetch user");

			const data = await response.json();
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
			error.value = e.message;
			console.error("Error fetching user:", e);
		} finally {
			loading.value = false;
		}
	}

	function clearUser() {
		user.value = null;
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
		// Getters
		isAuthenticated,
		userName,
		// Actions
		fetchUser,
		clearUser,
	};
});
