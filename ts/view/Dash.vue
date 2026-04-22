<template>
	<!-- Dashboard Header -->
	<div class="row mb-4">
		<div class="col-md-6">
			<h1>{{ session?.organization?.name }} Dashboard</h1>
			<p class="text-muted">
				Hey {{ session?.self?.display_name }}, here's an overview of mosquito
				control activities in your district
			</p>
		</div>
		<div
			class="col-md-6 text-md-end d-flex align-items-center justify-content-md-end"
		>
			<p v-if="dashboard.isSyncOngoing" class="last-refreshed mb-0">
				<i class="fas fa-sync-alt me-2 syncing"></i>Syncing now...
			</p>
			<p v-else class="last-refreshed mb-0">
				<i class="fas fa-sync-alt me-2"></i>Last updated:
				<span id="last-refreshed-time">{{
					formatTimeRelative(dashboard.lastSync)
				}}</span>
				<button
					class="btn btn-sm btn-outline-primary ms-3"
					@click="refreshData"
				>
					Refresh Data
				</button>
			</p>
		</div>
	</div>

	<!-- Key Metrics -->
	<div class="row g-4">
		<!-- Last Refreshed -->
		<div class="col-md-3">
			<div class="card stats-card h-100">
				<div class="card-body text-center">
					<div class="metric-icon bg-success text-white">
						<i class="fas fa-clock"></i>
					</div>
					<h5 class="card-title">Last Data Refresh</h5>
					<p class="metric-value">
						{{ formatTimeRelative(dashboard.lastSync) }}
					</p>
					<!-- <p class="card-text text-muted">Last sync: 12:45 PM</p> -->
				</div>
			</div>
		</div>

		<!-- Service Requests -->
		<div class="col-md-3">
			<div class="card stats-card h-100">
				<div class="card-body text-center">
					<div class="metric-icon bg-warning text-white">
						<i class="fas fa-ticket-alt"></i>
					</div>
					<h5 class="card-title">Service Requests</h5>
					<p v-if="dashboard.isSyncOngoing" class="metric-value">
						{{ formatBigNumber(dashboard.counts.service_requests) }}...?
					</p>
					<p v-else class="metric-value">
						{{ formatBigNumber(dashboard.counts.service_requests) }}
					</p>
					<!--<p class="card-text text-muted">
          <span class="text-success">
            <i class="fas fa-arrow-up"></i> 12%
          </span> since last week
        </p>-->
				</div>
			</div>
		</div>

		<!-- Mosquito Sources -->
		<div class="col-md-3">
			<div class="card stats-card h-100">
				<div class="card-body text-center">
					<div class="metric-icon bg-danger text-white">
						<i class="fas fa-bug"></i>
					</div>
					<h5 class="card-title">Mosquito Sources</h5>
					<p v-if="dashboard.isSyncOngoing" class="metric-value">
						{{ formatBigNumber(dashboard.counts.mosquito_sources) }}..?
					</p>
					<p v-else class="metric-value">
						{{ formatBigNumber(dashboard.counts.mosquito_sources) }}
					</p>
					<!-- <p class="card-text text-muted">
          <span class="text-danger">
            <i class="fas fa-arrow-up"></i> 8%
          </span> since last month
        </p> -->
				</div>
			</div>
		</div>

		<!-- Inspections -->
		<div class="col-md-3">
			<div class="card stats-card h-100">
				<div class="card-body text-center">
					<div class="metric-icon bg-info text-white">
						<i class="fas fa-clipboard-check"></i>
					</div>
					<h5 class="card-title">Traps</h5>
					<p v-if="dashboard.isSyncOngoing" class="metric-value">
						{{ formatBigNumber(dashboard.counts.traps) }}...?
					</p>
					<p v-else class="metric-value">
						{{ formatBigNumber(dashboard.counts.traps) }}
					</p>
					<!-- <p class="card-text text-muted">
          <span class="text-success">
            <i class="fas fa-arrow-up"></i> 15%
          </span> since last week
        </p> -->
				</div>
			</div>
		</div>
	</div>

	<!-- Map Section -->
	<h3 class="section-title">Mosquito Activity Heatmap</h3>
	<div class="row">
		<div class="col-12">
			<MapAggregate
				:bounds="mapBounds()"
				@cell-click="doClickMap"
				:markers="[]"
				:organizationId="session.organization?.id ?? 1"
				:tegola="session.urls?.tegola ?? ''"
			/>
		</div>
	</div>

	<!-- Recent Activity Section -->
	<h3 class="section-title">Recent Activity</h3>
	<div class="row">
		<div class="col-12">
			<div class="card">
				<div class="card-body">
					<div class="table-responsive">
						<table class="table table-hover">
							<thead>
								<tr>
									<th>Date</th>
									<th>Type</th>
									<th>Location</th>
									<th>Status</th>
									<th>Action</th>
								</tr>
							</thead>
						</table>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive } from "vue";
import MapAggregate from "@/components/MapAggregate.vue";
import { formatBigNumber, formatTimeRelative } from "@/format";
import { router } from "@/route/config";
import { useSessionStore } from "@/store/session";
import { useStoreServiceRequest } from "@/store/service_request";
import { useStoreSync } from "@/store/sync";
import type { Bounds } from "@/type/api";

const dashboard = reactive({
	counts: {
		service_requests: 0,
		mosquito_sources: 0,
		traps: 0,
	},
	organization: {
		name: "",
		id: 0,
	},
	isSyncOngoing: false,
	lastSync: new Date(),
	tegolaUrl: "",
	serviceArea: {
		min: { x: 0, y: 0 },
		max: { x: 0, y: 0 },
	},
	recentRequests: [],
});
const storeServiceRequest = useStoreServiceRequest();
const storeSync = useStoreSync();
const session = useSessionStore();
onMounted(async () => {
	const service_requests = await storeServiceRequest.fetchAll();
	const syncs = await storeSync.fetchAll();
	console.log("syncs", syncs);
});
function doClickMap(cell: string) {
	router.push("/cell/" + cell);
}
function mapBounds(): Bounds | undefined {
	if (session.organization?.service_area) {
		return session.organization?.service_area;
	}
	return undefined;
}
function refreshData() {
	console.log("fake refresh");
}
</script>
