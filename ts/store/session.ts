import * as axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";
import { SSEManager, type SSEMessageResource } from "@/SSEManager";
import {
	Organization,
	Session,
	SessionNotificationCounts,
	URLs,
	User,
} from "@/type/api";
import { apiClient, AxiosErrorJSON } from "@/client";

export class ErrorNotSignedIn extends Error {
	constructor() {
		super("not signed in");
		this.name = "ErrorNotSignedIn";
		Object.setPrototypeOf(this, ErrorNotSignedIn.prototype);
	}
}

export interface SigninResult {
	is_success: boolean;
	status: number;
}
export const useSessionStore = defineStore("session", () => {
	// State
	const hasSession = ref<boolean>(false);
	const isAuthenticated = ref<boolean>(false);
	const isLoading = ref(true);
	const impersonating = ref<string | null>(null);
	const error = ref<string | null>(null);
	const current = ref<Session | null>(null);
	const notification_counts = ref<SessionNotificationCounts | null>(null);
	const ongoingFetch = ref<Promise<Session> | null>(null);
	const organization = ref<Organization | null>(null);
	const self = ref<User | null>(null);
	const urls = ref<URLs | null>(null);

	// Subscription
	SSEManager.subscribe((msg: SSEMessageResource) => {
		if (msg.type == "sync:session") {
			fetchSession();
		}
	});

	// Actions
	async function doSignin(
		password: string,
		username: string,
	): Promise<SigninResult> {
		try {
			console.log("begin signin request");
			await apiClient.JSONPost("/api/signin", {
				password: password,
				username: username,
			});
			isAuthenticated.value = true;
			console.log("set authenticated to true after signin request");
			return {
				is_success: true,
				status: 200,
			};
		} catch (e: any) {
			const data: AxiosErrorJSON =
				e instanceof axios.AxiosError
					? (e.toJSON() as AxiosErrorJSON)
					: { status: 0 };
			if (!data) throw e;
			return {
				is_success: false,
				status: data.status,
			};
		}
	}
	async function fetchSession(): Promise<Session> {
		error.value = null;

		try {
			const data: Session = await apiClient.JSONGet("/api/session");
			isAuthenticated.value = true;
			console.log(
				"set authenticated",
				isAuthenticated.value,
				"due to successful GET /api/session",
			);
			impersonating.value = data.impersonating || null;
			notification_counts.value = data.notification_counts;
			organization.value = data.organization;
			self.value = data.self;
			urls.value = data.urls;
			return data;
		} catch (e: any) {
			const data: AxiosErrorJSON =
				e instanceof axios.AxiosError
					? (e.toJSON() as AxiosErrorJSON)
					: { status: 0 };
			if (data.status == 401) {
				throw new ErrorNotSignedIn();
			}
			console.error("Error fetching session:", e);
			throw e;
		} finally {
			hasSession.value = true;
			isLoading.value = false;
			console.log("no longer loading session");
		}
	}

	async function getAuthenticated(): Promise<boolean> {
		await get();
		return isAuthenticated.value;
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
	async function signout(): Promise<void> {
		isAuthenticated.value = false;
		console.log("set authenticated", isAuthenticated.value, "due to signout");
		apiClient.JSONPost("/api/signout", {});
	}
	return {
		// State
		error,
		getAuthenticated,
		hasSession,
		impersonating,
		isAuthenticated,
		isLoading,
		notification_counts,
		organization,
		self,
		urls,
		// Actions
		doSignin,
		fetchSession,
		get,
		signout,
	};
});
