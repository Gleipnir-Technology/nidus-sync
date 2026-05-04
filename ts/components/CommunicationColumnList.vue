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
				<div class="loading list-group-item" v-if="loading || all == null">
					Loading...
				</div>
				<div
					v-else-if="all.length > 0"
					v-for="comm in filteredCommunications"
					:key="comm.id"
				>
					<ListCardCommunication
						@click="handleClick(comm.id)"
						:comm="comm"
						:isSelected="selectedID == comm.id"
					/>
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
const filteredCommunications = computed((): Communication[] => {
	if (props.all == null) {
		return [];
	}
	return props.all;
});
</script>
