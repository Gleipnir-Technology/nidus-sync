import { defineStore } from "pinia";
import { ref } from "vue";

import { apiClient } from "@/client";
import {
	PublicReport,
	type PublicReportComplianceCreateRequest,
	type PublicReportDTO,
	type PublicReportUpdate,
} from "@/type/api";

export const useStorePublicReport = defineStore("publicreport", () => {
	// State
	const _byID = ref<Map<string, PublicReport>>(new Map());
	const loading = ref(false);
	const ongoingFetches = ref<Map<string, Promise<PublicReport> | null>>(
		new Map(),
	);

	function add(pr: PublicReport) {
		_byID.value.set(pr.public_id, pr);
	}
	async function byID(id: string): Promise<PublicReport> {
		const uri = "/api/rmo/publicreport/" + id;
		return byURI(uri);
	}
	async function byURI(uri: string): Promise<PublicReport> {
		let ongoing = ongoingFetches.value.get(uri);
		if (ongoing) return ongoing;
		ongoing = fetchByURI(uri).finally(() => {
			ongoingFetches.value.set(uri, null);
		});
		ongoingFetches.value.set(uri, ongoing);
		return ongoing;
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
	async function fetchByURI(uri: string): Promise<PublicReport> {
		loading.value = true;
		try {
			const body = (await apiClient.JSONGet(uri)) as PublicReportDTO;
			const report = PublicReport.fromJSON(body);
			_byID.value.set(report.public_id, report);
			return report;
		} catch (err) {
			console.error("Error loading users:", err);
			throw err;
		}
	}
	async function update(
		uri: string,
		updates: PublicReportUpdate,
	): Promise<PublicReport> {
		const resp = (await apiClient.JSONPut(uri, updates)) as PublicReportDTO;
		return PublicReport.fromJSON(resp);
	}
	return {
		// Actions
		add,
		byID,
		byURI,
		createCompliance,
		fetchByURI,
		update,
	};
});
