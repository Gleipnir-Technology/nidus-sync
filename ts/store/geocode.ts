import { defineStore } from "pinia";
import { ref } from "vue";
import type { Geocode, Location } from "@/type/api";

export const useStoreGeocode = defineStore("geocode", () => {
	// State
	const loading = ref(false);
	const error = ref(null);

	async function doReverse(url: string, location: Location): Promise<Geocode> {
		loading.value = true;
		error.value = null;
		try {
			//const url = `https://api.stadiamaps.com/geocoding/v2/reverse?point.lat=${location.lat}&point.lon=${location.lng}`;

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
	// Actions
	async function reverse(location: Location): Promise<Geocode> {
		return doReverse("/api/geocode/reverse", location);
	}
	async function reverseClosest(location: Location): Promise<Geocode> {
		return doReverse("/api/geocode/reverse/closest", location);
	}

	return {
		// Actions
		reverse,
		reverseClosest,
	};
});
