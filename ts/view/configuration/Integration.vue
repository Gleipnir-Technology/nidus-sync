<template>
	<div class="container py-4">
		<!-- Header -->
		<div class="mb-4">
			<h1>Integrations</h1>
			<div class="alert alert-warning">
				<i class="bi bi-exclamation-triangle me-2"></i>
				<strong>Important:</strong> This page allows you to configure
				integration with third-party services. The credentials and tokens stored
				here provide access to external systems and should be protected. Only
				authorized personnel should modify these settings.
			</div>
		</div>

		<!-- Esri ArcGIS Integration -->
		<div class="card mb-4 integration-card arcgis">
			<div
				class="card-header d-flex justify-content-between align-items-center"
			>
				<div>
					<h2 class="h5 mb-0">Esri's ArcGIS</h2>
				</div>
				<img
					src="https://via.placeholder.com/100x40?text=ArcGIS"
					alt="ArcGIS Logo"
					height="40"
				/>
			</div>
			<div class="card-body">
				<div class="table-responsive mb-3">
					<table class="table table-borderless">
						<tbody>
							<tr v-if="!arcGISConfig">
								<td>Not integrated</td>
							</tr>
							<template v-else>
								<tr>
									<td width="30%"><strong>OAuth Token Status</strong></td>
									<td>None</td>
									<td>
										<span
											v-if="arcGISConfig.invalidatedAt"
											class="status-inactive"
										>
											<i class="bi bi-x-circle-fill me-1"></i> Invalidated
										</span>
										<span
											v-else-if="
												isTokenExpired(arcGISConfig.accessTokenExpires)
											"
											class="status-inactive"
										>
											<i class="bi bi-x-circle-fill me-1"></i> Expired
										</span>
										<span v-else class="status-active">
											<i class="bi bi-check-circle-fill me-1"></i> Active
										</span>
									</td>
								</tr>
								<tr>
									<td><strong>Token Expiration</strong></td>
									<td>
										{{ formatRelativeTime(arcGISConfig.accessTokenExpires) }}
									</td>
								</tr>
								<tr>
									<td><strong>Integration Method</strong></td>
									<td>Polling</td>
								</tr>
								<tr>
									<td><strong>Permission Level</strong></td>
									<td>Read</td>
								</tr>
							</template>
						</tbody>
					</table>
				</div>
				<div class="d-flex gap-2">
					<a class="btn btn-primary" :href="urls.oauthRefreshArcGIS">
						<i class="bi bi-arrow-repeat me-2"></i>Refresh OAuth Token
					</a>
					<RouterLink to="/configuration/integration/arcgis">
						<button class="btn btn-outline-danger">
							<i class="bi bi-gear me-2"></i>Configure
						</button>
					</RouterLink>
				</div>
			</div>
		</div>

		<!-- FieldSeeker GIS Integration -->
		<div class="card mb-4 integration-card fieldseeker">
			<div
				class="card-header d-flex justify-content-between align-items-center"
			>
				<div>
					<h2 class="h5 mb-0">Frontier Precision's FieldSeeker GIS</h2>
				</div>
				<img
					src="https://via.placeholder.com/100x40?text=FieldSeeker"
					alt="FieldSeeker Logo"
					height="40"
				/>
			</div>
			<div class="card-body">
				<div class="table-responsive mb-3">
					<table class="table table-borderless">
						<tbody></tbody>
					</table>
				</div>
				<div class="d-flex gap-2"></div>
			</div>
		</div>

		<!-- VectorSurv Integration -->
		<div class="card mb-4 integration-card vectorsurv">
			<div
				class="card-header d-flex justify-content-between align-items-center"
			>
				<div>
					<h2 class="h5 mb-0">VectorSurv</h2>
				</div>
				<img
					src="https://via.placeholder.com/100x40?text=VectorSurv"
					alt="VectorSurv Logo"
					height="40"
				/>
			</div>
			<div class="card-body">
				<div class="table-responsive mb-3">
					<table class="table table-borderless">
						<tbody>
							<tr>
								<td width="30%"><strong>API Token</strong></td>
								<td>
									<span class="token-display">{{
										vectorSurvConfig.maskedToken
									}}</span>
								</td>
							</tr>
							<tr>
								<td><strong>Last Synchronization</strong></td>
								<td>{{ vectorSurvConfig.lastSync }}</td>
							</tr>
							<tr>
								<td><strong>Synchronization Status</strong></td>
								<td>
									<span class="status-active">
										<i class="bi bi-check-circle-fill me-1"></i> Active
										(Scheduled daily at 2:00 AM)
									</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div class="d-flex gap-2">
					<button class="btn btn-success" @click="openVectorSurvModal">
						<i class="bi bi-pencil-square me-2"></i>Edit Token
					</button>
					<button
						class="btn btn-outline-danger"
						@click="removeIntegration('vectorsurv')"
					>
						<i class="bi bi-trash me-2"></i>Remove Integration
					</button>
				</div>
			</div>
		</div>

		<!-- VeeMac Integration -->
		<div class="card mb-4 integration-card veemac">
			<div
				class="card-header d-flex justify-content-between align-items-center"
			>
				<div>
					<h2 class="h5 mb-0">VeeMac</h2>
				</div>
				<img
					src="https://via.placeholder.com/100x40?text=VeeMac"
					alt="VeeMac Logo"
					height="40"
				/>
			</div>
			<div class="card-body">
				<div class="table-responsive mb-3">
					<table class="table table-borderless">
						<tbody>
							<tr>
								<td width="30%"><strong>Username</strong></td>
								<td>{{ veeMacConfig.username }}</td>
							</tr>
							<tr>
								<td><strong>Password</strong></td>
								<td>••••••••••••</td>
							</tr>
							<tr>
								<td><strong>Last Synchronization</strong></td>
								<td>{{ veeMacConfig.lastSync }}</td>
							</tr>
							<tr>
								<td><strong>Synchronization Status</strong></td>
								<td>
									<span class="status-inactive">
										<i class="bi bi-x-circle-fill me-1"></i> Inactive (Manual
										sync only)
									</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div class="d-flex gap-2">
					<button class="btn btn-success" @click="openVeeMacModal">
						<i class="bi bi-pencil-square me-2"></i>Edit Credentials
					</button>
					<button
						class="btn btn-outline-danger"
						@click="removeIntegration('veemac')"
					>
						<i class="bi bi-trash me-2"></i>Remove Integration
					</button>
				</div>
			</div>
		</div>
	</div>

	<!-- VectorSurv Edit Token Modal -->
	<div
		class="modal fade"
		ref="vectorsurvModalEl"
		tabindex="-1"
		aria-labelledby="vectorsurvModalLabel"
		aria-hidden="true"
	>
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="vectorsurvModalLabel">
						Edit VectorSurv API Token
					</h5>
					<button
						type="button"
						class="btn-close"
						@click="closeVectorSurvModal"
						aria-label="Close"
					></button>
				</div>
				<div class="modal-body">
					<form @submit.prevent="saveVectorSurvConfig">
						<div class="mb-3">
							<label for="vectorsurvToken" class="form-label">API Token</label>
							<input
								type="text"
								class="form-control"
								id="vectorsurvToken"
								v-model="vectorSurvForm.token"
							/>
							<div class="form-text">
								You can find this token in your VectorSurv account settings.
							</div>
						</div>
						<div class="mb-3 form-check">
							<input
								type="checkbox"
								class="form-check-input"
								id="vectorsurvSyncCheck"
								v-model="vectorSurvForm.autoSync"
							/>
							<label class="form-check-label" for="vectorsurvSyncCheck">
								Enable automatic synchronization
							</label>
						</div>
						<div class="mb-3">
							<label for="vectorsurvSyncTime" class="form-label"
								>Sync Time</label
							>
							<input
								type="time"
								class="form-control"
								id="vectorsurvSyncTime"
								v-model="vectorSurvForm.syncTime"
							/>
						</div>
					</form>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click="closeVectorSurvModal"
					>
						Cancel
					</button>
					<button
						type="button"
						class="btn btn-success"
						@click="saveVectorSurvConfig"
					>
						Save Changes
					</button>
				</div>
			</div>
		</div>
	</div>

	<!-- VeeMac Edit Credentials Modal -->
	<div
		class="modal fade"
		ref="veemacModalEl"
		tabindex="-1"
		aria-labelledby="veemacModalLabel"
		aria-hidden="true"
	>
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="veemacModalLabel">
						Edit VeeMac Credentials
					</h5>
					<button
						type="button"
						class="btn-close"
						@click="closeVeeMacModal"
						aria-label="Close"
					></button>
				</div>
				<div class="modal-body">
					<form @submit.prevent="saveVeeMacConfig">
						<div class="mb-3">
							<label for="veemacUsername" class="form-label">Username</label>
							<input
								type="text"
								class="form-control"
								id="veemacUsername"
								v-model="veeMacForm.username"
							/>
						</div>
						<div class="mb-3">
							<label for="veemacPassword" class="form-label">Password</label>
							<input
								type="password"
								class="form-control"
								id="veemacPassword"
								v-model="veeMacForm.password"
							/>
						</div>
						<div class="mb-3 form-check">
							<input
								type="checkbox"
								class="form-check-input"
								id="veemacSyncCheck"
								v-model="veeMacForm.autoSync"
								@change="handleVeeMacSyncChange"
							/>
							<label class="form-check-label" for="veemacSyncCheck">
								Enable automatic synchronization
							</label>
						</div>
						<div class="mb-3">
							<label for="veemacSyncFrequency" class="form-label"
								>Sync Frequency</label
							>
							<select
								class="form-select"
								id="veemacSyncFrequency"
								v-model="veeMacForm.syncFrequency"
								:disabled="!veeMacForm.autoSync"
							>
								<option value="daily">Daily</option>
								<option value="weekly">Weekly</option>
								<option value="hourly">Hourly</option>
							</select>
						</div>
					</form>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click="closeVeeMacModal"
					>
						Cancel
					</button>
					<button
						type="button"
						class="btn btn-danger"
						@click="saveVeeMacConfig"
					>
						Save Changes
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from "vue";
import { Modal } from "bootstrap";

interface ArcGISConfig {
	invalidatedAt: string | null;
	accessTokenExpires: string;
}

interface VectorSurvConfig {
	maskedToken: string;
	lastSync: string;
}

interface VeeMacConfig {
	username: string;
	lastSync: string;
}

interface VectorSurvForm {
	token: string;
	autoSync: boolean;
	syncTime: string;
}

interface VeeMacForm {
	username: string;
	password: string;
	autoSync: boolean;
	syncFrequency: string;
}

interface URLs {
	oauthRefreshArcGIS: string;
	configurationArcGIS: string;
}

// Refs
const vectorsurvModalEl = ref<HTMLElement | null>(null);
const veemacModalEl = ref<HTMLElement | null>(null);

let vectorsurvModal: Modal | null = null;
let veemacModal: Modal | null = null;

// Data
const arcGISConfig = ref<ArcGISConfig | null>({
	invalidatedAt: null,
	accessTokenExpires: new Date(Date.now() + 86400000).toISOString(),
});

const urls = reactive<URLs>({
	oauthRefreshArcGIS: "/oauth/refresh/arcgis",
	configurationArcGIS: "/settings/configuration/arcgis",
});

const vectorSurvConfig = reactive<VectorSurvConfig>({
	maskedToken: "vs_9f72b5e3******************************c11d",
	lastSync: "December 5, 2025 at 08:34 AM (2 days ago)",
});

const veeMacConfig = reactive<VeeMacConfig>({
	username: "mosquito_district21",
	lastSync: "December 6, 2025 at 11:15 PM (Yesterday)",
});

const vectorSurvForm = reactive<VectorSurvForm>({
	token: "vs_9f72b5e3c8a1d492f6b7e54321098c11d",
	autoSync: true,
	syncTime: "02:00",
});

const veeMacForm = reactive<VeeMacForm>({
	username: "mosquito_district21",
	password: "password123",
	autoSync: false,
	syncFrequency: "daily",
});

// Lifecycle hooks
onMounted(() => {
	if (vectorsurvModalEl.value) {
		vectorsurvModal = new Modal(vectorsurvModalEl.value);
	}
	if (veemacModalEl.value) {
		veemacModal = new Modal(veemacModalEl.value);
	}

	// Load data from API
	loadIntegrations();
});

onUnmounted(() => {
	vectorsurvModal?.dispose();
	veemacModal?.dispose();
});

// Methods
const loadIntegrations = async (): Promise<void> => {
	// TODO: Fetch data from API
	// Example:
	// const response = await fetch('/api/integrations');
	// const data = await response.json();
	// Update reactive data with fetched data
};

const isTokenExpired = (expiresAt: string): boolean => {
	return new Date(expiresAt) < new Date();
};

const formatRelativeTime = (dateString: string): string => {
	const date = new Date(dateString);
	const now = new Date();
	const diffInSeconds = Math.floor((date.getTime() - now.getTime()) / 1000);

	if (diffInSeconds < 0) {
		return "Expired";
	}

	const days = Math.floor(diffInSeconds / 86400);
	const hours = Math.floor((diffInSeconds % 86400) / 3600);
	const minutes = Math.floor((diffInSeconds % 3600) / 60);

	if (days > 0) {
		return `in ${days} day${days > 1 ? "s" : ""}`;
	} else if (hours > 0) {
		return `in ${hours} hour${hours > 1 ? "s" : ""}`;
	} else {
		return `in ${minutes} minute${minutes > 1 ? "s" : ""}`;
	}
};

const openVectorSurvModal = (): void => {
	vectorsurvModal?.show();
};

const closeVectorSurvModal = (): void => {
	vectorsurvModal?.hide();
};

const openVeeMacModal = (): void => {
	veemacModal?.show();
};

const closeVeeMacModal = (): void => {
	veemacModal?.hide();
};

const saveVectorSurvConfig = async (): Promise<void> => {
	try {
		// TODO: Send data to API
		// await fetch('/api/integrations/vectorsurv', {
		//   method: 'PUT',
		//   headers: { 'Content-Type': 'application/json' },
		//   body: JSON.stringify(vectorSurvForm),
		// });

		console.log("Saving VectorSurv config:", vectorSurvForm);
		closeVectorSurvModal();
		await loadIntegrations();
	} catch (error) {
		console.error("Error saving VectorSurv config:", error);
	}
};

const saveVeeMacConfig = async (): Promise<void> => {
	try {
		// TODO: Send data to API
		// await fetch('/api/integrations/veemac', {
		//   method: 'PUT',
		//   headers: { 'Content-Type': 'application/json' },
		//   body: JSON.stringify(veeMacForm),
		// });

		console.log("Saving VeeMac config:", veeMacForm);
		closeVeeMacModal();
		await loadIntegrations();
	} catch (error) {
		console.error("Error saving VeeMac config:", error);
	}
};

const handleVeeMacSyncChange = (): void => {
	if (!veeMacForm.autoSync) {
		veeMacForm.syncFrequency = "daily";
	}
};

const removeIntegration = async (integration: string): Promise<void> => {
	if (
		!confirm(`Are you sure you want to remove the ${integration} integration?`)
	) {
		return;
	}

	try {
		// TODO: Send delete request to API
		// await fetch(`/api/integrations/${integration}`, {
		//   method: 'DELETE',
		// });

		console.log(`Removing ${integration} integration`);
		await loadIntegrations();
	} catch (error) {
		console.error(`Error removing ${integration} integration:`, error);
	}
};
</script>

<style scoped>
.integration-card {
	border-left: 5px solid #0d6efd;
	transition: all 0.2s;
}

.integration-card:hover {
	box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15);
}

.integration-card.fieldseeker {
	border-left-color: #0d6efd;
}

.integration-card.vectorsurv {
	border-left-color: #198754;
}

.integration-card.veemac {
	border-left-color: #dc3545;
}

.status-active {
	color: #198754;
}

.status-inactive {
	color: #dc3545;
}

.token-display {
	font-family: monospace;
	background-color: #f8f9fa;
	padding: 0.25rem 0.5rem;
	border-radius: 0.25rem;
	border: 1px solid #dee2e6;
}
</style>
