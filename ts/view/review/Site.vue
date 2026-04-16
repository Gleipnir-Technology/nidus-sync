<style scoped lang="scss">
body {
	background-color: #f5f5f5;
}

.left-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	border-right: 1px solid #dee2e6;
}

.middle-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	padding: 20px;
}

.right-panel {
	background-color: white;
	height: 100vh;
	overflow-y: auto;
	border-left: 1px solid #dee2e6;
	padding: 20px;
}
</style>
<template>
	<ThreeColumn>
		<template #left>
			<ReviewSiteColumnList
				@doSiteDeselect="siteDeselect"
				@doSiteSelect="siteSelect"
				:selectedSite="selectedSite"
				:sites="storeSite.all()"
			/>
		</template>
		<template #center>
			<ReviewSiteColumnDetail
				:mapFlyoverCamera="mapFlyoverCamera"
				:mapMarkers="mapMarkers"
				:selectedSite="selectedSite"
			/>
		</template>
		<template #right>
			<ReviewSiteColumnAction
				:selectedSite="selectedSite"
				:submitting="submitting"
			/>
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useStoreSite } from "@/store/site";
import { useSessionStore } from "@/store/session";
import maplibregl from "maplibre-gl";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ReviewSiteColumnAction from "@/components/ReviewSiteColumnAction.vue";
import ReviewSiteColumnDetail from "@/components/ReviewSiteColumnDetail.vue";
import ReviewSiteColumnList from "@/components/ReviewSiteColumnList.vue";
import { formatAddress } from "@/format";
import type { Changes } from "@/types";
import { Bounds, Contact, Location, Site } from "@/type/api";
import { Camera } from "@/type/map";
import { MapClickEvent, Marker } from "@/types";

// Props (you can pass these from parent component or environment)
interface Props {}

const props = withDefaults(defineProps<Props>(), {});

const mapFlyoverCamera = ref<Camera>(new Camera());
const storeSite = useStoreSite();
const selectedSiteID = ref<number>(0);
const submitting = ref<boolean>(false);
const selectedSite = computed((): Site | undefined => {
	if (!selectedSiteID.value) {
		return undefined;
	}
	return storeSite.byID(selectedSiteID.value);
});
const mapMarkers = computed<Marker[]>(() => {
	const site = selectedSite.value;
	if (!(site && site.address.location)) {
		return [];
	}
	const markers = {
		color: "#FF0000",
		draggable: false,
		id: "address",
		location: site.address.location,
	};
	return [markers];
});
function siteDeselect(id: number): void {
	if (selectedSiteID.value == id) {
		selectedSiteID.value = 0;
	}
}
function siteSelect(id: number): void {
	selectedSiteID.value = id;

	const site = storeSite.byID(id);
	if (!site) {
		console.log("no site", id);
		return;
	}
	mapFlyoverCamera.value = new Camera(site.address.location, 20);
	console.log("selecting site", id, site);
}

// Lifecycle
onMounted(async () => {
	storeSite.fetchAll();
});
</script>
