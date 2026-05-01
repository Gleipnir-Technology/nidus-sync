<template>
	<!-- First row: icon, type badge, and time -->
	<div class="justify-content-between align-items-center">
		<div class="row">
			<div class="d-flex align-items-center">
				<div class="col">
					<Tooltip placement="top" :title="tooltipTitleForCommunicationType()">
						<i class="bi fs-4 me-2" :class="iconForReportType()"></i>
					</Tooltip>
					<Tooltip placement="top" :title="tooltipTitleForReportType()">
						<i class="bi fs-4 me-2" :class="iconForCommunicationType()"></i>
					</Tooltip>
				</div>
				<div class="col-6 text-end">
					<small>
						<Tooltip placement="top" :title="tooltipTitleForCreated()">
							<TimeRelative :time="comm.created" />
						</Tooltip>
					</small>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import TimeRelative from "@/components/TimeRelative.vue";
import Tooltip from "@/components/Tooltip.vue";
import { formatAddress, formatDate } from "@/format";
import { Communication } from "@/type/api";
interface Props {
	comm: Communication;
}
const props = defineProps<Props>();
function iconForCommunicationType(): string {
	switch (props.comm.type) {
		case "publicreport.compliance":
			return "bi-card-checklist";
		case "publicreport.nuisance":
			return "bi-mosquito";
		case "publicreport.water":
			return "bi-droplet-fill";
		default:
			return "";
	}
}
function iconForReportType(): string {
	switch (props.comm.type) {
		case "publicreport.compliance":
		case "publicreport.nuisance":
		case "publicreport.water":
			return "bi-postcard";
		case "email":
			return "bi-envelope";
		case "text":
			return "bi-chat-dots";
		default:
			return "";
	}
}
function tooltipTitleForCommunicationType(): string {
	switch (props.comm.type) {
		case "publicreport.compliance":
		case "publicreport.nuisance":
		case "publicreport.water":
			return "A report made from a member of the public to report.mosquitoes.online";
		case "email":
			return "An email received from a member of the public";
		case "text":
			return "An SMS/MMS text message received from a member of the public";
		default:
			return "I'm actually not sure what this is. How are you even seeing this?";
	}
}
function tooltipTitleForReportType(): string {
	switch (props.comm.type) {
		case "publicreport.compliance":
		case "publicreport.nuisance":
		case "publicreport.water":
			return "A compliance report either made by scanning a door hanger or by receiving a personal letter through the mail";
		case "publicreport.nuisance":
			return "A report of a mosquito nuisance";
		case "publicreport.water":
			return "A report of standing water";
		default:
			return "I'm actually not sure what this is. This shouldn't be possible.";
	}
}
function tooltipTitleForCreated(): string {
	return `or at exactly ${formatDate(props.comm.created)}`;
}
</script>
