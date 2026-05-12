<style scoped>
.filters-section {
	font-size: 0.875rem;
}

.btn-group-sm .btn {
	font-size: 0.75rem;
	padding: 0.25rem 0.5rem;
}

.form-label.small {
	margin-bottom: 0.25rem;
	color: #6c757d;
}
</style>
<template>
	<div class="card shadow-sm h-100 reports-list">
		<div class="card-header bg-light pane-header">
			<!-- Filter Section -->
			<div class="filters-section">
				<!-- Status Filter -->
				<div class="mb-2">
					<label class="form-label small fw-bold mb-1">Status</label>
					<div class="btn-group btn-group-sm d-flex" role="group">
						<input
							type="checkbox"
							class="btn-check"
							id="status-new"
							v-model="statusFilters.new"
							autocomplete="off"
						/>
						<label class="btn btn-outline-success" for="status-new">New</label>

						<input
							type="checkbox"
							class="btn-check"
							id="status-opened"
							v-model="statusFilters.opened"
							autocomplete="off"
						/>
						<label class="btn btn-outline-primary" for="status-opened"
							>Opened</label
						>

						<input
							type="checkbox"
							class="btn-check"
							id="status-pending"
							v-model="statusFilters.pending"
							autocomplete="off"
						/>
						<label class="btn btn-outline-warning" for="status-pending"
							>Pending</label
						>

						<input
							type="checkbox"
							class="btn-check"
							id="status-closed"
							v-model="statusFilters.closed"
							autocomplete="off"
						/>
						<label class="btn btn-outline-secondary" for="status-closed"
							>Closed</label
						>
					</div>
				</div>

				<!-- Source Filter -->
				<div class="mb-2">
					<label class="form-label small fw-bold mb-1">Source</label>
					<div class="btn-group btn-group-sm d-flex" role="group">
						<button
							class="btn"
							:class="
								sourceFilter === 'all' ? 'btn-primary' : 'btn-outline-secondary'
							"
							@click="sourceFilter = 'all'"
						>
							All
						</button>
						<button
							class="btn"
							:class="
								sourceFilter === 'email'
									? 'btn-primary'
									: 'btn-outline-secondary'
							"
							@click="sourceFilter = 'email'"
						>
							<i class="bi bi-envelope"></i> Email
						</button>
						<button
							class="btn"
							:class="
								sourceFilter === 'text'
									? 'btn-primary'
									: 'btn-outline-secondary'
							"
							@click="sourceFilter = 'text'"
						>
							<i class="bi bi-chat-dots"></i> Text
						</button>
						<button
							class="btn"
							:class="
								sourceFilter === 'publicreport'
									? 'btn-primary'
									: 'btn-outline-secondary'
							"
							@click="sourceFilter = 'publicreport'"
						>
							<i class="bi bi-file-text"></i> Form
						</button>
					</div>
				</div>

				<!-- Type Filter -->
				<div class="mb-2" v-if="sourceFilter == 'publicreport'">
					<label class="form-label small fw-bold mb-1">Type</label>
					<div class="btn-group btn-group-sm d-flex" role="group">
						<button
							class="btn"
							:class="
								typeFilter === 'all' ? 'btn-primary' : 'btn-outline-secondary'
							"
							@click="typeFilter = 'all'"
						>
							All
						</button>
						<button
							class="btn"
							:class="
								typeFilter === 'publicreport.compliance'
									? 'btn-success'
									: 'btn-outline-secondary'
							"
							@click="typeFilter = 'publicreport.compliance'"
						>
							<i class="bi bi-check-circle"></i> Compliance
						</button>
						<button
							class="btn"
							:class="
								typeFilter === 'publicreport.nuisance'
									? 'btn-danger'
									: 'btn-outline-secondary'
							"
							@click="typeFilter = 'publicreport.nuisance'"
						>
							<i class="bi bi-mosquito"></i> Nuisance
						</button>
						<button
							class="btn"
							:class="
								typeFilter === 'publicreport.water'
									? 'btn-info'
									: 'btn-outline-secondary'
							"
							@click="typeFilter = 'publicreport.water'"
						>
							<i class="bi bi-droplet"></i> Water
						</button>
					</div>
				</div>

				<!-- Active Filters Summary -->
				<div class="small text-muted" v-if="activeFilterCount > 0">
					<i class="bi bi-funnel"></i> {{ filteredCommunications.length }} of
					{{ props.all?.length || 0 }} reports
				</div>
			</div>
		</div>

		<div class="card-body scroll-pane">
			<div class="list-group list-group-flush">
				<div class="loading list-group-item" v-if="loading || all == null">
					Loading...
				</div>
				<div
					v-else-if="filteredCommunications.length > 0"
					v-for="comm in filteredCommunications"
					:key="comm.id"
				>
					<ListCardCommunication
						@click="handleClick(comm.id)"
						:comm="comm"
						:isSelected="selectedID == comm.id"
					/>
				</div>
				<div v-else class="list-group-item text-center text-muted p-4">
					<i class="bi bi-inbox fs-1"></i>
					<p class="mt-2">No reports match the current filters</p>
					<button class="btn btn-sm btn-outline-primary" @click="resetFilters">
						Reset Filters
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import ListCardCommunication from "@/components/ListCardCommunication.vue";
import { Communication, LogEntry, PublicReport } from "@/type/api";

interface Props {
	all: Communication[] | null;
	loading: boolean;
	selectedID?: string;
}
interface Emits {
	(e: "deselect", id: string): void;
	(e: "select", id: string): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const handleClick = (id: string) => {
	if (props.selectedID == undefined) {
		emit("select", id);
	} else if (props.selectedID == id) {
		emit("deselect", id);
	} else {
		emit("select", id);
	}
};

// Filter state
const searchFilter = ref<string>("");
const typeFilter = ref<string>("all");
const sourceFilter = ref<string>("all");
const statusFilters = ref({
	new: true,
	opened: true,
	pending: false,
	closed: false,
});

// Reset filters to default
const resetFilters = () => {
	searchFilter.value = "";
	typeFilter.value = "all";
	sourceFilter.value = "all";
	statusFilters.value = {
		new: true,
		opened: true,
		pending: false,
		closed: false,
	};
};

// Count active filters
const activeFilterCount = computed(() => {
	let count = 0;
	if (typeFilter.value !== "all") count++;
	if (sourceFilter.value !== "all") count++;
	if (searchFilter.value.length > 0) count++;
	return count;
});

// Filtered communications
const filteredCommunications = computed((): Communication[] => {
	if (props.all == null) {
		return [];
	}

	return props.all.filter((comm) => {
		// Status filter
		const selectedStatuses = Object.entries(statusFilters.value)
			.filter(([_, isSelected]) => isSelected)
			.map(([status]) => status);

		if (
			selectedStatuses.length > 0 &&
			!selectedStatuses.includes(comm.status)
		) {
			return false;
		}

		// Source filter
		if (
			sourceFilter.value !== "all" &&
			!comm.type.startsWith(sourceFilter.value)
		) {
			return false;
		}

		// Type filter
		if (
			sourceFilter.value === "publicreport" &&
			typeFilter.value !== "all" &&
			comm.type !== typeFilter.value
		) {
			return false;
		}

		return true;
	});
});
</script>
