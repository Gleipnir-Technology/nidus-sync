<template>
	<div class="container-fluid py-3">
		<!-- Header -->
		<div class="row mb-3">
			<div class="col">
				<h3 class="mb-1">Daily Planning Workbench</h3>
				<div class="text-muted small">
					Signals and leads enter from the left, are investigated in the center,
					and transformed into structured field assignments using tools on the
					right.
				</div>
			</div>
		</div>

		<div class="row g-3">
			<!-- LEFT: Incoming Signals & Leads -->
			<div class="col-xl-3">
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
								@click="loadData()"
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
						<div v-if="loading && signals.length === 0" class="loading-spinner">
							<div class="spinner-border text-primary" role="status">
								<span class="visually-hidden">Loading...</span>
							</div>
						</div>

						<!-- Signals -->
						<div class="mb-3" v-show="!loading || signals.length > 0">
							<div class="fw-semibold mb-2">
								Signals
								<span
									class="badge bg-primary"
									v-show="selectedSignals.length > 0"
								>
									{{ selectedSignals.length }}
								</span>
							</div>

							<div
								v-if="signals.length === 0 && !loading"
								class="text-muted small fst-italic"
							>
								No signals found
							</div>

							<div
								v-for="signal in signals"
								:key="signal.id"
								class="border rounded p-2 mb-2 signal-item"
								:class="{ selected: isSelected(signal.id) }"
								@click="toggleSignal(signal)"
							>
								<div class="small fw-semibold">{{ signal.title }}</div>
								<div class="signal-address">
									{{ shortAddress(signal.address) }}
								</div>
								<div class="text-muted small">{{ signal.description }}</div>
								<span v-if="signal.badge" class="badge bg-secondary mt-1">
									{{ signal.badge }}
								</span>
							</div>
						</div>

						<hr />

						<!-- Mosquito Control Plan Followups -->
						<div class="mb-3" v-show="!loading || planFollowups.length > 0">
							<div class="fw-semibold mb-2">
								Mosquito Control Plan Follow-Ups
							</div>

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
			</div>

			<!-- CENTER: Active Workbench -->
			<div class="col-xl-6">
				<div class="card shadow-sm mb-3">
					<div class="card-header bg-white pane-header">
						Active Investigation Workbench
					</div>
					<div class="card-body">
						<div class="map-container">
							<map-multipoint
								ref="mapMultipoint"
								id="map"
								:organization-id="organizationId"
								:tegola="tegolaUrl"
								:xmin="serviceArea.xmin"
								:ymin="serviceArea.ymin"
								:xmax="serviceArea.xmax"
								:ymax="serviceArea.ymax"
							></map-multipoint>
						</div>

						<div class="row g-3">
							<div class="col-md-12">
								<div class="card border">
									<div class="card-body">
										<div class="fw-semibold">Selected Signals</div>
										<div class="text-muted small">
											{{ selectedSignals.length }} Signal{{
												selectedSignals.length !== 1 ? "s" : ""
											}}
											Selected
										</div>

										<div
											v-if="selectedSignals.length === 0"
											class="text-muted small mt-2 fst-italic"
										>
											Click signals from the left panel to select them
										</div>

										<table
											class="small mt-2 table"
											v-show="selectedSignals.length > 0"
										>
											<tbody>
												<tr v-for="signal in selectedSignals" :key="signal.id">
													<td>
														<button
															@click="toggleSignal(signal)"
															class="btn btn-sm btn-link text-danger p-0 ms-1"
															style="font-size: 0.7rem"
														>
															<i class="bi bi-x"></i>
														</button>
													</td>
													<td>
														<span v-if="signal.type === 'flyover pool'"
															>Pool</span
														>
														<span
															v-else-if="
																signal.type === 'publicreport nuisance'
															"
															>Nuisance</span
														>
														<span
															v-else-if="signal.type === 'publicreport water'"
															>Water</span
														>
													</td>
													<td>
														<time-relative
															:time="signal.created"
														></time-relative>
													</td>
													<td>{{ shortAddress(signal.address) }}</td>
												</tr>
											</tbody>
										</table>

										<button
											v-show="selectedSignals.length > 0"
											@click="clearSelection()"
											class="btn btn-sm btn-outline-secondary mt-2 w-100"
										>
											Clear Selection
										</button>
									</div>
								</div>
							</div>
						</div>

						<div v-show="showMapTile" class="map-container">
							<map-proxied-arcgis-tile
								ref="mapTile"
								class="map"
								:organization-id="organizationId"
								:tegola="tegolaUrl"
								:url-tiles="urlTiles"
								:latitude="selectedSignalLocation()?.latitude ?? 0.0"
								:longitude="selectedSignalLocation()?.longitude ?? 0.0"
								@click="updateSignalLocation"
							>
							</map-proxied-arcgis-tile>
						</div>
					</div>
				</div>
			</div>

			<!-- RIGHT: Transformation Tools -->
			<div class="col-xl-3">
				<div class="card shadow-sm h-100">
					<div class="card-header bg-white pane-header">
						Transformation Tools
					</div>
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
							<button class="btn btn-outline-dark tool-button">
								Set Priority
							</button>
							<button class="btn btn-outline-dark tool-button">
								Estimate Effort
							</button>
							<button class="btn btn-outline-dark tool-button">
								Send to Operations
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted, nextTick } from "vue";

// Props
const props = defineProps({
	organizationId: {
		type: String,
		required: true,
	},
	tegolaUrl: {
		type: String,
		required: true,
	},
	urlTiles: {
		type: String,
		required: true,
	},
	serviceArea: {
		type: Object,
		required: true,
		default: () => ({ xmin: 0, ymin: 0, xmax: 0, ymax: 0 }),
	},
});

// Refs
const mapMultipoint = ref(null);
const mapTile = ref(null);

// State
const apiBase = ref("/api");
const creating = ref(false);
const error = ref(null);
const planFollowups = ref([]);
const leads = ref([]);
const loading = ref(false);
const poolLocations = ref({});
const showMapTile = ref(false);
const selectedSignals = ref([]);
const signals = ref([]);

const filters = ref({
	species: "",
	type: "",
	sort: "newest",
});

// Helper functions (outside component)
const getBoundingBox = (points) => {
	if (!points || points.length === 0) {
		return null;
	}

	let minLat = points[0].latitude;
	let maxLat = points[0].latitude;
	let minLng = points[0].longitude;
	let maxLng = points[0].longitude;

	for (const point of points) {
		if (point.latitude < minLat) minLat = point.latitude;
		if (point.latitude > maxLat) maxLat = point.latitude;
		if (point.longitude < minLng) minLng = point.longitude;
		if (point.longitude > maxLng) maxLng = point.longitude;
	}

	return new window.maplibregl.LngLatBounds(
		new window.maplibregl.LngLat(minLng, minLat),
		new window.maplibregl.LngLat(maxLng, maxLat),
	);
};

const shortAddress = (a) => {
	if (!a) return "";
	return `${a.number} ${a.street}, ${a.locality}`;
};

const updateMap = (signals) => {
	if (!mapMultipoint.value) return;

	const locations = signals.map((s) => s.location);
	const markers = locations.map((l) =>
		new window.maplibregl.Marker({
			color: "#FF0000",
			draggable: false,
		}).setLngLat([l.longitude, l.latitude]),
	);

	mapMultipoint.value.SetMarkers(markers);

	const bounds = getBoundingBox(locations);
	if (bounds != null) {
		mapMultipoint.value.FitBounds(bounds, {
			padding: 50,
		});
	}
};

const configureMapTile = () => {
	if (!mapTile.value) return;

	mapTile.value.on("load", () => {
		mapTile.value.addLayer({
			id: "parcel",
			minzoom: 14,
			paint: {
				"line-color": "#0f0",
			},
			source: "tegola",
			"source-layer": "parcel",
			type: "line",
		});
		mapTile.value.addLayer({
			id: "signal-point",
			paint: {
				"circle-color": "#0D6EfD",
				"circle-radius": 7,
				"circle-stroke-width": 2,
				"circle-stroke-color": "#024AB6",
			},
			source: "tegola",
			"source-layer": "signal-point",
			type: "circle",
		});
	});
};

// Methods
const loadData = async () => {
	loading.value = true;
	error.value = null;

	try {
		await Promise.all([loadSignals(), loadLeads()]);
	} catch (err) {
		error.value = err.message;
		console.error("Error loading data:", err);
	} finally {
		loading.value = false;
	}
};

const loadSignals = async () => {
	try {
		const params = new URLSearchParams();
		if (filters.value.species) params.append("species", filters.value.species);
		if (filters.value.type) params.append("type", filters.value.type);
		if (filters.value.sort) params.append("sort", filters.value.sort);

		const response = await fetch(`${apiBase.value}/signal?${params}`);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		signals.value = data.signals || data;
	} catch (err) {
		console.error("Error loading signals:", err);
		throw err;
	}
};

const loadPlanFollowups = async () => {
	try {
		const response = await fetch(`${apiBase.value}/plan-followups`);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		planFollowups.value = data.followups || data;
	} catch (err) {
		console.error("Error loading plan followups:", err);
		throw err;
	}
};

const loadLeads = async () => {
	try {
		const response = await fetch(`${apiBase.value}/leads`);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		leads.value = data.leads || data;
	} catch (err) {
		console.error("Error loading leads:", err);
		throw err;
	}
};

const isSelected = (id) => {
	return selectedSignals.value.some((s) => s.id === id);
};

const clearSelection = () => {
	selectedSignals.value = [];
};

const createLead = async () => {
	if (selectedSignals.value.length === 0) return;

	creating.value = true;

	try {
		const response = await fetch(`${apiBase.value}/leads`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				pool_locations: poolLocations.value,
				signal_ids: selectedSignals.value.map((s) => s.id),
			}),
		});

		if (!response.ok) {
			const errorData = await response.json();
			throw new Error(
				errorData.message || `HTTP error! status: ${response.status}`,
			);
		}

		const newLead = await response.json();
		leads.value.unshift(newLead);
		clearSelection();
		await loadSignals();
	} catch (err) {
		console.error("Error creating lead:", err);
		alert(`Failed to create lead: ${err.message}`);
	} finally {
		creating.value = false;
	}
};

const markAsAddressed = async () => {
	if (selectedSignals.value.length === 0) return;

	try {
		const response = await fetch(`${apiBase.value}/signal/mark-addressed`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				signal_ids: selectedSignals.value.map((s) => s.id),
			}),
		});

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		signals.value = signals.value.filter(
			(signal) => !selectedSignals.value.some((s) => s.id === signal.id),
		);

		clearSelection();
		alert("Signals marked as addressed");
	} catch (err) {
		console.error("Error marking signals as addressed:", err);
		alert(`Failed to mark signals: ${err.message}`);
	}
};

const selectedSignalLocation = () => {
	const first_pool = selectedSignals.value.reduce((accumulator, current) => {
		if (accumulator == null && current.type === "flyover pool") {
			return current;
		}
		return accumulator;
	}, null);
	return first_pool?.location;
};

const toggleSignal = (signal) => {
	const index = selectedSignals.value.findIndex((s) => s.id === signal.id);

	if (index > -1) {
		selectedSignals.value.splice(index, 1);
	} else {
		selectedSignals.value.push(signal);
	}

	updateMap(selectedSignals.value);

	showMapTile.value = selectedSignals.value.reduce(
		(accumulator, current) => accumulator || current.type === "flyover pool",
		false,
	);

	console.log("show tile", showMapTile.value);

	if (showMapTile.value) {
		nextTick(() => {
			configureMapTile();
		});
	}
};

const updateSignalLocation = (event) => {
	const signalId = event.detail.signalId;
	console.log("map click", signalId, event.detail);

	const map = event.detail.map;
	const loc = {
		latitude: event.detail.lat,
		longitude: event.detail.lng,
	};

	map.SetMarkers([loc]);
	poolLocations.value[signalId] = loc;
};

// Lifecycle
onMounted(() => {
	loadData();

	// Subscribe to SSE events
	if (window.SSEManager) {
		window.SSEManager.subscribe("*", (e) => {
			if (e.resource === "sync:signal") {
				loadData();
			}
		});
	}

	// Configure map multipoint
	const map = mapMultipoint.value;
	if (map) {
		map.on("load", () => {
			map.addLayer({
				id: "parcel",
				minzoom: 14,
				paint: {
					"line-color": "#0f0",
				},
				source: "tegola",
				"source-layer": "parcel",
				type: "line",
			});
			map.addLayer({
				id: "signal-point",
				paint: {
					"circle-color": "#0D6EfD",
					"circle-radius": 7,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#024AB6",
				},
				source: "tegola",
				"source-layer": "signal-point",
				type: "circle",
			});
			console.log("Added parcel and signal layers");
		});
	}
});
</script>

<style scoped>
/* Add any component-specific styles here */
.map-container {
	height: 400px;
	margin-bottom: 1rem;
}

.scroll-pane {
	max-height: calc(100vh - 200px);
	overflow-y: auto;
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

.signal-address {
	font-size: 0.875rem;
	color: #6c757d;
}

.tool-button {
	width: 100%;
	margin-bottom: 0.5rem;
	text-align: left;
}

.loading-spinner {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 2rem;
}

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

.pane-header {
	font-weight: 600;
	border-bottom: 2px solid #dee2e6;
}
</style>
