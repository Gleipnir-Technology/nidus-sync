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
		<div class="card-header bg-light pane-header">Transformation Tools</div>
		<div class="card-body scroll-pane">
			<div class="mb-3">
				<div class="text-muted small mb-2">Signal → Lead</div>
				<button
					class="btn btn-outline-primary tool-button"
					:disabled="selectedSignalIDs.size === 0 || creating"
					@click="emit('doCreateLead')"
				>
					<span v-if="!creating">Create New Lead from Selection</span>
					<span v-else>
						<span class="spinner-border spinner-border-sm me-1"></span>
						Creating...
					</span>
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					:disabled="selectedSignalIDs.size === 0"
					@click="emit('doAddToLead')"
				>
					Add Signals to Existing Lead
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					:disabled="selectedSignalIDs.size === 0"
					@click="emit('doMarkSignalAddressed')"
				>
					Mark Signal as Addressed
				</button>
			</div>

			<hr />

			<div class="mb-3">
				<div class="text-muted small mb-2">Lead → Field Assignment</div>
				<button
					class="btn btn-outline-success tool-button"
					@click="emit('doCreateProposedAssignment')"
				>
					Create Proposed Assignment
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					@click="emit('doAddLeadsToAssignment')"
				>
					Add Leads to Existing Assignment
				</button>
				<button
					class="btn btn-outline-secondary tool-button"
					@click="emit('doSplitLead')"
				>
					Split Lead
				</button>
			</div>

			<hr />

			<div class="mb-3">
				<div class="text-muted small mb-2">Assignment Controls</div>
				<button
					class="btn btn-outline-dark tool-button"
					@click="emit('doSetPriority')"
				>
					Set Priority
				</button>
				<button
					class="btn btn-outline-dark tool-button"
					@click="emit('doEstimateEffort')"
				>
					Estimate Effort
				</button>
				<button
					class="btn btn-outline-dark tool-button"
					@click="emit('doSendToOperations')"
				>
					Send to Operations
				</button>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
interface Emits {
	(e: "doAddToLead"): void;
	(e: "doAddLeadsToAssignment"): void;
	(e: "doCreateLead"): void;
	(e: "doCreateProposedAssignment"): void;
	(e: "doEstimateEffort"): void;
	(e: "doMarkSignalAddressed"): void;
	(e: "doSetPriority"): void;
	(e: "doSendToOperations"): void;
	(e: "doSplitLead"): void;
}
interface Props {
	creating: boolean;
	selectedSignalIDs: Set<number>;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
</script>
