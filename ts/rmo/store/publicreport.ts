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
	const _byURI = ref<Map<string, PublicReport>>(new Map());
	const loading = ref(false);
	const ongoingFetches = ref<Map<string, Promise<PublicReport> | null>>(
		new Map(),
	);

	function add(pr: PublicReport) {}
	async function byID(id: string): Promise<PublicReport> {
		const uri = "/api/rmo/publicreport/" + id;
		return byURI(uri);
	}
	async function byURI(uri: string): Promise<PublicReport> {
		let result = _byURI.value.get(uri);
		if (result) return result;
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
		_byURI.value.set(result.uri, result);
		return result;
	}
	async function fetchByURI(uri: string): Promise<PublicReport> {
		loading.value = true;
		try {
			const body = (await apiClient.JSONGet(uri)) as PublicReportDTO;
			const report = PublicReport.fromJSON(body);
			_byURI.value.set(report.uri, report);
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
		byID,
		byURI,
		createCompliance,
		fetchByURI,
		update,
	};
});
