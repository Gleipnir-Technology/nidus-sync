<style scoped lang="scss">
.results-container {
	max-width: 1400px;
}

.summary-card {
	transition: transform 0.2s;
}

.summary-card:hover {
	transform: translateY(-5px);
}

.has-error {
	background-color: #fff3cd;
}

.badge.status,
.badge.condition {
	font-size: 0.875rem;
	padding: 0.35em 0.65em;
}

.table-responsive {
	max-height: 600px;
	overflow-y: auto;
}

thead tr.header {
	position: sticky;
	top: 0;
	z-index: 10;
	background-color: #f8f9fa;
}

.map-container {
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 400px;
	margin-bottom: 20px;
	margin-top: 20px;
	align-items: center;
	justify-content: center;
	/* Prevent touch scrolling issues */
	touch-action: pan-y pinch-zoom;
}
.badge.dry {
	background-color: $info;
}
.badge.empty {
	background-color: #9c9bc0;
}
.badge.false.pool {
	background-color: #6b2737;
}
.badge.green {
	background-color: #4b6827;
}
.badge.murky {
	background-color: #88bc4e;
}
.badge.unknown {
	background-color: gray;
}
.badge.committed {
	background-color: $success;
}
.badge.committing {
	background-color: $success;
}
.badge.discarded {
	background-color: gray;
}
.badge.error {
	background-color: $danger;
}
.badge.parsed {
	background-color: $secondary;
}
.badge.uploaded {
	background-color: $info;
}
.summary-card {
	transition: transform 0.2s;
}
.summary-card:hover {
	transform: translateY(-5px);
}
.badge.status {
	font-size: 0.85rem;
}
.badge.status.existing {
	background-color: $secondary;
}
.badge.status.new {
	background-color: $primary;
}
.badge.status.outside {
	background-color: $warning;
}
.badge.status.unknown {
	background-color: gray;
}
tr.has-error {
	background-color: rgba(255, 193, 7, 0.15) !important;
}
</style>
<template>
	<div class="container mt-4 results-container">
		<div v-if="upload">
			<div class="d-flex justify-content-between align-items-center mb-4">
				<h2>Upload Results: {{ upload.filename ?? "" }}</h2>
				<span class="badge rounded-pill" :class="upload.status">
					<i class="bi me-1" :class="getUploadStatusIcon(upload.status)"></i>
					{{ getUploadStatusDisplay(upload.status) }}
				</span>
			</div>

			<div class="row mb-4">
				<div class="col-md-4">
					<div class="card summary-card h-100 border-primary">
						<div class="card-body text-center">
							<h1 class="display-4 text-primary">
								{{ upload.csv_pool?.count.existing }}
							</h1>
							<h5>Existing Pools</h5>
							<p class="text-muted">Matches found in previous records</p>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="card summary-card h-100 border-success">
						<div class="card-body text-center">
							<h1 class="display-4 text-success">
								{{ upload.csv_pool?.count.new }}
							</h1>
							<h5>New Pools</h5>
							<p class="text-muted">Not found in existing records</p>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="card summary-card h-100 border-warning">
						<div class="card-body text-center">
							<h1 class="display-4 text-warning">
								{{ upload.csv_pool?.count.outside }}
							</h1>
							<h5>Outside District</h5>
							<p class="text-muted">Potential geocoding errors</p>
						</div>
					</div>
				</div>
			</div>

			<div class="card mb-4">
				<div v-if="!(session.organization && session.self)">
					<p>loading</p>
				</div>
				<div
					class="map-container"
					v-show="session.organization && session.self"
				>
					<MapLocator
						:markers="markers"
						:organizationId="session.organization!.id"
						:tegola="session.urls?.tegola ?? ''"
						v-model="mapCamera"
					/>
				</div>
			</div>

			<div class="card mb-4">
				<div
					class="card-header bg-light d-flex justify-content-between align-items-center"
				>
					<h5 class="mb-0">Data Preview</h5>
					<div class="form-check form-switch">
						<input
							class="form-check-input"
							type="checkbox"
							id="showIssuesOnly"
							v-model="showIssuesOnly"
							@change="handleShowIssuesOnly"
						/>
						<label class="form-check-label" for="showIssuesOnly">
							Show issues only
						</label>
					</div>
				</div>
				<div class="card-body">
					<div
						v-for="error in upload.csv_pool?.errors"
						:key="error.message"
						class="alert alert-danger"
						role="alert"
					>
						<i class="bi bi-exclamation-triangle me-2"></i>
						<strong>Error:</strong> {{ error.message }}
					</div>

					<div v-if="upload == null">Loading...</div>
					<div
						v-else-if="
							upload.status === 'uploaded' || upload.status === 'parsing'
						"
						class="alert alert-info"
						role="alert"
					>
						<i class="bi bi-exclamation-triangle me-2"></i>
						<strong>Working:</strong> File is still processing... refresh this
						page in a bit to see updates.
					</div>

					<template v-else>
						<div v-if="upload.error" class="alert alert-error" role="alert">
							<i class="bi bi-exclamation-triangle me-2"></i>
							<strong>Error:</strong> Your upload failed to parse correctly. The
							specific error was: '{{ upload.error }}'
						</div>
						<div
							v-else-if="
								!upload.csv_pool?.pools || upload.csv_pool.pools.length === 0
							"
							class="alert alert-warning"
							role="alert"
						>
							<i class="bi bi-exclamation-triangle me-2"></i>
							<strong>Warning:</strong> No pools could be understood from your
							file.
						</div>

						<div v-else class="table-responsive">
							<table class="table table-hover table-striped">
								<thead class="table-light">
									<tr class="header">
										<th></th>
										<th>Number</th>
										<th>Street</th>
										<th>City</th>
										<th>Postal</th>
										<th>Status</th>
										<th>Condition</th>
										<th>Tags</th>
									</tr>
								</thead>
								<tbody>
									<tr
										v-for="(pool, index) in upload.csv_pool.pools"
										:key="index"
										:class="{
											'has-error': hasError(pool),
										}"
										:style="getRowStyle(pool)"
									>
										<td>
											<Tooltip
												placement="top"
												:title="errorTooltip(pool)"
												v-show="hasError(pool)"
											>
												<i class="bi bi-info-circle-fill text-primary ms-1"></i>
											</Tooltip>
										</td>
										<td>{{ pool.address?.number }}</td>
										<td>{{ pool.address?.street }}</td>
										<td>{{ pool.address?.locality }}</td>
										<td>{{ pool.address?.postal_code }}</td>
										<td>
											<span class="badge status" :class="pool.status">
												{{ titleCase(pool.status) }}
											</span>
										</td>
										<td>
											<span class="badge condition" :class="pool.condition">
												{{ titleCase(pool.condition) }}
											</span>
										</td>
										<td>
											<ul>
												<li v-for="(v, k) in pool.tags">{{ k }}={{ v }}</li>
											</ul>
										</td>
									</tr>
								</tbody>
							</table>
						</div>
					</template>
				</div>
			</div>

			<div class="d-flex justify-content-between mt-4 mb-5">
				<button
					type="button"
					class="btn btn-outline-danger"
					@click="handleDiscard"
					:disabled="isSubmitting"
				>
					Discard
				</button>
				<button
					class="btn btn-primary"
					id="confirmUploadBtn"
					@click="handleCommit"
					:disabled="isSubmitting"
				>
					<i class="bi bi-check2 me-1"></i> Confirm and Commit Data
				</button>
			</div>
		</div>
		<div v-else>
			<p>loading...</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import * as bootstrap from "bootstrap";
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import MapLocator from "@/components/MapLocator.vue";
import Tooltip from "@/components/Tooltip.vue";
import { useUploadStore } from "@/store/upload";
import { useSessionStore } from "@/store/session";
import {
	CSVPoolDetail,
	CSVPoolError,
	Upload,
	UploadPoolError,
	UploadPoolRow,
} from "@/type/api";
import { Camera } from "@/type/map";
import type { Marker } from "@/types";

interface ErrorMessage {
	message: string;
}

interface Props {
	id: number;
}

const props = defineProps<Props>();

const mapCamera = ref<Camera>(new Camera());
const router = useRouter();
const showIssuesOnly = ref(false);
const isSubmitting = ref(false);
const uploadStore = useUploadStore();
const session = useSessionStore();

const upload = ref<Upload | null>(null);

const getUploadStatusIcon = (status?: string): string => {
	const icons: Record<string, string> = {
		uploaded: "bi-cloud-upload",
		parsing: "bi-hourglass-split",
		parsed: "bi-check-circle",
		error: "bi-x-circle",
	};
	return icons[status || ""] || "bi-question-circle";
};

const getUploadStatusDisplay = (status?: string): string => {
	const displays: Record<string, string> = {
		committed: "Committed",
		committing: "Committing",
		discarded: "Discarded",
		error: "Error",
		parsed: "Parsed",
		parsing: "Parsing",
		uploaded: "Uploaded",
	};
	return displays[status || ""] || "Unknown";
};

const titleCase = (str?: string): string => {
	if (!str) return "";
	return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
};
const getRowStyle = (pool: UploadPoolRow) => {
	if (showIssuesOnly.value) {
		const hasError = pool.errors && pool.errors.length > 0;
		return { display: hasError ? "table-row" : "none" };
	}
	return { display: "table-row" };
};

const handleShowIssuesOnly = () => {
	// The reactive display is handled by getRowStyle
};

const initializeMap = () => {
	/*
	if (!map) return;

	map.addEventListener("load", () => {
		map.addSource("tegola-nidus", {
			type: "vector",
			tiles: [
				`${props.tegolaUrl}maps/nidus/{z}/{x}/{y}?csv_file=${props.id}&id=${user.organization.id}`,
			],
		});

		map.addLayer({
			id: "pool",
			source: "tegola-nidus",
			"source-layer": "fileupload-pool",
			type: "circle",
			paint: {
				"circle-color": "#91b979",
				"circle-radius": 7,
				"circle-stroke-width": 2,
				"circle-stroke-color": "#7aab5f",
			},
		});
	});
	*/
};

const handleDiscard = async () => {
	isSubmitting.value = true;
	try {
		const response = await fetch(`/api/upload/${props.id}/discard`, {
			body: JSON.stringify({}),
			headers: {
				"Content-Type": "application/json",
			},
			method: "POST",
		});

		if (!response.ok) throw new Error("Failed to discard upload");

		// Navigate to uploads list or appropriate page
		router.push("/_/configuration/upload");
	} catch (error) {
		console.error("Error discarding upload:", error);
		alert("Failed to discard upload. Please try again.");
	} finally {
		isSubmitting.value = false;
	}
};

const handleCommit = async () => {
	isSubmitting.value = true;
	try {
		const response = await fetch(`/api/upload/${props.id}/commit`, {
			body: JSON.stringify({}),
			headers: {
				"Content-Type": "application/json",
			},
			method: "POST",
		});

		if (!response.ok) throw new Error("Failed to confirm upload");

		// Navigate to success page or appropriate page
		router.push("/_/configuration/upload");
	} catch (error) {
		console.error("Error confirming upload:", error);
		alert("Failed to confirm upload. Please try again.");
	} finally {
		isSubmitting.value = false;
	}
};
const markers = computed((): Marker[] => {
	if (!upload.value?.csv_pool?.pools) {
		return [];
	}
	let markers: Marker[] = [];
	upload.value.csv_pool.pools.forEach((p: UploadPoolRow) => {
		if (p.address.location) {
			markers.push({
				color: "#FF0000",
				draggable: true,
				id: "x",
				location: p.address.location,
			});
		}
	});
	console.log("updated markers to", markers);
	return markers;
});
function hasError(row: UploadPoolRow): boolean {
	return row.errors.length > 0;
}
function errorTooltip(row: UploadPoolRow): string {
	return row.errors.map((e) => e.message).join(", ");
}
onMounted(() => {
	initializeMap();
	uploadStore.fetchOne(props.id).then((u) => {
		console.log("got upload", u);
		upload.value = u;
	});
});
</script>
