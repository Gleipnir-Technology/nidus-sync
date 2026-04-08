import { defineStore } from "pinia";
import { ref } from "vue";
import type { Publicreport } from "@/type/api";

export const useStorePublicreport = defineStore("publicreport", () => {
	// State
	const _byID = ref<Map<string, Publicreport>>(new Map());
	const error = ref(null);
	const loading = ref(false);
	//const ongoingFetch = ref<Promise<Publicreport[]> | null>(null);

	function add(pr: Publicreport) {
		_byID.value.set(pr.id, pr);
	}
	// Actions
	async function byID(id: string): Promise<Publicreport | undefined> {
		loading.value = true;
		error.value = null;
		try {
			const url = `/api/publicreport/${id}`;
			const response = await fetch(url);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const body = await response.json();
			_byID.value.set(id, body);
			return body;
		} catch (err) {
			console.error("Error loading users:", err);
			throw err;
		}
	}

	return {
		// Actions
		add,
		byID,
	};
});
