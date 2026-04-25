<style scoped>
.map-container {
	background-color: #e9ecef;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-top: 20px;
}
#map {
	height: 500px;
	width: 100%;
	margin-bottom: 10px;
}
#map img {
	max-width: none;
	min-width: 0px;
	height: auto;
}
.search-box {
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	border-radius: 8px;
}
@media (max-width: 768px) {
	.map-container {
		height: 300px;
	}
}
/* Base styles for circular checkboxes */
.custom-circle-checkbox .form-check-input {
	border-radius: 50%;
	width: 20px;
	height: 20px;
	cursor: pointer;
	margin-top: 0.25rem;
	background-image: none; /* Remove Bootstrap's default checkmark */
	position: relative;
}

/* Style for when checked */
.custom-circle-checkbox .form-check-input:checked {
	background-color: currentColor;
	border-color: currentColor;
}

/* Style for when unchecked - just an outline */
.custom-circle-checkbox .form-check-input:not(:checked) {
	background-color: transparent;
}

/* Colors based on data attribute */
.custom-circle-checkbox .form-check-input[data-color="danger"] {
	border-color: $red;
	color: $red;
}

.custom-circle-checkbox .form-check-input[data-color="success"] {
	border-color: $blue;
	color: $blue;
}

.custom-circle-checkbox .form-check-input[data-color="info"] {
	border-color: $green;
	color: $green;
}
</style>
<template>
	<div class="container my-4" v-if="tegola">
		<!-- Search Box -->
		<div class="card search-box mb-4">
			<div class="card-body">
				<form class="row g-3 align-items-center" action="#" id="lookup-form">
					<div class="col-md-9">
						<AddressOrReportSuggestionInput
							name="address-or-report"
							placeholder="Enter a report ID, address, neighborhood, or zip code"
						/>
					</div>
					<div class="col-md-3">
						<span
							data-bs-toggle="tooltip"
							id="lookup-tooltip"
							title="You can look up a report once you type in the full report ID. Start typing and I'll suggest complete IDs"
						>
							<button
								type="submit"
								class="btn btn-primary btn-lg w-100 disabled"
								disabled
								id="lookup"
							>
								Lookup Report by ID
							</button>
						</span>
					</div>
					<div class="col-12">
						<div class="form-check custom-circle-checkbox">
							<input
								class="form-check-input"
								type="checkbox"
								id="checkboxNuisance"
								data-color="danger"
								checked
							/>
							<label class="form-check-label" for="checkboxNuisance"
								>Mosquito Nuisance</label
							>
						</div>

						<div class="form-check custom-circle-checkbox">
							<input
								class="form-check-input"
								type="checkbox"
								id="checkboxWater"
								data-color="success"
								checked
							/>
							<label class="form-check-label" for="checkboxWater"
								>Standing Water</label
							>
						</div>
					</div>
				</form>
			</div>
		</div>

		<!-- Map Section -->
		<div class="card mb-4">
			<div class="card-header bg-info text-white">
				<h5 class="mb-0"><i class="bi bi-pin-map-fill me-2"></i>Reports Map</h5>
			</div>
			<div class="card-body p-0">
				<div class="map-container">
					<Map class="map" tegola="tegola">
						<Layer
							id="nuisance"
							:paint="paintConfigNuisance"
							source="tegola"
							sourceLayer="nuisance_location"
							type="circle"
							v-model="renderedReportsNuisance"
						/>
						<Source id="tegola" type="vector" :tiles="[tegola]" />
					</Map>
				</div>
			</div>
		</div>

		<!-- Results Section -->
		<div class="card">
			<div
				class="card-header bg-primary text-white d-flex justify-content-between align-items-center"
			>
				<h5 class="mb-0">
					<i class="bi bi-geo-fill me-2"></i>Reports Near You
				</h5>
				<span class="badge bg-light text-dark" id="report-count"
					>- Reports Found</span
				>
			</div>
			<div class="card-body p-0">
				<div class="table-responsive">
					<TableReport :reports="renderedReports" />
				</div>
			</div>
			<!--
		<div class="card-footer">
			<nav aria-label="Page navigation">
				<ul class="pagination justify-content-center mb-0">
					<li class="page-item disabled">
						<a class="page-link" href="#" tabindex="-1" aria-disabled="true">Previous</a>
					</li>
					<li class="page-item active"><a class="page-link" href="#">1</a></li>
					<li class="page-item"><a class="page-link" href="#">2</a></li>
					<li class="page-item"><a class="page-link" href="#">3</a></li>
					<li class="page-item">
						<a class="page-link" href="#">Next</a>
					</li>
				</ul>
			</nav>
		</div>
		-->
		</div>
	</div>
	<div v-else>
		<p>loading...</p>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import TableReport, { Report } from "@/rmo/components/TableReport.vue";
import AddressOrReportSuggestionInput from "@/components/AddressOrReportSuggestionInput.vue";
import { apiClient } from "@/client";
import Map from "@/map/Map.vue";
import Layer, { Feature, MouseEvent } from "@/map/Layer.vue";
import Source from "@/map/Source.vue";
import { useStoreAPI } from "@/store/api";

const paintConfigNuisance = {
	"circle-color": "#DC4535",
	"circle-radius": 7,
	"circle-stroke-color": "#9C1C28",
	"circle-stroke-width": 2,
};
const paintConfigWater = {
	"circle-color": "#0D6EfD",
	"circle-radius": 7,
	"circle-stroke-color": "#024AB6",
	"circle-stroke-width": 2,
};
const renderedReportsNuisance = ref<Feature[]>([]);
const storeAPI = useStoreAPI();
const tegola = ref<string | null>(null);

const renderedReports = computed((): Report[] => {
	let reports: Report[] = [];
	renderedReportsNuisance.value.forEach((f) => {
		const p = f.properties;
		reports.push({
			id: p.public_id,
			created: new Date(p.created),
			type: "nuisance",
			address: p.address_raw,
			status: p.status,
		});
	});
	return reports;
});
onMounted(() => {
	const a = storeAPI.get().then((a) => {
		tegola.value = a.tegola.rmo;
	});
});
</script>
