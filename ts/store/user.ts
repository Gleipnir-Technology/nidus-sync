import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useUserStore = defineStore("user", () => {
	// State
	const user = ref(null);
	const loading = ref(false);
	const error = ref(null);

	// Getters
	const isAuthenticated = computed(() => user.value !== null);
	const userName = computed(() => user.value?.name ?? "");
	const organization = computed(() => user.value?.organization ?? "");

	// Actions
	async function fetchUser() {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/user/self");
			if (!response.ok) throw new Error("Failed to fetch user");

			user.value = await response.json();
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
		user,
		loading,
		error,
		// Getters
		isAuthenticated,
		userName,
		organization,
		// Actions
		fetchUser,
		clearUser,
	};
});
