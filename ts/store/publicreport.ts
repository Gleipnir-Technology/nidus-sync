import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import {
	PublicReport,
	type PublicReportCreateRequest,
	type PublicReportDTO,
} from "@/type/api";

export const useStorePublicReport = defineStore("publicreport", () => {
	// State
	const _byID = ref<Map<string, PublicReport>>(new Map());
	const loading = ref(false);
	//const ongoingFetch = ref<Promise<PublicReport[]> | null>(null);

	function add(pr: PublicReport) {
		_byID.value.set(pr.public_id, pr);
	}
	// Actions
	async function byID(id: string): Promise<PublicReport | undefined> {
		loading.value = true;
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
	async function byURI(uri: string): Promise<PublicReport | undefined> {
		const id = uri.split("/").pop() || "";
		if (!id) {
			throw new Error(`${uri} is not a recognized public report URI`);
		}
		return byID(id);
	}
	async function create(
		data: PublicReportCreateRequest,
	): Promise<PublicReport> {
		const resp = (await apiClient.JSONPost(
			"/api/rmo/compliance",
			data,
		)) as PublicReportDTO;
		const result = PublicReport.fromJSON(resp);
		_byID.value.set(result.public_id, result);
		return result;
	}
	return {
		// Actions
		add,
		byID,
		byURI,
		create,
	};
});
