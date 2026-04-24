<template>
	<!-- Renderless component -->
</template>

<script setup lang="ts">
import maplibregl from "maplibre-gl";
import { inject, onMounted, onBeforeUnmount, Ref, useAttrs, watch } from "vue";

export type MouseEvent = maplibregl.MapLayerMouseEvent;
type LayerType = maplibregl.LayerSpecification["type"];
interface Emits {
	(e: "click", evt: MouseEvent): void;
	(e: "mouseenter"): void;
	(e: "mouseleave"): void;
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

type OnCallbackFunc = (e?: MouseEvent) => void;
type RegisterOnFunc = (
	eventname: string,
	layerid: string,
	callback: OnCallbackFunc,
) => void;
type RegisterLayerFunc = (id: string, config: any) => void;
type UnregisterLayerFunc = (id: string) => void;

const map: Ref<maplibregl.Map | null> | undefined = inject("map");
const registerOn: RegisterOnFunc | undefined = inject("registerOn");
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
