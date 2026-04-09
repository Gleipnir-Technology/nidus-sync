<template>
	<TimeRelative :time="signal.created"></TimeRelative>
	<p>{{ shortAddress(signal.address) }}</p>
	<div v-if="signal.type == 'flyover pool' && signal.pool">
		<FlyoverPoolCard :location="signal.location" :markers="[]" />
	</div>
	<div v-else-if="signal.type == 'publicreport nuisance' && signal.report">
		<PublicreportCard :report="signal.report" />
	</div>
	<div v-else-if="signal.type == 'publicreport water' && signal.report">
		<PublicreportCard :report="signal.report" />
	</div>
	<div v-else>
		<p>No report or pool</p>
	</div>
</template>

<script setup lang="ts">
import FlyoverPoolCard from "@/components/FlyoverPoolCard.vue";
import PublicreportCard from "@/components/PublicreportCard.vue";
import TimeRelative from "@/components/TimeRelative.vue";
import { shortAddress } from "@/format";
import { Signal } from "@/type/api";
interface Props {
	signal: Signal;
}
const props = defineProps<Props>();
</script>
