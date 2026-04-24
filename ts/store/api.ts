import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import { APIProperties } from "@/type/api";

export const useStoreAPI = defineStore("api", () => {
	// State
	const _response = ref<APIProperties | null>(null);
	const loading = ref(false);
	const ongoingFetch = ref<Promise<APIProperties> | null>(null);

	// Actions
	async function doFetch(): Promise<APIProperties> {
		loading.value = true;
		const url = "/api";
		const resp = (await apiClient.JSONGet(url)) as APIProperties;
		return resp;
	}
	async function get(): Promise<APIProperties> {
		if (_response.value) {
			return _response.value;
		}
		if (ongoingFetch.value !== null) {
			return ongoingFetch.value;
		}
		ongoingFetch.value = doFetch().finally(() => {
			ongoingFetch.value = null;
		});
		return ongoingFetch.value;
	}
	return {
		// Actions
		get,
	};
});
