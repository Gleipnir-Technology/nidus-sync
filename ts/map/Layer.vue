<template>
	<!-- Renderless component -->
</template>

<script setup lang="ts">
import maplibregl from "maplibre-gl";
import { inject, onMounted, onBeforeUnmount, Ref, watch } from "vue";

type LayerType = maplibregl.LayerSpecification["type"];
export interface Props {
	filter?: maplibregl.FilterSpecification;
	id: string;
	minzoom?: number;
	paint: Object;
	source: string;
	sourceLayer: string;
	type: LayerType;
}
const props = withDefaults(defineProps<Props>(), {});

type RegisterLayerFunc = (id: string, config: any) => void;
type UnregisterLayerFunc = (id: string) => void;
const map: Ref<maplibregl.Map | null> | undefined = inject("map");
const registerLayer: RegisterLayerFunc | undefined = inject("registerLayer");
const unregisterLayer: UnregisterLayerFunc | undefined =
	inject("unregisterLayer");

const getLayerConfig = (): maplibregl.LayerSpecification => {
	let result: maplibregl.LayerSpecification = {
		id: props.id,
		source: props.source,
		"source-layer": props.sourceLayer,
		type: props.type,
		...(props.filter && { filter: props.filter }),
		...(props.minzoom && { minzoom: props.minzoom }),
		...(props.paint && { paint: props.paint }),
	} as maplibregl.LayerSpecification;
	return result;
};

onMounted(() => {
	if (registerLayer) {
		registerLayer(props.id, getLayerConfig());
	} else {
		console.log("registerLayer is nully");
	}
});

onBeforeUnmount(() => {
	if (unregisterLayer) {
		unregisterLayer(props.id);
	} else {
		console.log("unregisterLayer is nully");
	}
});

// Update paint/layout properties reactively
watch(
	() => props.paint,
	(newPaint) => {
		if (map && map.value?.getLayer(props.id)) {
			Object.entries(newPaint || {}).forEach(([key, value]) => {
				if (!(map && map.value)) return;
				map.value.setPaintProperty(props.id, key, value);
			});
		}
	},
	{ deep: true },
);
</script>
