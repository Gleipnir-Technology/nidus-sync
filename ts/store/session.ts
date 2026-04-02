import { defineStore } from "pinia";
import { ref } from "vue";
import { SSEManager, type SSEMessage } from "@/SSEManager";
import {
	Organization,
	Session,
	SessionNotificationCounts,
	URLs,
	User,
} from "@/types";

export const useSessionStore = defineStore("session", () => {
	// State
	const impersonating = ref<string | null>(null);
	const error = ref<string | null>(null);
	const loading = ref(false);
	const current = ref<Session | null>(null);
	const notification_counts = ref<SessionNotificationCounts | null>(null);
	const ongoingFetch = ref<Promise<Session> | null>(null);
	const organization = ref<Organization | null>(null);
	const self = ref<User | null>(null);
	const urls = ref<URLs | null>(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessage) => {
		if (msg.type !== "sync:session") {
			fetchSession();
		}
	});

	// Actions
	async function fetchSession(): Promise<Session> {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/session");
			if (!response.ok) throw new Error("Failed to fetch user");

			const data: Session = await response.json();
			impersonating.value = data.impersonating || null;
			notification_counts.value = data.notification_counts;
			organization.value = data.organization;
			self.value = data.self;
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

	async function get(): Promise<Session> {
		if (current.value != null) {
			return current.value;
		}

		if (ongoingFetch.value !== null) {
			return ongoingFetch.value;
		}

		const s = await fetchSession();
		current.value = s;
		ongoingFetch.value = null;
		return s;
	}
	return {
		// State
		error,
		impersonating,
		loading,
		notification_counts,
		organization,
		self,
		urls,
		// Actions
		fetchSession,
		get,
		isAuthenticated,
	};
});
