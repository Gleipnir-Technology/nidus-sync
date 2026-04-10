<template>
	<div class="card-header bg-info bg-opacity-10">
		<i class="bi bi-droplet"></i> Standing Water Details
	</div>
	<div class="card-body">
		<div
			v-if="
				report.access_gate ||
				report.access_fence ||
				report.access_locked ||
				report.access_dog ||
				report.access_other
			"
			class="col-md-6"
		>
			<label class="form-label text-muted small mb-0">
				<i class="bi bi-droplet"></i> Access
			</label>
			<div>
				<ul>
					<li v-if="report.access_gate">Gate</li>
					<li v-if="report.access_fence">Fence</li>
					<li v-if="report.access_locked">Locked</li>
					<li v-if="report.access_dog">Dog</li>
					<li v-if="report.access_other">Other access obstacle</li>
				</ul>
			</div>
		</div>
		<div v-if="report.access_comments" class="col-12">
			<label class="form-label text-muted small mb-0">
				<i class="bi bi-chat-text"></i> Access Comments
			</label>
			<div class="p-2 bg-light rounded">
				{{ report.access_comments }}
			</div>
		</div>
		<label class="form-label text-muted small mb-0">
			<i class="bi bi-eye"></i> Mosquito Life Stages Observed
		</label>
		<div class="mt-2">
			<span
				class="badge me-2"
				:class="report.has_larvae ? 'badge-larvae' : 'bg-light text-muted'"
			>
				<i
					class="bi"
					:class="report.has_larvae ? 'bi-check-circle' : 'bi-circle'"
				></i>
				Larvae
			</span>
			<span
				class="badge me-2"
				:class="report.has_pupae ? 'badge-pupae' : 'bg-light text-muted'"
			>
				<i
					class="bi"
					:class="report.has_pupae ? 'bi-check-circle' : 'bi-circle'"
				></i>
				Pupae
			</span>
			<span
				class="badge"
				:class="report.has_adult ? 'badge-adult' : 'bg-light text-muted'"
			>
				<i
					class="bi"
					:class="report.has_adult ? 'bi-check-circle' : 'bi-circle'"
				></i>
				Adult Mosquitoes
			</span>
		</div>
		<div v-if="report.comments" class="col-12">
			<label class="form-label text-muted small mb-0">
				<i class="bi bi-chat-text"></i> Comments
			</label>
			<div class="p-2 bg-light rounded">
				{{ report.comments }}
			</div>
		</div>
		<div class="col-md-6">
			<label class="form-label text-muted small mb-0">
				<i class="bi bi-person"></i> Owner Name
			</label>
			<div class="fw-medium">
				{{ report.owner.name || "not given" }}
			</div>
		</div>
		<div class="col-md-6">
			<label
				v-if="report.owner.has_email"
				class="form-label text-muted small mb-0"
			>
				<i class="bi bi-envelope"></i>
			</label>
			<label
				v-if="report.owner.has_phone"
				class="form-label text-muted small mb-0"
			>
				<i class="bi bi-phone"></i>
			</label>
		</div>
	</div>
</template>
<script setup lang="ts">
import { PublicReportWater } from "@/type/api";
interface Props {
	report: PublicReportWater;
}
const props = defineProps<Props>();
</script>
