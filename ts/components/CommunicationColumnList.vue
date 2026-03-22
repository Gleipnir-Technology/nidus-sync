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
</style>

<template>
	<div class="col-md-3 border-end p-0 reports-list">
		<div class="p-3 bg-light border-bottom">
			<div class="input-group input-group-sm">
				<span class="input-group-text"><i class="bi bi-search"></i></span>
				<input
					type="text"
					class="form-control"
					placeholder="Filter reports..."
					v-model="searchFilter"
				/>
			</div>
			<div class="mt-2 d-flex gap-2">
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
					<i class="bi bi-mosquito"></i>Mosquito Nuisance
				</button>
				<button
					class="btn btn-sm"
					:class="typeFilter === 'water' ? 'btn-info' : 'btn-outline-secondary'"
					@click="typeFilter = 'water'"
				>
					<i class="bi bi-droplet"></i> Water
				</button>
			</div>
		</div>

		<div class="list-group list-group-flush">
			<div v-if="loading" class="loading">Loading...</div>
			<div
				v-else-if="all.length > 0"
				v-for="comm in filteredCommunications"
				:key="comm.id"
				class="list-group-item report-card p-3"
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
								comm.type === 'publicreport.nuisance' ? 'bg-danger' : 'bg-info'
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
							comm.public_report.address.postal_code
						}}</span>
					</div>
					<small>{{ formatAddress(comm.public_report.address) }}</small>
					<div
						v-if="
							comm.public_report.images && comm.public_report.images.length > 0
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
import TimeRelative from "../components/TimeRelative.vue";
import { Communication } from "../types";

interface Props {
	all: Communication[] | null;
	loading: boolean;
	selectedId?: string | null;
}
interface Emits {
	(e: "select", id: string): void;
}

const props = withDefaults(defineProps<Props>(), {
	selectedIndex: null,
});

const emit = defineEmits<Emits>();
const handleClick = (id: string) => {
	emit("select", id);
};
const searchFilter = ref("");
const typeFilter = ref("all");

function selectCommunication(communication: Communication) {
	// Emit both events - one for general use, one for v-model
	console.log("selected", communication);
	emit("select-item", communication);
	emit("update:selectedItem", communication);
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
function filterMatches(filter, comm) {
	// Implement your filter logic here
	return true;
}
function formatAddress(a) {
	if (a.number === "" && a.street === "") {
		return "no address provided";
	}
	return `${a.number} ${a.street}, ${a.locality}`;
}
</script>
