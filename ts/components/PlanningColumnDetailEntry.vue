<template>
	<td>
		<button
			@click="toggleSignal(signal)"
			class="btn btn-sm btn-link text-danger p-0 ms-1"
			style="font-size: 0.7rem"
		>
			<i class="bi bi-x"></i>
		</button>
	</td>
	<td>
		<span v-if="signal.type === 'flyover pool'">Pool</span>
		<span v-else-if="signal.type === 'publicreport nuisance'"
			>Nuisance</span
		>
		<span v-else-if="signal.type === 'publicreport water'"
			>Water</span
		>
	</td>
	<td>
		<TimeRelative :time="signal.created"></TimeRelative>
	</td>
	<td>{{ shortAddress(signal.address) }}</td>
</template>

<script setup lang="ts">
import { shortAddress } from "../format";
import TimeRelative from "@/components/TimeRelative.vue";
interface Props {
	signal: Signal;
}
const props = defineProps<Props>();
</script>
