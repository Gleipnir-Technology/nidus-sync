<template>
	<!-- First row: icon, type badge, and time -->
	<div class="justify-content-between align-items-center">
		<div class="row">
			<div class="d-flex align-items-center">
				<div class="col">
					<i class="bi fs-4 me-2" :class="iconForType()"></i>
				</div>
				<div class="col-6 text-end">
					<span class="badge" :class="colorForType()">
						{{ titleForType() }}
					</span>
				</div>
			</div>
		</div>
		<div class="row">
			<small>
				<TimeRelative :time="comm.created" />
			</small>
		</div>
	</div>
</template>

<script setup lang="ts">
import TimeRelative from "@/components/TimeRelative.vue";
import { formatAddress } from "@/format";
import { Communication } from "@/type/api";
interface Props {
	comm: Communication;
}
const props = defineProps<Props>();
function colorForType(): string {
	if (props.comm.type == "publicreport.compliance") {
		return "bg-secondary";
	} else if (props.comm.type == "publicreport.nuisance") {
		return "bg-danger";
	} else if (props.comm.type == "publicreport.water") {
		return "bg-info";
	} else {
		return "";
	}
}
function iconForType(): string {
	if (props.comm.type == "publicreport.compliance") {
		return "bi-postcard";
	} else if (props.comm.type == "publicreport.nuisance") {
		return "bi-mosquito";
	} else if (props.comm.type == "publicreport.water") {
		return "bi-droplet-fill";
	} else {
		return "";
	}
}
function titleForType(): string {
	if (props.comm.type == "publicreport.compliance") {
		return "Compliance";
	} else if (props.comm.type == "publicreport.nuisance") {
		return "Nuisance";
	} else if (props.comm.type == "publicreport.water") {
		return "Standing Water";
	} else {
		return "Unknown";
	}
}
</script>
