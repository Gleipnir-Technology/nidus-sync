<template>
	<!-- Error Alert -->
	<div v-if="error" class="mt-3 alert alert-danger alert-dismissible">
		{{ error }}
		<button type="button" class="btn-close" @click="error = null"></button>
	</div>

	<div class="p-3 border-bottom bg-primary text-white">
		<h5 class="mb-0"><i class="bi bi-list-ul"></i> Review Queue</h5>
		<small>{{ total }} entries pending</small>
	</div>

	<!-- Loading State -->
	<div v-if="tasks == null" class="p-4 text-center">
		<span class="spinner-border" role="status"></span>
		<p class="mt-2">Loading tasks...</p>
	</div>

	<!-- Empty State -->
	<div v-else-if="tasks.length === 0" class="p-4 text-center text-muted">
		<i class="bi bi-check-circle" style="font-size: 48px"></i>
		<p class="mt-2">No entries to review!</p>
	</div>

	<!-- Task List -->
	<div
		v-for="task in tasks"
		:key="task.id"
		class="entry-item"
		:class="{ active: selectedTaskID === task.id }"
		@click="selectTask(task)"
	>
		<div class="d-flex justify-content-between align-items-start">
			<div>
				<i class="bi bi-droplet"></i>
				<strong>Pool {{ task.id }}</strong>
			</div>
			<small class="text-muted">{{ task.condition }}</small>
		</div>
		<small class="text-muted d-block mt-1">
			{{ formatAddress(task.address) }}
		</small>
	</div>
</template>
<script setup lang="ts">
import { formatAddress } from "@/format";

interface Props {
	error: string | null;
	selectedTaskID: int | null;
	tasks: ReviewTask[];
	total: int;
}
const props = defineProps<Props>();
</script>
