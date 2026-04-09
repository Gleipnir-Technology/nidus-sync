<style scoped>
.report-card {
	cursor: pointer;
	transition: background-color 0.2s;
}

.report-card:hover {
	background-color: $secondary;
}

.report-card.active {
	background-color: $primary;
	color: white;
}
.reports-list {
	overflow-y: auto;
	max-height: 100vh;
}
</style>

<template>
	<div class="card shadow-sm h-100 reports-list">
		<div class="card-header bg-light pane-header">
			<div class="input-group input-group-sm">
				<span class="input-group-text"><i class="bi bi-search"></i></span>
				<input
					type="text"
					class="form-control"
					placeholder="Filter reports..."
					v-model="searchFilter"
				/>
			</div>
		</div>
		<div class="card-body scroll-pane">
			<div class="mb-3">
				<button
					class="btn btn-sm"
					:class="
						typeFilter === 'all' ? 'btn-primary' : 'btn-outline-secondary'
					"
					@click="typeFilter = 'all'"
				>
					All
				</button>
				<button
					class="btn btn-sm"
					:class="
						typeFilter === 'nuisance' ? 'btn-danger' : 'btn-outline-secondary'
					"
					@click="typeFilter = 'nuisance'"
				>
					<i class="bi bi-mosquito"></i>Nuisance
				</button>
				<button
					class="btn btn-sm"
					:class="typeFilter === 'water' ? 'btn-info' : 'btn-outline-secondary'"
					@click="typeFilter = 'water'"
				>
					<i class="bi bi-droplet"></i> Water
				</button>
			</div>

			<div class="list-group list-group-flush">
				<div v-if="loading || all == null" class="loading">Loading...</div>
				<div
					v-else-if="all.length > 0"
					v-for="comm in filteredCommunications"
					:key="comm.id"
					class="border rounded list-group-item report-card p-3"
					:class="{
						active: selectedId && selectedId === comm.id,
					}"
					@click="handleClick(comm.id)"
				>
					<!-- First row: icon, type badge, and time -->
					<div class="d-flex justify-content-between align-items-center mb-2">
						<div class="d-flex align-items-center">
							<i
								v-if="comm.type === 'publicreport.nuisance'"
								class="bi bi-mosquito icon-nuisance fs-4 me-2"
							>
							</i>
							<i
								v-if="comm.type === 'publicreport.water'"
								class="bi bi-droplet-fill icon-standing-water fs-4 me-2"
							></i>
							<span
								class="badge"
								:class="
									comm.type === 'publicreport.nuisance'
										? 'bg-danger'
										: 'bg-info'
								"
							>
								{{
									comm.type === "publicreport.nuisance"
										? "Nuisance"
										: "Standing Water"
								}}
							</span>
						</div>
						<small>
							<TimeRelative :time="comm.created" />
						</small>
					</div>

					<!-- Details section: full width -->
					<div>
						<div>
							<i class="bi bi-geo-alt text-muted"></i>
							<span class="fw-medium">{{
								comm.public_report?.address.postal_code
							}}</span>
						</div>
						<small>{{ formatAddress(comm.public_report?.address) }}</small>
						<div
							v-if="
								comm.public_report?.images &&
								comm.public_report?.images.length > 0
							"
							class="mt-1"
						>
							<small class="text-muted">
								<i class="bi bi-camera"></i>
								{{ comm.public_report.images.length }} photo(s)
							</small>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div
			v-if="filteredCommunications.length === 0"
			class="text-center text-muted p-4"
		>
			<i class="bi bi-inbox fs-1"></i>
			<p class="mt-2">No reports found</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import TimeRelative from "@/components/TimeRelative.vue";
import { formatAddress } from "@/format";
import { Communication, LogEntry, PublicReport } from "@/type/api";

interface Props {
	all: Communication[] | null;
	loading: boolean;
	selectedId?: string | null;
}
interface Emits {
	(e: "deselect", id: string): void;
	(e: "select", id: string): void;
}

const props = withDefaults(defineProps<Props>(), {
	selectedId: null,
});

const emit = defineEmits<Emits>();
const handleClick = (id: string) => {
	if (props.selectedId == null) {
		emit("select", id);
	} else if (props.selectedId == id) {
		emit("deselect", id);
	} else {
		emit("select", id);
	}
};
const searchFilter = ref<string>("");
const typeFilter = ref<string>("all");

function selectCommunication(communication: Communication) {
	// Emit both events - one for general use, one for v-model
	console.log("selected", communication);
	emit("select", communication.id);
	//emit("update:selectedItem", communication);
	//messageText.value = "";
	//updateMap();
}

// Computed properties
const filteredCommunications = computed(() => {
	if (props.all == null) {
		return [];
	}
	return props.all.filter((c) => {
		const matchesType =
			typeFilter.value === "all" || c.type === typeFilter.value;
		return matchesType && filterMatches(searchFilter.value, c);
	});
});
// Methods
function filterMatches(filter: string, comm: Communication) {
	const pr = comm.public_report;
	// When we have non-public-report communications fix this.
	if (pr == null) {
		return false;
	}
	return filterMatchesPublicReport(filter, pr);
}
function filterMatchesLogEntry(filter: string, logs: LogEntry[]) {
	for (const le of logs) {
		if (le.message.includes(filter)) {
			return true;
		}
	}
}
function filterMatchesPublicReport(filter: string, pr: PublicReport) {
	if (
		pr.address.raw.includes(filter) ||
		pr.id.includes(filter) ||
		filterMatchesLogEntry(filter, pr.log)
	) {
		return true;
	}
	return false;
}
</script>
