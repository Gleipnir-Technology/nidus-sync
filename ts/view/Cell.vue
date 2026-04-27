<style scoped>
.section-header {
	margin-bottom: 1rem;
	font-size: 1.5rem;
	font-weight: 600;
}

.address-container {
	padding: 1.5rem;
}

.location-divider {
	margin: 1.5rem 0;
}

.coordinates-table {
	width: 100%;
}

.coordinates-table td {
	padding: 0.5rem;
}

.coordinates-table td:first-child {
	width: 40%;
}
</style>
<template>
	<div class="container mt-4 mb-5">
		<!-- Location Header Section -->
		<div class="row mb-4">
			<div class="col-12">
				<h1>Location Data View</h1>
			</div>
		</div>

		<!-- Map and Address Section - Side by Side -->
		<div class="row mb-4">
			<div class="col-md-8">
				<!-- map here -->
			</div>
			<div class="col-md-4">
				<div class="card h-100">
					<div class="card-body address-container">
						<h5>Approximate Address:</h5>
						<p class="lead" id="location-address">
							{{ address }}
						</p>

						<hr class="location-divider" />

						<h5>Cell Coordinates (Hexagon):</h5>
						<div class="table-responsive">
							<table class="coordinates-table">
								<tbody>
									<tr v-for="(coord, index) in cellBoundary" :key="index">
										<td>
											<strong>Vertex {{ index }}:</strong>
										</td>
										<td>{{ formatLatLng(coord) }}</td>
									</tr>
								</tbody>
							</table>
						</div>
						<hr class="location-divider" />
					</div>
				</div>
			</div>
		</div>

		<!-- Two-Column Layout for Tables -->
		<div class="row">
			<!-- Left Column -->
			<div class="col-md-6">
				<!-- Breeding Sources Section -->
				<h2 class="section-header">Mosquito Breeding Sources</h2>
				<div class="card mb-4">
					<div class="card-body">
						<div class="table-responsive">
							<table class="table table-striped table-hover">
								<thead>
									<tr>
										<th>ID</th>
										<th>Source Type</th>
										<th>Last Inspected</th>
										<th>Last Treated</th>
									</tr>
								</thead>
								<tbody>
									<tr v-for="source in breedingSources" :key="source.id">
										<td>
											<a :href="`/source/${source.id}`">{{
												shortUuid(source.id)
											}}</a>
										</td>
										<td>{{ source.type }}</td>
										<td>{{ formatRelativeTime(source.lastInspected) }}</td>
										<td>{{ formatRelativeTime(source.lastTreated) }}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
				</div>

				<!-- Inspections Section -->
				<h2 class="section-header">Inspections History</h2>
				<div class="card mb-4">
					<div class="card-body">
						<div class="table-responsive">
							<table class="table table-striped table-hover">
								<thead>
									<tr>
										<th>LocationID</th>
										<th>Location</th>
										<th>Date</th>
										<th>Action</th>
										<th>Notes</th>
									</tr>
								</thead>
								<tbody>
									<tr
										v-for="inspection in inspections"
										:key="inspection.locationID"
									>
										<td>
											<a :href="`/source/${inspection.locationID}`">
												{{ shortUuid(inspection.locationID) }}
											</a>
										</td>
										<td>{{ inspection.location }}</td>
										<td>{{ formatRelativeTime(inspection.date) }}</td>
										<td>{{ inspection.action }}</td>
										<td>{{ inspection.notes }}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</div>

			<!-- Right Column -->
			<div class="col-md-6">
				<h2 class="section-header">Traps</h2>
				<div v-if="traps.length > 0" class="card mb-4">
					<div class="card-body">
						<div class="table-responsive">
							<table class="table table-striped table-hover">
								<thead>
									<tr>
										<th>ID</th>
										<th>Active</th>
										<th>Comments</th>
									</tr>
								</thead>
								<tbody>
									<tr v-for="trap in traps" :key="trap.globalID">
										<td>
											<a :href="`/trap/${trap.globalID}`">
												{{ shortUuid(trap.globalID) }}
											</a>
										</td>
										<td>{{ trap.active }}</td>
										<td>{{ trap.comments }}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
				</div>
				<p v-else>No traps</p>

				<!-- Treatments Section -->
				<h2 class="section-header">Treatment History</h2>
				<div class="card mb-4">
					<div class="card-body">
						<div class="table-responsive">
							<table class="table table-striped table-hover">
								<thead>
									<tr>
										<th>Location</th>
										<th>Treatment Date</th>
										<th>Insecticide Used</th>
										<th>Technician Notes</th>
									</tr>
								</thead>
								<tbody>
									<tr
										v-for="treatment in treatments"
										:key="treatment.locationID"
									>
										<td>
											<a :href="`/source/${treatment.locationID}`">
												{{ shortUuid(treatment.locationID) }}
											</a>
										</td>
										<td>{{ formatRelativeTime(treatment.date) }}</td>
										<td>{{ treatment.product }}</td>
										<td>{{ treatment.notes }}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";

import Map from "@/map/Map.vue";
import { Address } from "@/type/api";
// Types
interface LatLng {
	lat: number;
	lng: number;
}

interface MapData {
	geoJSON: object;
	center: LatLng;
	zoom: number;
}

interface Organization {
	id: string;
}

interface User {
	organization: Organization;
}

interface BreedingSource {
	id: string;
	type: string;
	lastInspected: Date | null;
	lastTreated: Date | null;
}

interface Inspection {
	locationID: string;
	location: string;
	date: Date | null;
	action: string;
	notes: string;
}

interface Trap {
	globalID: string;
	active: boolean;
	comments: string;
}

interface Treatment {
	locationID: string;
	date: Date | null;
	product: string;
	notes: string;
}

interface Props {
	cell: string;
}

const props = withDefaults(defineProps<Props>(), {});

/*

  mapData: MapData;
  user: User;
  cellBoundary: LatLng[];
  breedingSources: BreedingSource[];
  inspections: Inspection[];
  traps: Trap[];
  treatments: Treatment[];
  address?: string;
*/

const address = ref<Address>(new Address());
const breedingSources = ref<BreedingSource[]>([]);
const cellBoundary = ref<LatLng[]>([]);
const inspections = ref<Inspection[]>([]);
const mapData = ref<MapData>({
	geoJSON: {},
	center: {
		lat: 0,
		lng: 0,
	},
	zoom: 0,
});
const traps = ref<Trap[]>([]);
const treatments = computed((): Treatment[] => {
	return [];
});
// Methods
const formatLatLng = (coord: LatLng): string => {
	return `${coord.lat.toFixed(6)}, ${coord.lng.toFixed(6)}`;
};

const shortUuid = (uuid: string): string => {
	return uuid.split("-")[0];
};

const formatRelativeTime = (date: Date | null): string => {
	if (!date) return "N/A";

	const now = new Date();
	const diff = now.getTime() - new Date(date).getTime();
	const seconds = Math.floor(diff / 1000);
	const minutes = Math.floor(seconds / 60);
	const hours = Math.floor(minutes / 60);
	const days = Math.floor(hours / 24);
	const months = Math.floor(days / 30);
	const years = Math.floor(days / 365);

	if (years > 0) return `${years} year${years > 1 ? "s" : ""} ago`;
	if (months > 0) return `${months} month${months > 1 ? "s" : ""} ago`;
	if (days > 0) return `${days} day${days > 1 ? "s" : ""} ago`;
	if (hours > 0) return `${hours} hour${hours > 1 ? "s" : ""} ago`;
	if (minutes > 0) return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
	return "Just now";
};

const generateGISStatement = (boundary: LatLng[]): string => {
	const coords = boundary
		.map((coord) => `${coord.lng} ${coord.lat}`)
		.join(", ");
	return `POLYGON((${coords}))`;
};

// Lifecycle
onMounted(() => {
	// Load MapLibre GL if not already loaded
	if (!window.maplibregl) {
		const script = document.createElement("script");
		script.src = "//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.js";
		script.async = true;
		document.head.appendChild(script);
	}
});
</script>
