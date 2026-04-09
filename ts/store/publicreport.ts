import { defineStore } from "pinia";
import { ref } from "vue";
import { PublicReport, type PublicReportDTO } from "@/type/api";

export const useStorePublicReport = defineStore("publicreport", () => {
	// State
	const _byID = ref<Map<string, PublicReport>>(new Map());
	const error = ref(null);
	const loading = ref(false);
	//const ongoingFetch = ref<Promise<PublicReport[]> | null>(null);

	function add(pr: PublicReport) {
		_byID.value.set(pr.id, pr);
	}
	// Actions
	async function byID(id: string): Promise<PublicReport | undefined> {
		loading.value = true;
		error.value = null;
		try {
			const url = `/api/publicreport/${id}`;
			const response = await fetch(url);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const body: PublicReportDTO = await response.json();
			const report = PublicReport.fromJSON(body);
			_byID.value.set(id, report);
			return report;
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
