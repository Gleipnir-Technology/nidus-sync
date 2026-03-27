<template>
	<TimeRelative :time="signal.created"></TimeRelative>
	<p>{{ shortAddress(signal.address) }}</p>
	<div v-if="signal.type == 'flyover pool' && signal.pool">
		<FlyoverPoolCard :pool="signal.pool" />
	</div>
	<div v-else-if="signal.type == 'publicreport nuisance'">
		<PublicreportCard :report="signal.report" />
	</div>
	<div v-else-if="signal.type == 'publicreport water'">
		<PublicreportCard :report="signal.report" />
	</div>
</template>

<script setup lang="ts">
import { shortAddress } from "@/format";
import FlyoverPoolCard from "@/components/FlyoverPoolCard.vue";
import PublicreportCard from "@/components/PublicreportCard.vue";
import TimeRelative from "@/components/TimeRelative.vue";
interface Props {
	signal: Signal;
}
const props = defineProps<Props>();
</script>
