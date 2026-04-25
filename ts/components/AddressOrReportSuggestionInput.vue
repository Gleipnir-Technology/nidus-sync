<style scoped>
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
	top: 48px;
}

.suggestion-item {
	cursor: pointer;
	padding: 10px 12px;
	border-bottom: 1px solid #f0f0f0;
}

.suggestion-item:hover {
	background-color: #f8f9fa;
}

.report-id {
	font-weight: 500;
}

.report-type {
	font-size: 0.85rem;
	color: #6c757d;
	margin-top: 2px;
}
</style>
<template>
	<div class="position-relative">
		<div class="input-group">
			<span class="input-group-text">
				<i class="bi bi-search"></i>
			</span>
			<input
				ref="inputRef"
				type="text"
				class="form-control form-control-lg"
				:placeholder="placeholder"
				:value="modelValue"
				maxlength="200"
				@input="handleInput"
			/>
			<div v-if="showSuggestions" class="suggestions-container list-group">
				<!-- Report Suggestions -->
				<div
					v-for="(report, index) in suggestions.reports"
					:key="`report-${report.id}`"
					class="suggestion-item list-group-item"
					@click="handleReportClick(report, index)"
				>
					<div class="report-id">{{ formatReportID(report.id) }}</div>
					<div class="report-type">{{ formatReportType(report.type) }}</div>
				</div>

				<!-- Address Suggestions -->
				<div
					v-for="address in suggestions.addresses"
					:key="`address-${address.properties.gid}`"
					class="suggestion-item list-group-item"
					@click="handleAddressClick(address)"
				>
					<div class="main-address">{{ address.properties.name }}</div>
					<div class="place-info">{{ address.properties.coarse_location }}</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from "vue";
import { useStoreSuggestion } from "@/rmo/store/address-or-report-suggestion";

const props = defineProps({
	modelValue: {
		type: String,
		default: "",
	},
	placeholder: {
		type: String,
		default: "Enter address",
	},
	apiKey: {
		type: String,
		default: "",
	},
});

const emit = defineEmits(["update:modelValue", "suggestion-selected"]);

const inputRef = ref(null);
const debounceTimer = ref(null);
const store = useStoreSuggestion();

const suggestions = computed(() => ({
	addresses: store.addresses,
	reports: store.reports,
}));

const showSuggestions = computed(() => {
	return (
		props.modelValue.length >= 3 &&
		(suggestions.value.addresses.length > 0 ||
			suggestions.value.reports.length > 0)
	);
});

const handleInput = (event) => {
	const searchText = event.target.value;
	emit("update:modelValue", searchText);

	// Clear previous timer
	clearTimeout(debounceTimer.value);

	// Clear suggestions if input is less than 3 characters
	if (searchText.trim().length < 3) {
		store.clearSuggestions();
		return;
	}

	// Debounce API calls (wait 300ms after typing stops)
	debounceTimer.value = setTimeout(() => {
		store.fetchSuggestions(searchText.trim());
	}, 300);
};

const handleReportClick = (report, index) => {
	const formattedId = formatReportID(report.id);
	emit("update:modelValue", formattedId);
	store.clearSuggestions();

	emit("suggestion-selected", {
		content: report,
		type: "report",
	});
};

const handleAddressClick = async (address) => {
	try {
		const detailedAddress = await store.fetchAddressDetails(
			address.properties.gid,
		);

		if (detailedAddress) {
			const formattedAddress =
				detailedAddress.properties.formatted_address_line;
			emit("update:modelValue", formattedAddress);
			store.clearSuggestions();

			emit("suggestion-selected", {
				content: detailedAddress,
				type: "address",
			});
		}
	} catch (error) {
		console.error("Error handling address click:", error);
	}
};

const formatReportID = (id) => {
	if (id.length === 12) {
		return `${id.substring(0, 4)}-${id.substring(4, 8)}-${id.substring(8)}`;
	}
	return id;
};

const formatReportType = (type) => {
	const types = {
		nuisance: "Mosquito Nuisance Report",
		water: "Standing Water Report",
	};
	return types[type] || "Unknown Report Type";
};

const clear = () => {
	emit("update:modelValue", "");
	store.clearSuggestions();
};

const setValue = (suggestion) => {
	if (suggestion?.properties?.formatted_address_line) {
		emit("update:modelValue", suggestion.properties.formatted_address_line);
		store.clearSuggestions();
	}
};

// Expose public methods
defineExpose({
	clear,
	setValue,
});

// Cleanup on unmount
onBeforeUnmount(() => {
	clearTimeout(debounceTimer.value);
});
</script>
