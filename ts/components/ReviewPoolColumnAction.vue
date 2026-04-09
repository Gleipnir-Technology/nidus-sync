<template>
	<h5 class="mb-4">Actions</h5>

	<button
		class="btn btn-success action-btn"
		@click="markReviewed"
		:disabled="!selectedTask || submitting"
	>
		<span v-if="!submitting">
			<i class="bi bi-check-circle"></i> Complete Review
		</span>
		<span v-else>
			<span class="spinner-border spinner-border-sm" role="status"></span>
			Submitting...
		</span>
	</button>

	<div v-if="selectedTask" class="card mt-3">
		<div v-if="changes.updated.length > 0" class="card-body">
			<h6 class="card-title">
				<i class="bi bi-pencil-square text-warning"></i> Updates
			</h6>
			<ul class="mb-0">
				<li v-for="item in changes.updated" :key="item">{{ item }}</li>
			</ul>
		</div>
		<div v-if="changes.unchanged.length > 0" class="card-body">
			<h6 class="card-title">
				<i class="bi bi-dash-circle text-muted"></i> Not changed
			</h6>
			<ul class="mb-0">
				<li v-for="item in changes.unchanged" :key="item">{{ item }}</li>
			</ul>
		</div>
	</div>

	<button
		class="btn btn-danger action-btn"
		@click="discardEntry"
		:disabled="!selectedTask || submitting"
	>
		<i class="bi bi-trash"></i> Discard Entry
	</button>

	<!-- Keyboard Shortcuts Help -->
	<div class="card mt-3">
		<div class="card-body">
			<h6 class="card-title"><i class="bi bi-keyboard"></i> Tips</h6>
			<small class="text-muted">
				Fields with a yellow border have been modified from their original
				values.
			</small>
		</div>
	</div>
</template>
<script setup lang="ts">
import MapMultipoint from "@/components/MapMultipoint.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { Changes } from "@/types";
import { ReviewTask } from "@/type/api";

interface Props {
	changes: Changes;
	selectedTask?: ReviewTask;
	submitting: boolean;
}
const props = defineProps<Props>();
function discardEntry() {
	console.log("Fake discard entry");
}
function markReviewed() {
	console.log("Fake mark reviewed");
}
</script>
