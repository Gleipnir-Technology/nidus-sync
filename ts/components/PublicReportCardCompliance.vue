<template>
	<div class="card-header bg-danger bg-opacity-10">
		<i class="bi bi-exclamation-triangle"></i> Compliance Details
	</div>
	<div class="card-body">
		<div class="row g-3">
			<div class="col-md-12">
				<label
					class="form-label small mb-0"
					v-if="report.permission_type == 'granted'"
				>
					<i class="bi bi-check"></i> Access granted for tech to enter even if
					not at home.
				</label>
				<label
					class="form-label small mb-0"
					v-else-if="report.permission_type == 'with-owner'"
				>
					<i class="bi bi-info-square"></i> Access granted for tech, but owner
					must be present
				</label>
				<label
					class="form-label small mb-0"
					v-else-if="report.permission_type == 'denied'"
				>
					<i class="bi bi-exclamation-triangle"></i> Access denied
				</label>
				<label
					class="form-label small mb-0"
					v-else-if="report.permission_type == 'unselected'"
				>
					<i class="bi bi-question"></i> No access information provided
				</label>
			</div>
			<div class="col-md-6" v-if="report.wants_scheduled">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-clock"></i> Owner requests a scheduled visit
				</label>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-key"></i> Access Instructions
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.access_instructions || "none" }}
				</div>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-clock"></i> Availability Notes
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.availability_notes || "none" }}
				</div>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-pencil-square"></i> Comments
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.comments || "none" }}
				</div>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0">
					<i class="bi bi-key"></i> Gate Code
				</label>
				<div class="p-2 bg-light rounded">
					{{ report.gate_code || "none" }}
				</div>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0" v-if="report.has_dog">
					<i class="bi bi-dog"></i> Has Dog
				</label>
			</div>
			<div class="col-md-6">
				<label class="form-label text-muted small mb-0" v-if="report.has_dog">
					<i class="bi bi-dog"></i> Has Dog
				</label>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { PublicReportCompliance } from "@/type/api";
interface Props {
	report: PublicReportCompliance;
}
const props = defineProps<Props>();
</script>
