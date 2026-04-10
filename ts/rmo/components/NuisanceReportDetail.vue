<template>
	<ReportDetailEntry
		label="Active early morning (5a-8a)?"
		:value="nuisance.time_of_day_early.toString()"
	/>
	<ReportDetailEntry
		label="Active daytime (5a-8a)?"
		:value="nuisance.time_of_day_day.toString()"
	/>
	<ReportDetailEntry
		label="Active evening (5p-9p)?"
		:value="nuisance.time_of_day_evening.toString()"
	/>
	<ReportDetailEntry
		label="Active night (9p-5a)?"
		:value="nuisance.time_of_day_night.toString()"
	/>
	<ReportDetailEntry label="Duration" :value="nuisance.duration" />
	<ReportDetailEntry
		label="Active in backyard?"
		:value="nuisance.is_location_backyard.toString()"
	/>
	<ReportDetailEntry
		label="Active in frontyard?"
		:value="nuisance.is_location_frontyard.toString()"
	/>
	<ReportDetailEntry
		label="Active in garden?"
		:value="nuisance.is_location_garden.toString()"
	/>
	<ReportDetailEntry
		label="Active in other?"
		:value="nuisance.is_location_other.toString()"
	/>
	<ReportDetailEntry
		label="Active in pool?"
		:value="nuisance.is_location_pool.toString()"
	/>
	<ReportDetailEntry
		label="Stagnant water?"
		:value="nuisance.source_stagnant.toString()"
	/>
	<ReportDetailEntry
		label="container?"
		:value="nuisance.source_container.toString()"
	/>
	<ReportDetailEntry
		label="Sprinklers and/or gutters?"
		:value="nuisance.source_gutter.toString()"
	/>
</template>
<script setup lang="ts">
import type { PublicReportNuisance } from "@/type/api";
import ReportDetailEntry from "@/rmo/components/ReportDetailEntry.vue";
interface Props {
	nuisance: PublicReportNuisance;
}
const props = defineProps<Props>();
</script>
