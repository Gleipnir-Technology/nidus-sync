<template>
	<div>
		<div class="text-muted">
			<i class="bi" :class="typeToIcon()" />{{ typeToTitle() }}
		</div>
		<div>{{ entry.message }}</div>
		<small class="text-muted">{{ formatDate(entry.created) }}</small>
	</div>
</template>

<script setup lang="ts">
import { LogEntry } from "@/type/api";

interface Props {
	entry: LogEntry;
}
const props = defineProps<Props>();
function formatDate(date: Date) {
	return date.toLocaleString();
}
function typeToIcon(): string {
	if (props.entry.type == "message-text-incoming") {
		return "bi-box-arrow-in-left";
	} else if (props.entry.type == "message-text-outgoing") {
		return "bi-box-arrow-out-right";
	} else if (props.entry.type == "created") {
		return "bi-stars";
	} else {
		return "bi-question";
	}
}
function typeToTitle(): string {
	if (props.entry.type == "message-text-incoming") {
		return "Incoming Text";
	} else if (props.entry.type == "message-text-outgoing") {
		return "Outgoing Text";
	} else if (props.entry.type == "created") {
		return "Report Created";
	} else {
		return props.entry.type;
	}
}
</script>
