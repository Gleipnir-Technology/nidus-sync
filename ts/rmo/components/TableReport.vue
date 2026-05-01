<style scoped>
.table {
	width: 100%;
	margin-bottom: 0;
	border-collapse: collapse;
}

.table-light {
	background-color: #f8f9fa;
}

.table-hover tbody tr:hover {
	background-color: rgba(0, 0, 0, 0.075);
}

th,
td {
	padding: 0.75rem;
	border-bottom: 1px solid #dee2e6;
	text-align: left;
}

.clickable-row {
	cursor: pointer;
	transition: background-color 0.15s ease-in-out;
}

.clickable-row:hover {
	background-color: rgba(13, 110, 253, 0.1);
}

.badge {
	display: inline-block;
	padding: 0.35em 0.65em;
	font-size: 0.75em;
	font-weight: 700;
	line-height: 1;
	color: #fff;
	text-align: center;
	white-space: nowrap;
	vertical-align: baseline;
	border-radius: 0.25rem;
}

.bg-danger {
	background-color: #dc3545;
}

.bg-primary {
	background-color: #0d6efd;
}

.bg-success {
	background-color: #198754;
}

.bg-warning {
	background-color: #ffc107;
}

.bg-info {
	background-color: #0dcaf0;
}

.bg-secondary {
	background-color: #6c757d;
}

.report-type-badge {
	font-size: 0.85rem;
}

.text-dark {
	color: #212529 !important;
}
</style>
<template>
	<table class="table table-hover mb-0">
		<thead class="table-light">
			<tr>
				<th scope="col">Report ID</th>
				<th scope="col">Reported</th>
				<th scope="col">Type</th>
				<th scope="col">Address</th>
				<th scope="col">Status</th>
			</tr>
		</thead>
		<tbody id="report-table-body">
			<tr
				v-if="reports.length > 0"
				v-for="report in reports"
				:key="report.id"
				class="clickable-row"
				:data-report-id="report.id"
				@click="handleRowClick(report.id)"
			>
				<td>
					<strong>{{ formatId(report.id) }}</strong>
				</td>
				<td><TimeRelative :time="report.created" /></td>
				<td>
					<span
						class="badge report-type-badge"
						:class="getTypeClass(report.type)"
					>
						{{ report.type }}
					</span>
				</td>
				<td>{{ report.address || "N/A" }}</td>
				<td>
					<span class="badge" :class="getStatusClass(report.status)">
						{{ report.status }}
					</span>
				</td>
			</tr>
			<tr v-else>
				<td colspan="5">No reports</td>
			</tr>
		</tbody>
	</table>
</template>

<script setup lang="ts">
import TimeRelative from "@/components/TimeRelative.vue";
export interface Report {
	id: string;
	created: Date;
	type: string;
	address?: string;
	status: string;
}

interface Props {
	reports?: Report[];
}

// Define props with defaults
const props = withDefaults(defineProps<Props>(), {
	reports: () => [],
});

// Define emits
const emit = defineEmits<{
	(e: "rowClicked", reportId: string): void;
}>();

/**
 * Get badge color class based on report type
 */
const getTypeClass = (type: string): string => {
	switch (type) {
		case "nuisance":
			return "bg-danger";
		case "quick":
			return "bg-primary";
		case "water":
			return "bg-success";
		default:
			return "bg-secondary";
	}
};

/**
 * Get badge color class based on report status
 */
const getStatusClass = (status: string): string => {
	switch (status) {
		case "Reported":
			return "bg-warning text-dark";
		case "Assigned":
			return "bg-info text-dark";
		case "On-Hold":
			return "bg-secondary";
		case "Complete":
			return "bg-success";
		default:
			return "bg-secondary";
	}
};

/**
 * Format the report ID with hyphens
 */
const formatId = (id: string): string => {
	if (id.length === 12) {
		return `${id.substring(0, 4)}-${id.substring(4, 8)}-${id.substring(8)}`;
	}
	return id;
};

/**
 * Handle row click event
 */
const handleRowClick = (report_id: string): void => {
	console.log("row clicked", report_id);
	emit("rowClicked", report_id);
};
</script>
