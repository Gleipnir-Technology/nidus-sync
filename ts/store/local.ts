import { defineStore } from "pinia";

export const useStoreLocal = defineStore("local", () => {
	function getClientID(): string {
		let id = localStorage.getItem("session_id");
		if (id) {
			return id;
		}
		id = crypto.randomUUID();
		localStorage.setItem("session_id", id.toString());
		return id;
	}
	function getExistingComplianceReportURI(): string | null {
		return localStorage.getItem("working_compilance_report_uri");
	}
	function setExistingComplianceReportURI(uri: string) {
		localStorage.setItem("working_compilance_report_uri", uri);
	}
	return {
		getClientID,
		getExistingComplianceReportURI,
		setExistingComplianceReportURI,
	};
});
