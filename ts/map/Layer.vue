<template>
	<!-- Renderless component -->
</template>

<script setup lang="ts">
import maplibregl from "maplibre-gl";
import { inject, onMounted, onBeforeUnmount, Ref, useAttrs, watch } from "vue";

export type MapEventType = maplibregl.MapEventType;
export type MouseEvent = maplibregl.MapLayerMouseEvent;
export type Feature = maplibregl.MapGeoJSONFeature;
type LayerType = maplibregl.LayerSpecification["type"];
interface Emits {
	(e: "click", evt: MouseEvent): void;
	(e: "mouseenter"): void;
	(e: "mouseleave"): void;
	(e: "update:modelValue", features: Feature[]): void;
}
export interface Props {
	filter?: maplibregl.FilterSpecification;
	id: string;
	minzoom?: number;
	paint: Object;
	source: string;
	sourceLayer: string;
	type: LayerType;
}
const attrs = useAttrs();
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {});

type OnCallbackFunc = (e?: any) => void;
type RegisterOnFunc = (
	eventname: string,
	layerid: string,
	callback: OnCallbackFunc,
) => void;
type RegisterLayerFunc = (id: string, config: any) => void;
type UnregisterLayerFunc = (id: string) => void;

const map: Ref<maplibregl.Map | null> | undefined = inject("map");
const registerOn: RegisterOnFunc | undefined = inject("registerOn");
const registerOnce: RegisterOnFunc | undefined = inject("registerOnce");
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

function updateModel() {
	if (!(map && map.value)) return;
	const query: maplibregl.QueryRenderedFeaturesOptions = {
		layers: [props.id],
	};
	const features = map.value.queryRenderedFeatures(query);
	const features_from_source = features.filter(
		(feature: any) => feature.source == props.source,
	);
	emit("update:modelValue", features_from_source);
	//emit("mouseleave");
}
onMounted(() => {
	if (registerLayer) {
		registerLayer(props.id, getLayerConfig());
	}
	if (registerOn) {
		registerOn("click", props.id, (e?: MouseEvent) => {
			if (e) {
				emit("click", e);
			}
		});
	}
	if (registerOn) {
		registerOn("mouseenter", props.id, () => {
			emit("mouseenter");
		});
	}
	if (registerOn) {
		registerOn("mouseleave", props.id, () => {
			emit("mouseleave");
		});
	}
	if (registerOn) {
		registerOn("moveend", props.id, updateModel);
	}
	if (registerOnce) {
		registerOnce("idle", props.id, updateModel);
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
