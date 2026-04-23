<template>
	<!-- Renderless component -->
</template>

<script setup>
import { inject, onMounted, onBeforeUnmount, watch } from "vue";

const props = defineProps({
	id: { type: String, required: true },
	type: { type: String, required: true },
	tiles: Array,
	url: String,
	// ... other source properties
});

const map = inject("map");
const registerSource = inject("registerSource");
const unregisterSource = inject("unregisterSource");

const getSourceConfig = () => {
	const { id, ...config } = props;
	return config;
};

onMounted(() => {
	registerSource(props.id, getSourceConfig());
});

onBeforeUnmount(() => {
	unregisterSource(props.id);
});

// Watch for prop changes and update source
watch(
	() => getSourceConfig(),
	(newConfig) => {
		if (map.value?.getSource(props.id)) {
			// MapLibre doesn't support updating sources directly
			// You'd need to remove and re-add, or handle specific updates
			unregisterSource(props.id);
			registerSource(props.id, newConfig);
		}
	},
	{ deep: true },
);
</script>
