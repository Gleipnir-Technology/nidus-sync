<style scoped>
.scroll-pane {
	max-height: calc(100vh - 200px);
	overflow-y: auto;
}

.tool-button {
	width: 100%;
	margin-bottom: 0.5rem;
	text-align: left;
}
</style>
<template>
	<div class="card shadow-sm h-100">
		<div class="card-header bg-white pane-header">Transformation Tools</div>
		<div class="card-body scroll-pane">
			<div class="mb-3">
				<div class="text-muted small mb-2">Signal → Lead</div>
				<button
					class="btn btn-outline-primary tool-button"
					:disabled="selectedSignals.length === 0 || creating"
					@click="createLead()"
				>
					<span v-if="!creating">Create New Lead from Selection</span>
					<span v-else>
						<span class="spinner-border spinner-border-sm me-1"></span>
						Creating...
					</span>
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					:disabled="selectedSignals.length === 0"
				>
					Add Signals to Existing Lead
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					:disabled="selectedSignals.length === 0"
					@click="markAsAddressed()"
				>
					Mark Signal as Addressed
				</button>
			</div>

			<hr />

			<div class="mb-3">
				<div class="text-muted small mb-2">Lead → Field Assignment</div>
				<button class="btn btn-outline-success tool-button">
					Create Proposed Assignment
				</button>
				<button class="btn btn-outline-secondary tool-button">
					Add Leads to Existing Assignment
				</button>
				<button class="btn btn-outline-secondary tool-button">
					Split Lead
				</button>
			</div>

			<hr />

			<div class="mb-3">
				<div class="text-muted small mb-2">Assignment Controls</div>
				<button class="btn btn-outline-dark tool-button">Set Priority</button>
				<button class="btn btn-outline-dark tool-button">
					Estimate Effort
				</button>
				<button class="btn btn-outline-dark tool-button">
					Send to Operations
				</button>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
interface Props {
	creating: boolean;
	selectedSignals: int[];
}
const props = defineProps<Props>();
</script>
