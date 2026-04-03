import { ref } from "vue";
import { defineStore } from "pinia";
import { District } from "@/rmo/type";

export const useDistrictStore = defineStore("district", () => {
	const districts = ref<District[] | null>(null);
	const error = ref<string | null>(null);
	const loading = ref<boolean>(false);
	const ongoingFetch = ref<Promise<District[]> | null>(null);

	async function fetchDistricts(): Promise<District[]> {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/district");
			if (!response.ok) throw new Error("Failed to fetch districts");

			const data: District[] = await response.json();
			districts.value = data;
			return data;
		} catch (e) {
			error.value = e instanceof Error ? e.message : "an error ocurred";
			console.error("Error fetching districts:", e);
			throw new Error(error.value);
		} finally {
			loading.value = false;
		}
	}
	async function get(): Promise<District[]> {
		if (districts.value != null) {
			return districts.value;
		}

		if (ongoingFetch.value !== null) {
			return ongoingFetch.value;
		}

		const s = await fetchDistricts();
		districts.value = s;
		ongoingFetch.value = null;
		return s;
	}
	return {
		error,
		get,
	};
});
