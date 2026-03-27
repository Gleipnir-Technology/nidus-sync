<style scoped>
.selected {
	background-color: $primary;
}
</style>

<template>
	<div class="row" :class="selected ? 'selected' : ''">
		<div class="col-1">
			<i class="bi" :class="icon(signal)"></i>
		</div>
		<div class="col-5">
			<div class="small fw-semibold">{{ title(signal) }}</div>
		</div>
		<div class="col-6">
			<div class="text-muted small fst-italic">
				{{ location(signal) }}
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { shortAddress } from "../format";
interface Props {
	selected: boolean;
	signal: Signal;
}
const props = defineProps<Props>();
function icon(signal: Signal): string {
	if (signal.type == "flyover pool") {
		return "bi-pond";
	} else if (signal.type == "publicreport nuisance") {
		return "bi-mosquito";
	} else if (signal.type == "publicreport water") {
		return "bi-water";
	} else {
		return "bi-mosquito";
	}
}
function location(signal: Signal): string {
	if (signal.address != null) {
		return shortAddress(signal.address);
	} else {
		return `${signal.location.latitude}, ${signal.location.longitude}`;
	}
}
function title(signal: Signal): string {
	if (signal.type == "flyover pool") {
		return "Green pool";
	} else if (signal.type == "publicreport nuisance") {
		return "Nuisance";
	} else if (signal.type == "publicreport water") {
		return "Standing water";
	}
}
</script>
