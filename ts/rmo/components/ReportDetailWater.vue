<template>
	<ReportDetailEntry
		label="Has a gate that affects access?"
		:value="water.access_gate.toString()"
	/>
	<ReportDetailEntry
		label="Has dog that affects access?"
		:value="water.access_dog.toString()"
	/>
	<ReportDetailEntry
		label="Has a fence that affects access?"
		:value="water.access_fence.toString()"
	/>
	<ReportDetailEntry
		label="Has a locked entracnce that affects access?"
		:value="water.access_locked.toString()"
	/>
	<ReportDetailEntry
		label="Reporter observed larvae (wigglers)?"
		:value="water.has_larvae.toString()"
	/>
	<ReportDetailEntry
		label="Reporter observed pupae (tumblers)?"
		:value="water.has_pupae.toString()"
	/>
	<ReportDetailEntry
		label="Reporter observed adult mosquitoes?"
		:value="water.has_adult.toString()"
	/>
</template>
<script setup lang="ts">
import type { PublicReportWater } from "@/type/api";
import ReportDetailEntry from "@/rmo/components/ReportDetailEntry.vue";
interface Props {
	water: PublicReportWater;
}
const props = defineProps<Props>();
</script>
