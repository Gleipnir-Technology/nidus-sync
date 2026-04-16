<style scoped>
#address-input {
	font-size: 16px;
}
.map-container {
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	margin-bottom: 20px;
	margin-top: 20px;
	align-items: center;
	justify-content: center;
	/* Prevent touch scrolling issues */
	touch-action: pan-y pinch-zoom;
}
#map {
	width: 100%;
	height: 100%;
}

/* Mobile-specific adjustments */
@media (max-width: 768px) {
	.map-container {
		height: 400px;
		margin-bottom: 15px;
		margin-top: 15px;
	}
}

/* Extra small devices */
@media (max-width: 576px) {
	.map-container {
		height: 350px;
		border-radius: 5px;
	}
}
</style>
<template>
	<div class="mb-4">
		<AddressSuggestion
			v-model="modelValue"
			placeholder="Start typing an address (min 3 characters)"
			@suggestion-selected="doAddressSuggestionSelected"
		>
		</AddressSuggestion>
	</div>

	<!-- Map Placeholder -->
	<div class="mb-4">
		<label class="form-label fw-semibold">Location Preview</label>
		<div class="map-container">
			<MapLocator
				:initialCamera="initialCamera"
				:markers="markers"
				@click="doMapClick"
				@marker-drag-end="doMapMarkerDragEnd"
				v-model="currentCamera"
			/>
		</div>
	</div>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import AddressSuggestion from "@/components/AddressSuggestion.vue";
import MapLocator from "@/components/MapLocator.vue";
import { Address } from "@/type/api";
import type { Geocode, GeocodeSuggestion, Location } from "@/type/api";
import { useStoreGeocode } from "@/store/geocode";
import { Camera } from "@/type/map";
import type { MapClickEvent, Marker } from "@/types";

interface Emits {
	(e: "update:modelValue", value: Address): void;
}
interface Props {
	initialCamera?: Camera;
	modelValue: Address;
}
const address = ref<string>("");
const currentCamera = ref<Camera>(new Camera());
const emit = defineEmits<Emits>();
const geocode = useStoreGeocode();
const markers = computed((): Marker[] => {
	if (!props.modelValue.location) {
		return [];
	}
	if (
		props.modelValue.location.latitude == 0.0 ||
		props.modelValue.location.longitude == 0.0
	) {
		return [];
	}
	const marker = {
		color: "#FF0000",
		draggable: true,
		id: "x",
		location: props.modelValue.location,
	};
	return [marker];
});
const modelValue = computed({
	get: () => props.modelValue,
	set: (value: Address) => emit("update:modelValue", value),
});
const props = defineProps<Props>();
function doAddressSuggestionSelected(suggestion: GeocodeSuggestion) {
	console.log("Address suggestion selected", suggestion);

	doAddressSuggestionDetails(suggestion);
}
async function doAddressSuggestionDetails(suggestion: GeocodeSuggestion) {
	// Fetch full details for the selected suggestion
	updateModel(suggestion.gid, suggestion.detail, props.modelValue.location);
	const url = `/api/geocode/by-gid/${suggestion.gid}`;
	const response = await fetch(url);
	if (!response.ok) {
		console.error("Failed to get suggestion detail", response.statusText);
		return;
	}
	const data = (await response.json()) as Geocode;

	if (currentCamera.value) {
		console.log("suggestion located, zooming", data);
		currentCamera.value.zoom = 15;
	}
	updateModel(data.address.gid, data.address.raw, data.address.location);
}
function doMapClick(event: MapClickEvent) {
	updateModel(props.modelValue.gid, props.modelValue.raw, event.location);
	geocode
		.reverse(event.location)
		.then((code: Geocode) => {
			updateModel(
				code.address.gid,
				code.address.raw,
				props.modelValue.location,
			);
			console.log("reverse geocoded", code);
		})
		.catch((e) => {
			console.error("failed to reverse geocode after map click", e);
		});
}
function doMapMarkerDragEnd(location: Location) {
	updateModel(props.modelValue.gid, props.modelValue.raw, location);
}
function updateModel(
	address_gid: string,
	address_raw: string,
	location?: Location,
) {
	const newAddress = new Address(
		"",
		address_gid,
		"",
		"",
		"",
		address_raw,
		"",
		"",
		"",
		location,
	);
	emit("update:modelValue", newAddress);
}
</script>
