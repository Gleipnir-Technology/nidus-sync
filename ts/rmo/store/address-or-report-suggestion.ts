import { defineStore } from "pinia";

// Type definitions
interface AddressProperties {
	gid: string;
	name: string;
	coarse_location: string;
	formatted_address_line?: string;
	[key: string]: any; // Allow other properties from the API
}

interface Address {
	properties: AddressProperties;
	geometry?: any;
	type?: string;
}

interface Report {
	id: string;
	type: "nuisance" | "water" | string;
	[key: string]: any; // Allow other properties from the API
}

interface SuggestionsState {
	addresses: Address[];
	reports: Report[];
	loading: boolean;
	error: string | null;
}

interface GeocodeResponse {
	features?: Address[];
	[key: string]: any;
}

interface ReportResponse {
	reports?: Report[];
	[key: string]: any;
}

interface PlaceDetailsResponse {
	features?: Address[];
	[key: string]: any;
}

export const useStoreSuggestion = defineStore("suggestions", {
	state: (): SuggestionsState => ({
		addresses: [],
		reports: [],
		loading: false,
		error: null,
	}),

	actions: {
		async fetchSuggestions(searchText: string): Promise<void> {
			this.loading = true;
			this.error = null;

			try {
				await Promise.all([
					this.fetchAddressSuggestions(searchText),
					this.fetchReportSuggestions(searchText),
				]);
			} catch (error) {
				this.error =
					error instanceof Error ? error.message : "Unknown error occurred";
				console.error("Error fetching suggestions:", error);
			} finally {
				this.loading = false;
			}
		},

		async fetchAddressSuggestions(text: string): Promise<void> {
			try {
				const url = `https://api.stadiamaps.com/geocoding/v2/autocomplete?text=${encodeURIComponent(text)}&focus.point.lat=35&focus.point.lon=-115`;

				const response = await fetch(url);

				if (!response.ok) {
					throw new Error(`Address API error: ${response.status}`);
				}

				const data: GeocodeResponse = await response.json();
				this.addresses = data.features || [];
			} catch (error) {
				console.error("Error fetching geocoding suggestions:", error);
				this.addresses = [];
				throw error;
			}
		},

		async fetchReportSuggestions(text: string): Promise<void> {
			try {
				const url = `/report/suggest?r=${encodeURIComponent(text)}`;
				const response = await fetch(url);

				if (!response.ok) {
					throw new Error(`Report API error: ${response.status}`);
				}

				const data: ReportResponse = await response.json();
				this.reports = data.reports || [];
			} catch (error) {
				console.error("Error fetching report suggestions:", error);
				this.reports = [];
				throw error;
			}
		},

		async fetchAddressDetails(gid: string): Promise<Address | null> {
			try {
				const url = `https://api.stadiamaps.com/geocoding/v2/place_details?ids=${gid}`;
				const response = await fetch(url);

				if (!response.ok) {
					throw new Error(`Address details API error: ${response.status}`);
				}

				const data: PlaceDetailsResponse = await response.json();
				return data.features?.[0] || null;
			} catch (error) {
				console.error("Error fetching address details:", error);
				throw error;
			}
		},

		clearSuggestions(): void {
			this.addresses = [];
			this.reports = [];
			this.error = null;
		},
	},

	getters: {
		hasAddresses: (state): boolean => state.addresses.length > 0,
		hasReports: (state): boolean => state.reports.length > 0,
		hasSuggestions: (state): boolean =>
			state.addresses.length > 0 || state.reports.length > 0,
		totalSuggestions: (state): number =>
			state.addresses.length + state.reports.length,
	},
});

// Export types for use in components
export type { Address, Report, AddressProperties, SuggestionsState };
