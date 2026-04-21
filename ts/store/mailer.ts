import { defineStore } from "pinia";
import { ref } from "vue";
import { useSessionStore } from "@/store/session";

import { apiClient } from "@/client";
import { Mailer, type MailerDTO } from "@/type/api";

export const useStoreMailer = defineStore("publicreport", () => {
	// State
	const _all = ref<Mailer[] | null>(null);
	const _byID = ref<Map<string, Mailer>>(new Map());
	const loading = ref(false);
	const ongoingFetch = ref<Promise<Mailer[]> | null>(null);

	// Actions
	async function byID(id: string): Promise<Mailer> {
		const r = _byID.value.get(id);
		if (r) {
			return r;
		}
		return fetchByID(id);
	}
	async function byURI(uri: string): Promise<Mailer> {
		const id = uri.split("/").pop() || "";
		if (!id) {
			throw new Error(`${uri} is not a recognized public report URI`);
		}
		return byID(id);
	}
	async function fetchAll(): Promise<Mailer[]> {
		const sessionStore = useSessionStore();
		const session = await sessionStore.get();
		loading.value = true;
		const params = new URLSearchParams();
		params.append("sort", "-created");
		const url = `${session.urls.api.mailer}?${params}`;
		const mailers = (await apiClient.JSONGet(url)) as Mailer[];
		_all.value = mailers;
		for (const m of mailers) {
			_byID.value.set(m.id, m);
		}
		return mailers;
	}
	async function fetchByID(id: string): Promise<Mailer> {
		const uri = `/api/publicreport/${id}`;
		return fetchByURI(uri);
	}
	async function fetchByURI(uri: string): Promise<Mailer> {
		loading.value = true;
		try {
			const response = await fetch(uri);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const body: MailerDTO = await response.json();
			const report = Mailer.fromJSON(body);
			_byID.value.set(report.id, report);
			return report;
		} catch (err) {
			console.error("Error loading users:", err);
			throw err;
		}
	}
	async function list(): Promise<Mailer[]> {
		if (_all.value) {
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
	return {
		// Actions
		byID,
		byURI,
		fetchByID,
		fetchByURI,
		list,
	};
});
