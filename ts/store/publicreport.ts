import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import {
	PublicReport,
	type PublicReportComplianceCreateRequest,
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
	async function byID(id: string): Promise<PublicReport> {
		const r = _byID.value.get(id);
		if (r) {
			return r;
		}
		return fetchByID(id);
	}
	async function byURI(uri: string): Promise<PublicReport> {
		const id = uri.split("/").pop() || "";
		if (!id) {
			throw new Error(`${uri} is not a recognized public report URI`);
		}
		return byID(id);
	}
	async function createCompliance(
		data: PublicReportComplianceCreateRequest,
	): Promise<PublicReport> {
		const resp = (await apiClient.JSONPost(
			"/api/rmo/compliance",
			data,
		)) as PublicReportDTO;
		const result = PublicReport.fromJSON(resp);
		_byID.value.set(result.public_id, result);
		return result;
	}
	async function fetchByID(id: string): Promise<PublicReport> {
		const uri = `/api/publicreport/${id}`;
		return fetchByURI(uri);
	}
	async function fetchByURI(uri: string): Promise<PublicReport> {
		loading.value = true;
		try {
			const response = await fetch(uri);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const body: PublicReportDTO = await response.json();
			const report = PublicReport.fromJSON(body);
			_byID.value.set(report.public_id, report);
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
		byURI,
		createCompliance,
		fetchByID,
		fetchByURI,
	};
});
