<style scoped>
.error-message {
	background-color: #f8d7da;
	border: 1px solid #f5c2c7;
	border-radius: 0.25rem;
	padding: 1rem;
	margin-bottom: 1rem;
	color: #842029;
}

.filter-label {
	font-size: 0.875rem;
	font-weight: 500;
}

.loading-spinner {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 2rem;
}

.signal-address {
	font-size: 0.875rem;
	color: #6c757d;
}

.signal-item {
	cursor: pointer;
	transition: all 0.2s;
}

.signal-item:hover {
	background-color: #f8f9fa;
}

.signal-item.selected {
	background-color: #e7f3ff;
	border-color: #0d6efd;
}
</style>
<template>
	<div class="card shadow-sm h-100">
		<div class="card-header bg-white pane-header">
			Incoming Signals & Leads
			<span
				v-show="loading"
				class="spinner-border spinner-border-sm ms-2"
				role="status"
			></span>
		</div>
		<div class="card-body scroll-pane">
			<!-- Error Display -->
			<div v-if="error" class="error-message">
				<strong>Error:</strong> <span>{{ error }}</span>
				<button
					@click="emit('refresh')"
					class="btn btn-sm btn-outline-danger mt-2 w-100"
				>
					Retry
				</button>
			</div>

			<!-- FILTERS -->
			<div class="mb-3">
				<div class="filter-label mb-1">Species</div>
				<select
					class="form-select form-select-sm mb-2 disabled"
					disabled
					v-model="filters.species"
					@change="loadSignals()"
				>
					<option value="">All Species</option>
					<option value="aedes_aegypti">Aedes aegypti</option>
					<option value="aedes_albopictus">Aedes albopictus</option>
					<option value="culex_pipiens">Culex pipiens</option>
					<option value="culex_tarsalis">Culex tarsalis</option>
				</select>

				<div class="filter-label mb-1">Signal Type</div>
				<select
					class="form-select form-select-sm mb-2 disabled"
					disabled
					v-model="filters.type"
					@change="loadSignals()"
				>
					<option value="">All Types</option>
					<option value="public_report">Public Report</option>
					<option value="trap_spike">Trap Spike</option>
					<option value="surveillance">Surveillance Observation</option>
					<option value="residual_expiring">Residual Expiring</option>
					<option value="plan_followup">Plan Follow-Up</option>
				</select>

				<div class="filter-label mb-1">Sort By</div>
				<select
					class="form-select form-select-sm disabled"
					disabled
					v-model="filters.sort"
					@change="loadSignals()"
				>
					<option value="newest">Newest First</option>
					<option value="priority">Highest Priority</option>
					<option value="linked">Most Signals Linked</option>
					<option value="species_signal">Strongest Species Signal</option>
				</select>
			</div>

			<hr />

			<!-- Loading State -->
			<div v-if="loading" class="loading-spinner">
				<div class="spinner-border text-primary" role="status">
					<span class="visually-hidden">Loading...</span>
				</div>
			</div>
			<div v-else>
				<div class="mb-3">
					<div class="fw-semibold mb-2">
						Signals
						<span class="badge bg-primary" v-show="selectedSignalIDs.size > 0">
							{{ selectedSignalIDs.size }}
						</span>
					</div>

					<div
						v-if="signals != null && signals.length == 0"
						class="text-muted small fst-italic"
					>
						No signals found
					</div>

					<div
						v-for="signal in signals"
						v-if="signals != null && signals.length > 0"
						:key="signal.id"
						class="border rounded p-2 mb-2 signal-item"
						:class="{ selected: isSelected(signal.id) }"
						@click="toggleSignal(signal)"
					>
						<PlanningColumnListEntry :selected="selectedSignalIDs.has(signal.id)" :signal="signal"/>
					</div>
				</div>
			</div>

			<hr />

			<!-- Mosquito Control Plan Followups -->
			<div class="mb-3" v-show="!loading || planFollowups.length > 0">
				<div class="fw-semibold mb-2">Mosquito Control Plan Follow-Ups</div>

				<div
					v-if="planFollowups.length === 0 && !loading"
					class="text-muted small fst-italic"
				>
					No plan follow-ups
				</div>

				<div
					v-for="followup in planFollowups"
					:key="followup.id"
					class="border rounded p-2 mb-2 signal-item"
					:class="{ selected: isSelected(followup.id) }"
					@click="toggleSignal(followup)"
				>
					<div class="small fw-semibold">{{ followup.title }}</div>
					<div class="text-muted small">{{ followup.description }}</div>
					<span class="badge bg-secondary">Plan</span>
				</div>
			</div>

			<hr />

			<!-- Leads -->
			<div v-show="!loading || leads.length > 0">
				<div class="fw-semibold mb-2">Existing Leads</div>

				<div
					v-if="leads.length === 0 && !loading"
					class="text-muted small fst-italic"
				>
					No existing leads
				</div>

				<div
					v-for="lead in leads"
					:key="lead.id"
					class="border rounded p-2 mb-2 signal-item"
				>
					<div class="small fw-semibold">{{ lead.title }}</div>
					<div class="text-muted small">{{ lead.description }}</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import PlanningColumnListEntry from "@/components/PlanningColumnListEntry.vue";

interface Emits {
	(e: "refresh"): void;
	(e: "signalDeselect", integer): void;
	(e: "signalSelect", integer): void;
}
interface Props {
	error: string | null;
	leads: Lead[] | null;
	loading: boolean;
	planFollowups: Followup[] | null;
	selectedSignalIDs: Set<int>;
	signals: Signal[] | null;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
const filters = ref({
	species: "",
	type: "",
	sort: "newest",
});
const isSelected = (id) => {
	return props.selectedSignalIDs.values().some((s) => s.id === id);
};
const toggleSignal = (signal: int) => {
	if (props.selectedSignalIDs.has(signal)) {
		emit("signalDeselect", signal);
	} else {
		emit("signalSelect", signal);
	}
};
</script>
