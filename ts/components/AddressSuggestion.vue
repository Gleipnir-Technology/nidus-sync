<style scoped>
.address-input-wrapper {
	position: relative;
}

.detail-label {
	font-size: 0.8rem;
	text-transform: uppercase;
	color: #6c757d;
	margin-bottom: 2px;
	font-weight: 600;
}

.detail-value {
	font-weight: 500;
}

.main-address {
	font-weight: 500;
}

.place-info {
	font-size: 0.85rem;
	color: #6c757d;
	margin-top: 2px;
}

.suggestions-container {
	position: absolute;
	width: 100%;
	max-height: 300px;
	overflow-y: auto;
	z-index: 3;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	background: white;
}

.suggestion-item {
	cursor: pointer;
	padding: 10px 12px;
	border-bottom: 1px solid #f0f0f0;
}

.suggestion-item:hover {
	background-color: #f8f9fa;
}
</style>
<template>
	<div class="address-input-wrapper">
		<label for="addressInput" class="form-label">Enter address</label>
		<input
			autocomplete="off"
			id="addressInput"
			v-model="searchText"
			class="form-control"
			type="text"
			:placeholder="placeholder"
			maxlength="200"
			@input="handleInput"
		/>
		<div v-if="suggestions.length > 0" class="suggestions-container list-group">
			<div
				v-for="(suggestion, index) in suggestions"
				:key="suggestion.gid"
				class="suggestion-item list-group-item"
				@click="selectSuggestion(suggestion)"
			>
				<div class="main-address">{{ suggestion.detail }}</div>
				<div class="place-info">
					{{ suggestion.type }} {{ suggestion.locality }}
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";

import { formatAddress } from "@/format";
import { Address, GeocodeSuggestion } from "@/type/api";

// Props
interface Props {
	modelValue: Address;
	organizationID?: string;
	placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
	placeholder: "Enter address",
});

// Emits
const emit = defineEmits<{
	"update:modelValue": [value: Address];
	"suggestion-selected": [address: GeocodeSuggestion];
}>();

// State
const searchText = ref("");
const suggestions = ref<GeocodeSuggestion[]>([]);
const debounceTimer = ref<ReturnType<typeof setTimeout> | null>(null);

// Watch for external changes to modelValue
watch(
	() => props.modelValue,
	(newValue) => {
		searchText.value = newValue.raw;
	},
	{ immediate: true },
);

// Methods
function handleInput() {
	const text = searchText.value;

	// Clear previous timer
	if (debounceTimer.value) {
		clearTimeout(debounceTimer.value);
	}

	// Clear suggestions if input is less than 3 characters
	if (text.length < 3) {
		suggestions.value = [];
		return;
	}

	// Debounce API calls (wait 300ms after typing stops)
	debounceTimer.value = setTimeout(async () => {
		await fetchAddressSuggestions(text);
	}, 300);

	// Update the model
	emit("update:modelValue", {
		country: "",
		gid: props.modelValue.gid,
		locality: "",
		number: "",
		postal_code: "",
		raw: text,
		region: "",
		street: "",
		unit: "",
	});
}

async function fetchAddressSuggestions(text: string) {
	try {
		const q = encodeURIComponent(text);
		let url = `/api/geocode/suggestion?query=${q}`;
		if (props.organizationID) {
			url = url + "&org=" + props.organizationID;
		}

		const response = await fetch(url);
		const data = await response.json();
		suggestions.value = data || [];
	} catch (error) {
		console.error("Error fetching geocoding suggestions:", error);
		suggestions.value = [];
	}
}

async function selectSuggestion(suggestion: GeocodeSuggestion) {
	try {
		// Update display text
		searchText.value = suggestion.detail + ", " + suggestion.locality;

		// Clear suggestions
		suggestions.value = [];

		// Emit the full address object
		emit("update:modelValue", {
			country: "",
			gid: suggestion.gid,
			locality: "",
			number: "",
			postal_code: "",
			raw: suggestion.detail,
			region: "",
			street: "",
			unit: "",
		});
		emit("suggestion-selected", suggestion);
	} catch (error) {
		console.error("Error fetching place details:", error);
	}
}
onMounted(() => {
	// We may get an address from the API which doesn't have a raw value because it wasn't
	// ever typed or it matched an address from a vendor's database. In that case, populate it,
	// so we see something nice
	if (props.modelValue.raw == "" && props.modelValue.gid != "") {
		const raw = formatAddress(props.modelValue);
		searchText.value = raw;
		console.log("Set raw address", raw, props.modelValue);
	}
});
</script>
