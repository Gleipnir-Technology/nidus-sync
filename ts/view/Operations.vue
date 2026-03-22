<template>
	<div class="operations-command-center">
		<!-- HEADER -->
		<div class="row mb-3 align-items-center">
			<div class="col-md-6">
				<h4 class="fw-bold mb-0">Operations Command Center</h4>
			</div>
			<div class="col-md-6 text-end">
				<button
					class="btn btn-outline-primary me-2"
					@click="addEmergentAssignment"
				>
					Add Emergent Assignment
				</button>
				<button class="btn btn-outline-danger" @click="closeDay">
					Close Day
				</button>
			</div>
		</div>

		<!-- MODE TOGGLE -->
		<ul class="nav nav-tabs mode-toggle mb-4" role="tablist">
			<li class="nav-item">
				<button
					class="nav-link"
					:class="{ active: activeTab === 'planning' }"
					@click="activeTab = 'planning'"
				>
					Planning Mode
				</button>
			</li>
			<li class="nav-item">
				<button
					class="nav-link"
					:class="{ active: activeTab === 'live' }"
					@click="activeTab = 'live'"
				>
					Live Mode
				</button>
			</li>
		</ul>

		<div class="tab-content">
			<!-- ================= PLANNING MODE ================= -->
			<div
				class="tab-pane fade"
				:class="{ 'show active': activeTab === 'planning' }"
			>
				<div class="row mb-3">
					<!-- LEFT: ASSIGNMENTS -->
					<div class="col-lg-3">
						<div class="card">
							<div
								class="card-header d-flex justify-content-between align-items-center fw-semibold"
							>
								<span>Assignments & Work Requests</span>
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										v-model="selectAllAssignments"
										@change="toggleAllAssignments"
									/>
									<label class="form-check-label small">Select All</label>
								</div>
							</div>
							<div class="card-body scroll-panel">
								<input
									class="form-control form-control-sm mb-2"
									placeholder="Filter by section, equipment, expertise"
									v-model="assignmentFilter"
								/>
								<div class="list-group">
									<div
										v-for="assignment in filteredAssignments"
										:key="assignment.id"
										class="list-group-item d-flex"
									>
										<div class="form-check me-2">
											<input
												class="form-check-input"
												type="checkbox"
												v-model="assignment.selected"
											/>
										</div>
										<div>
											<div class="d-flex justify-content-between">
												<strong>{{ assignment.name }}</strong>
												<span
													class="badge"
													:class="{
														'bg-primary': assignment.status === 'Planned',
														'bg-warning text-dark':
															assignment.status === 'Emergent',
													}"
												>
													{{ assignment.status }}
												</span>
											</div>
											<small>{{ assignment.details }}</small>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- CENTER: MAP -->
					<div class="col-lg-6">
						<div class="card">
							<div class="card-header fw-semibold">Routing Map</div>
							<div class="map-placeholder" ref="planningMap">
								Map: Selected Assignments, Selected Technicians, Proposed Routes
							</div>
							<div
								class="card-footer d-flex justify-content-between align-items-center"
							>
								<div
									class="selection-counter small"
									:class="selectionCounterClass"
								>
									{{ selectedAssignmentsCount }} Assignments Selected ·
									{{ selectedTechniciansCount }} Technicians Selected
								</div>
								<div>
									<button class="btn btn-success" @click="computeOptimalRoutes">
										Compute Optimal Routes
									</button>
									<button
										class="btn btn-outline-primary"
										@click="manualAssignment"
									>
										Manual Assignment
									</button>
								</div>
							</div>
						</div>
					</div>

					<!-- RIGHT: TECHNICIANS -->
					<div class="col-lg-3">
						<div class="card">
							<div
								class="card-header d-flex justify-content-between align-items-center fw-semibold"
							>
								<span>Technicians</span>
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										v-model="selectAllTechnicians"
										@change="toggleAllTechnicians"
									/>
									<label class="form-check-label small">Select All</label>
								</div>
							</div>
							<div class="card-body scroll-panel">
								<input
									class="form-control form-control-sm mb-2"
									placeholder="Filter technicians"
									v-model="technicianFilter"
								/>
								<div class="list-group">
									<div
										v-for="technician in filteredTechnicians"
										:key="technician.id"
										class="list-group-item d-flex"
										:class="{ overload: technician.overload }"
									>
										<div class="form-check me-2">
											<input
												class="form-check-input"
												type="checkbox"
												v-model="technician.selected"
											/>
										</div>
										<div>
											<div class="d-flex justify-content-between">
												<strong>{{ technician.name }}</strong>
												<span>
													<span
														class="status-dot"
														:class="{
															'bg-success': technician.status === 'Available',
															'bg-warning': technician.status === 'In Field',
														}"
													></span>
													{{ technician.status }}
												</span>
											</div>
											<small>{{ technician.details }}</small>
											<div class="mt-2">
												<label class="form-label form-label-sm mb-1">
													Assigned Vehicle
												</label>
												<select
													class="form-select form-select-sm"
													v-model="technician.vehicle"
												>
													<option
														v-for="vehicle in vehicles"
														:key="vehicle.id"
														:value="vehicle.id"
													>
														{{ vehicle.name }}
													</option>
												</select>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- ROUTES LIST (STACKED) -->
				<div class="row mb-4">
					<div class="col-12">
						<div class="card">
							<div class="card-header fw-semibold">Proposed Routes</div>
							<div class="card-body">
								<div
									v-for="route in proposedRoutes"
									:key="route.id"
									class="card route-card p-3 mb-3"
								>
									<strong>{{ route.title }}</strong
									><br />
									<small>{{ route.summary }}</small>
									<div class="mt-2">
										<button
											class="btn btn-sm btn-outline-secondary"
											@click="viewAssignments(route.id)"
										>
											View Assignments
										</button>
										<button
											class="btn btn-sm btn-outline-primary"
											@click="modifyRoute(route.id)"
										>
											Modify Route
										</button>
										<button
											class="btn btn-sm btn-outline-secondary"
											@click="shiftAssignment(route.id)"
										>
											Shift Assignment
										</button>
										<button
											class="btn btn-sm btn-outline-secondary"
											@click="swapTechnician(route.id)"
										>
											Swap Technician
										</button>
									</div>
								</div>
							</div>
							<div class="card-footer text-end">
								<button
									class="btn btn-primary send-routes-btn"
									@click="sendRoutes"
								>
									Send Routes to Technicians and Begin Live Operations
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- ================= LIVE MODE ================= -->
			<div
				class="tab-pane fade"
				:class="{ 'show active': activeTab === 'live' }"
			>
				<div class="row mb-3">
					<!-- LEFT: ASSIGNMENTS IN ROUTE ORDER -->
					<div class="col-lg-3">
						<div class="card">
							<div class="card-header fw-semibold">
								Assignments in Route Order
							</div>
							<div class="card-body scroll-panel">
								<div class="alert alert-warning">
									<strong>Unassigned Assignments</strong><br />
									{{ unassignedCount }} awaiting routing
								</div>
								<ul class="list-group list-group-flush">
									<li
										v-for="assignment in liveAssignments"
										:key="assignment.id"
										class="list-group-item"
									>
										<strong>{{ assignment.name }}</strong
										><br />
										{{ assignment.status }}
									</li>
								</ul>
							</div>
						</div>
					</div>

					<!-- CENTER: LIVE MAP -->
					<div class="col-lg-6">
						<div class="live-map" ref="liveMap">
							Live Map: Active Routes, Technician Position, Route Progress
						</div>
					</div>

					<!-- RIGHT: TECHNICIAN STATUS -->
					<div class="col-lg-3">
						<div class="card">
							<div class="card-header fw-semibold">Technician Status</div>
							<div class="card-body scroll-panel">
								<div class="list-group">
									<div
										v-for="technician in liveTechnicians"
										:key="technician.id"
										class="list-group-item"
									>
										<div class="d-flex justify-content-between">
											<strong>{{ technician.name }}</strong>
											<span
												class="badge"
												:class="{
													'bg-success': technician.liveStatus === 'On Track',
													'bg-danger':
														technician.liveStatus === 'Support Requested',
												}"
											>
												{{ technician.liveStatus }}
											</span>
										</div>
										<small>{{ technician.liveDetails }}</small>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- ROUTE INTELLIGENCE -->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-header fw-semibold">
								Active Routes Intelligence
							</div>
							<div class="card-body">
								<table class="table table-sm align-middle">
									<thead>
										<tr>
											<th>Technician</th>
											<th>Assignments</th>
											<th>Estimated Completion</th>
											<th>Remaining Time</th>
											<th>Status</th>
											<th>Actions</th>
										</tr>
									</thead>
									<tbody>
										<tr v-for="route in activeRoutes" :key="route.id">
											<td>{{ route.technician }}</td>
											<td>{{ route.assignmentCount }}</td>
											<td>{{ route.estimatedCompletion }}</td>
											<td>{{ route.remainingTime }}</td>
											<td>
												<span
													class="badge"
													:class="{
														'bg-success': route.status === 'On Track',
														'bg-danger': route.status === 'Blocked',
													}"
												>
													{{ route.status }}
												</span>
											</td>
											<td>
												<button
													class="btn btn-sm btn-outline-secondary"
													@click="viewRoute(route.id)"
												>
													View Route
												</button>
												<button
													class="btn btn-sm btn-outline-primary"
													@click="reallocateOrAssist(route.id)"
												>
													{{
														route.status === "Blocked" ? "Assist" : "Reallocate"
													}}
												</button>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";

// Active tab state
const activeTab = ref("planning");

// Planning Mode - Assignments
const assignmentFilter = ref("");
const selectAllAssignments = ref(false);
const assignments = ref([
	{
		id: 1,
		name: "Larval Habitat Inspection",
		status: "Planned",
		details: "Residential · Backpack Blower",
		selected: false,
	},
	{
		id: 2,
		name: "Green Pool Treatment",
		status: "Emergent",
		details: "Residential · Larvicide · Access Clearance",
		selected: false,
	},
]);

const filteredAssignments = computed(() => {
	if (!assignmentFilter.value) return assignments.value;
	const filter = assignmentFilter.value.toLowerCase();
	return assignments.value.filter(
		(a) =>
			a.name.toLowerCase().includes(filter) ||
			a.details.toLowerCase().includes(filter),
	);
});

const selectedAssignmentsCount = computed(() => {
	return assignments.value.filter((a) => a.selected).length;
});

const toggleAllAssignments = () => {
	assignments.value.forEach((a) => (a.selected = selectAllAssignments.value));
};

// Planning Mode - Technicians
const technicianFilter = ref("");
const selectAllTechnicians = ref(false);
const technicians = ref([
	{
		id: 1,
		name: "Technician A",
		status: "Available",
		details: "Residential · ULV · Backpack",
		vehicle: 1,
		selected: false,
		overload: false,
	},
	{
		id: 2,
		name: "Technician B",
		status: "In Field",
		details: "Agricultural · Capacity Exceeded",
		vehicle: 2,
		selected: false,
		overload: true,
	},
]);

const vehicles = ref([
	{ id: 1, name: "Truck 12 · ULV · Backpack · Larvicide Kit" },
	{ id: 2, name: "ATV 3 · Dipper · Granular Spreader" },
	{ id: 3, name: "Reserve Vehicle · Minimal Equipment" },
]);

const filteredTechnicians = computed(() => {
	if (!technicianFilter.value) return technicians.value;
	const filter = technicianFilter.value.toLowerCase();
	return technicians.value.filter(
		(t) =>
			t.name.toLowerCase().includes(filter) ||
			t.details.toLowerCase().includes(filter),
	);
});

const selectedTechniciansCount = computed(() => {
	return technicians.value.filter((t) => t.selected).length;
});

const toggleAllTechnicians = () => {
	technicians.value.forEach((t) => (t.selected = selectAllTechnicians.value));
};

// Selection counter class
const selectionCounterClass = computed(() => {
	const assignmentCount = selectedAssignmentsCount.value;
	const technicianCount = selectedTechniciansCount.value;

	if (assignmentCount > 0 && technicianCount > 0) {
		return "valid";
	} else if (assignmentCount > 0 || technicianCount > 0) {
		return "warning";
	}
	return "invalid";
});

// Proposed Routes
const proposedRoutes = ref([
	{
		id: 1,
		title: "Route: Technician A",
		summary: "5 Assignments · Est. 4.5 hrs · Equipment: Backpack",
	},
	{
		id: 2,
		title: "Route: Technician B",
		summary: "6 Assignments · Est. 6 hrs · Equipment: ATV",
	},
]);

// Live Mode - Assignments
const unassignedCount = ref(2);
const liveAssignments = ref([
	{ id: 1, name: "Green Pool Reinspection", status: "Communication Pending" },
	{ id: 2, name: "Storm Drain Treatment", status: "In Progress" },
]);

// Live Mode - Technicians
const liveTechnicians = ref([
	{
		id: 1,
		name: "Technician A",
		liveStatus: "On Track",
		liveDetails: "72% Complete · 1.5 hrs Remaining",
	},
	{
		id: 2,
		name: "Technician C",
		liveStatus: "Support Requested",
		liveDetails: "Equipment Issue",
	},
]);

// Live Mode - Active Routes
const activeRoutes = ref([
	{
		id: 1,
		technician: "Technician A",
		assignmentCount: 5,
		estimatedCompletion: "3:45 PM",
		remainingTime: "1 hr 30 min",
		status: "On Track",
	},
	{
		id: 2,
		technician: "Technician C",
		assignmentCount: 4,
		estimatedCompletion: "4:30 PM",
		remainingTime: "2 hrs",
		status: "Blocked",
	},
]);

// Map refs
const planningMap = ref(null);
const liveMap = ref(null);

// Methods
const addEmergentAssignment = () => {
	console.log("Add emergent assignment");
};

const closeDay = () => {
	console.log("Close day");
};

const computeOptimalRoutes = () => {
	console.log("Computing optimal routes...");
};

const manualAssignment = () => {
	console.log("Manual assignment mode");
};

const viewAssignments = (routeId) => {
	console.log("View assignments for route:", routeId);
};

const modifyRoute = (routeId) => {
	console.log("Modify route:", routeId);
};

const shiftAssignment = (routeId) => {
	console.log("Shift assignment for route:", routeId);
};

const swapTechnician = (routeId) => {
	console.log("Swap technician for route:", routeId);
};

const sendRoutes = () => {
	console.log("Sending routes to technicians...");
	activeTab.value = "live";
};

const viewRoute = (routeId) => {
	console.log("View route:", routeId);
};

const reallocateOrAssist = (routeId) => {
	console.log("Reallocate/Assist route:", routeId);
};

// Initialize map when component is mounted
onMounted(() => {
	// Initialize MapLibre GL maps here if needed
	// You'll need to import maplibre-gl separately
	console.log("Component mounted, initialize maps");
});
</script>

<style scoped>
.card {
	box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
}

.status-dot {
	width: 10px;
	height: 10px;
	border-radius: 50%;
	display: inline-block;
	margin-right: 6px;
}

.overload {
	border-left: 4px solid #dc3545;
}

.map-placeholder {
	height: 520px;
	background: #e9ecef;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 600;
	color: #6c757d;
}

.live-map {
	height: 620px;
	background: #dee2e6;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 600;
	color: #6c757d;
}

.scroll-panel {
	max-height: 500px;
	overflow-y: auto;
}

.mode-toggle .nav-link {
	font-weight: 600;
}

.route-card {
	border-left: 4px solid #0d6efd;
}

.send-routes-btn {
	font-size: 1.1rem;
	padding: 0.75rem 1.5rem;
}

.selection-counter.valid {
	color: #198754;
}

.selection-counter.invalid {
	color: #dc3545;
}

.selection-counter.warning {
	color: #fd7e14;
}
</style>
