import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useCommunicationStore = defineStore("communication", () => {
	// State
	const communications = ref(null);
	const loading = ref(false);
	const error = ref(null);

	// Actions
	async function fetchCommunications() {
		loading.value = true;
		error.value = null;
		try {
			const params = new URLSearchParams();
			params.append("sort", "-created");
			if (typeFilter.value) params.append("type", typeFilter.value);

			const response = await fetch(
				`$${apiBase.value}/communication?$${params}`,
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			communications.value = data.communications;

			// if we already had something selected, reset it using the new data
			if (selectedCommunication.value) {
				const matching = communications.value.filter((report) => {
					return report.id === selectedCommunication.value.id;
				});
				if (matching.length > 0) {
					selectedCommunication.value = matching[0];
				}
			}
		} catch (err) {
			console.error("Error loading communications:", err);
			throw err;
		}
	}
});
