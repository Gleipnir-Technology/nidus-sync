<style scoped>
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
.placeholder-box {
	background-color: #e9ecef;
	border: 2px dashed #adb5bd;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #6c757d;
	font-size: 18px;
	margin-bottom: 20px;
}

.map-placeholder {
	height: 300px;
}

.image-placeholder {
	height: 400px;
}

.action-btn {
	width: 100%;
	margin-bottom: 10px;
	padding: 12px;
	font-size: 16px;
}

.status-badge {
	font-size: 11px;
}

.map-container {
	margin-bottom: 20px;
}
</style>
<template>
	<ThreeColumn>
		<template #left>
			<ReviewSiteColumnList />
		</template>
		<template #center>
			<ReviewSiteColumnDetail />
		</template>
		<template #right>
			<ReviewSiteColumnAction />
		</template>
	</ThreeColumn>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useStoreReviewTask } from "@/store/review-task";
import { useSessionStore } from "@/store/session";
import maplibregl from "maplibre-gl";
import ThreeColumn from "@/components/layout/ThreeColumn.vue";
import ReviewSiteColumnAction from "@/components/ReviewSiteColumnAction.vue";
import ReviewSiteColumnDetail from "@/components/ReviewSiteColumnDetail.vue";
import ReviewSiteColumnList from "@/components/ReviewSiteColumnList.vue";
import { formatAddress } from "@/format";
import type { Changes } from "@/types";
import { Bounds, Contact, Location, ReviewTask } from "@/type/api";
import { Camera } from "@/type/map";
import { MapClickEvent, Marker } from "@/types";

// Props (you can pass these from parent component or environment)
interface Props {}

const props = withDefaults(defineProps<Props>(), {});

// Lifecycle
onMounted(async () => {});
</script>
