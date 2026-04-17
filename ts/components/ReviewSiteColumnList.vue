<style scoped lang="scss">
.entry-item {
	padding: 15px;
	border-bottom: 1px solid #e9ecef;
	cursor: pointer;
	transition: background-color 0.2s;
}

.entry-item:hover {
	background-color: #f8f9fa;
}

.entry-item.active {
	background-color: #e7f3ff;
	border-left: 4px solid #0d6efd;
}
</style>
<template>
	<div class="p-3 border-bottom bg-primary text-white">
		<h5 class="mb-0"><i class="bi bi-list-ul"></i> Sites</h5>
	</div>
	<div
		v-for="site in sites"
		:key="site.id"
		class="entry-item"
		:class="{ active: selectedSite?.id === site.id }"
		@click="doClick(site)"
	>
		<div class="d-flex">
			<div>
				<i
					class="bi"
					:class="{
						'bi-house-fill': site.leads.length > 0,
						'bi-house': site.leads.length == 0,
					}"
				></i>
			</div>
			<strong>{{ formatAddress(site.address) }}</strong>
		</div>
	</div>
</template>
<script setup lang="ts">
import { Site } from "@/type/api";
import { formatAddress } from "@/format";
interface Emits {
	(e: "doSiteDeselect", id: number): void;
	(e: "doSiteSelect", id: number): void;
}
interface Props {
	selectedSite: Site | undefined;
	sites: Site[];
}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {});
function doClick(site: Site) {
	if (props.selectedSite && site.id == props.selectedSite.id) {
		emit("doSiteDeselect", site.id);
	} else {
		emit("doSiteSelect", site.id);
	}
}
</script>
