import { defineStore } from "pinia";
import { ref } from "vue";
import type { District } from "@/type/api";

export const useStoreDistrict = defineStore("district", () => {
	// State
	const _byURI = ref<Map<string, District>>(new Map());
	const error = ref<string | null>(null);
	const loading = ref<boolean>(false);
	const ongoingFetch = ref<Promise<District[]> | null>(null);

	// Actions
	async function byURI(uri: string): Promise<District | undefined> {
		let district = _byURI.value.get(uri);
		console.log("district by uri", uri, district);
		if (district) {
			return district;
		}
		loading.value = true;
		error.value = null;
		try {
			const response = await fetch(uri);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const body = await response.json();
			_byURI.value.set(uri, body);
			return body;
		} catch (e) {
			console.error("Error loading users:", e);
			error.value = e instanceof Error ? e.message : "an error ocurred";
			throw e;
		} finally {
			loading.value = false;
		}
	}
	async function fetchDistricts(): Promise<District[]> {
		loading.value = true;
		error.value = null;

		try {
			const response = await fetch("/api/district");
			if (!response.ok) throw new Error("Failed to fetch districts");

			const data: District[] = await response.json();
			data.forEach((d: District) => {
				_byURI.value.set(d.uri, d);
				console.log("district", d.uri);
			});
			return data;
		} catch (e) {
			error.value = e instanceof Error ? e.message : "an error ocurred";
			console.error("Error fetching districts:", e);
			throw new Error(error.value);
		} finally {
			loading.value = false;
		}
	}
	async function list(): Promise<District[]> {
		if (_byURI.value.size > 0) {
			return Array.from(_byURI.value.values());
		}

		if (ongoingFetch.value !== null) {
			return ongoingFetch.value;
		}

		const s = await fetchDistricts();
		ongoingFetch.value = null;
		return s;
	}

	return {
		// Actions
		byURI,
		list,
	};
});
