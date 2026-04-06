import { defineStore } from "pinia";
import { ref } from "vue";
import type { Geocode, Location } from "@/type/api";

export const useGeocodeStore = defineStore("geocode", () => {
	// State
	const loading = ref(false);
	const error = ref(null);

	// Actions
	async function reverse(location: Location): Promise<Geocode> {
		loading.value = true;
		error.value = null;
		try {
			//const url = `https://api.stadiamaps.com/geocoding/v2/reverse?point.lat=${location.lat}&point.lon=${location.lng}`;

			const url = "/api/geocode/reverse";
			const response = await fetch(url, {
				body: JSON.stringify(location),
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
			});
			const data = (await response.json()) as Geocode;
			return data;
		} catch (err) {
			console.error("Error loading signals:", err);
			throw err;
		}
	}

	return {
		// Actions
		reverse,
	};
});
