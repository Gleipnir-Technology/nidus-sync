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
	z-index: 1000;
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
				:key="suggestion.properties.gid || index"
				class="suggestion-item list-group-item"
				@click="selectSuggestion(suggestion)"
			>
				<div class="main-address">{{ suggestion.properties.name }}</div>
				<div class="place-info">
					{{ suggestion.properties.coarse_location }}
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { Address } from "@/type/stadia";

// Props
interface Props {
	modelValue?: Address | null;
	placeholder?: string;
	apiKey?: string;
}

const props = withDefaults(defineProps<Props>(), {
	modelValue: null,
	placeholder: "Enter address",
	apiKey: "",
});

// Emits
const emit = defineEmits<{
	"update:modelValue": [value: Address | null];
	"address-selected": [address: Address];
}>();

// State
const searchText = ref("");
const suggestions = ref<Address[]>([]);
const debounceTimer = ref<ReturnType<typeof setTimeout> | null>(null);

// Watch for external changes to modelValue
watch(
	() => props.modelValue,
	(newValue) => {
		if (newValue) {
			searchText.value = formatAddressDisplay(newValue);
		} else {
			searchText.value = "";
		}
	},
	{ immediate: true },
);

// Methods
function handleInput() {
	const text = searchText.value.trim();

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
}

async function fetchAddressSuggestions(text: string) {
	try {
		const url = `https://api.stadiamaps.com/geocoding/v2/autocomplete?text=${encodeURIComponent(
			text,
		)}&focus.point.lat=35&focus.point.lon=-115`;

		const response = await fetch(url);
		const data = await response.json();
		suggestions.value = data.features || [];
	} catch (error) {
		console.error("Error fetching geocoding suggestions:", error);
		suggestions.value = [];
	}
}

async function selectSuggestion(suggestion: Address) {
	try {
		// Fetch full details for the selected suggestion
		const url = `https://api.stadiamaps.com/geocoding/v2/place_details?ids=${suggestion.properties.gid}`;
		const response = await fetch(url);
		const data = await response.json();
		const fullAddress: Address = data.features[0];

		// Update display text
		searchText.value = formatAddressDisplay(fullAddress);

		// Clear suggestions
		suggestions.value = [];

		// Emit the full address object
		emit("update:modelValue", fullAddress);
		emit("address-selected", fullAddress);
	} catch (error) {
		console.error("Error fetching place details:", error);
	}
}

function formatAddressDisplay(address: Address): string {
	const props = address.properties;

	if (props.formatted_address_line) {
		return props.formatted_address_line;
	} else if (props.address_components) {
		const num = props.address_components.number ?? "";
		const street = props.address_components.street ?? "";
		const location = props.coarse_location ?? "";
		return `${num} ${street}, ${location}`.trim();
	} else {
		return `${props.name ?? ""}, ${props.coarse_location ?? ""}`.trim();
	}
}
</script>
