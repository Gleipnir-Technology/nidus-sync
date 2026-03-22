<template>
	<div class="container py-5">
		<div class="row justify-content-center">
			<div class="col-lg-8">
				<!-- Header -->
				<div class="d-flex align-items-center mb-4">
					<i class="bi bi-globe2 text-primary fs-1 me-3"></i>
					<div>
						<h1 class="mb-0">ArcGIS Integration</h1>
						<p class="text-muted mb-0">Configure your Esri ArcGIS connection</p>
					</div>
				</div>

				<!-- Main Card -->
				<div class="card shadow-sm">
					<div class="card-body p-4">
						<form @submit.prevent="handleSubmit">
							<!-- OAuth Authentication Section -->
							<div class="mb-4">
								<h5 class="card-title border-bottom pb-2 mb-3">
									<i class="bi bi-key-fill text-success me-2"></i>OAuth
									Authentication
								</h5>

								<div class="row g-3">
									<div class="col-md-6">
										<label class="form-label fw-semibold">ArcGIS ID</label>
										<input
											type="text"
											class="form-control"
											:value="arcGISAccount.id"
											readonly
										/>
									</div>

									<div class="col-md-6">
										<label class="form-label fw-semibold"
											>Organization Name</label
										>
										<input
											type="text"
											class="form-control"
											:value="arcGISAccount.name"
											readonly
										/>
									</div>

									<div class="col-md-12">
										<label class="form-label fw-semibold">Authorized By</label>
										<input
											type="text"
											class="form-control"
											:value="arcGISOAuth.username"
											readonly
										/>
									</div>

									<div class="col-md-6">
										<label class="form-label fw-semibold">Token Age</label>
										<input
											type="text"
											class="form-control"
											:value="formatTimeRelative(arcGISOAuth.created)"
											readonly
										/>
									</div>

									<div class="col-md-6">
										<label class="form-label fw-semibold"
											>Token Expiration</label
										>
										<input
											type="text"
											class="form-control"
											:value="
												formatTimeRelative(arcGISOAuth.refreshTokenExpires)
											"
											readonly
										/>
									</div>
								</div>

								<!-- Token Actions -->
								<div class="mt-3 d-flex gap-2">
									<a
										class="btn btn-outline-primary"
										:href="urls.oAuthRefreshArcGIS"
									>
										<i class="bi bi-arrow-clockwise me-1"></i>Refresh Token
									</a>
									<button
										type="button"
										class="btn btn-outline-danger"
										@click="showDeleteModal"
									>
										<i class="bi bi-trash me-1"></i>Delete Token
									</button>
								</div>
							</div>

							<hr class="my-4" />

							<!-- Feature Layers Section -->
							<div class="mb-4">
								<h5 class="card-title border-bottom pb-2 mb-3">
									<i class="bi bi-layers-fill text-info me-2"></i>Feature Layer
									Configuration
								</h5>

								<div class="row g-3">
									<div class="col-md-12">
										<label for="map-service" class="form-label fw-semibold">
											Map Service (Aerial Imagery)
											<span class="text-danger">*</span>
										</label>
										<select
											class="form-select"
											id="map-service"
											v-model="selectedMapService"
											required
										>
											<option
												v-for="service in serviceMaps"
												:key="service.arcgisId"
												:value="service.arcgisId"
											>
												{{ service.name }}
											</option>
										</select>
										<div class="form-text">
											Select the feature layer for aerial imagery data
										</div>
									</div>
								</div>
							</div>

							<!-- Save Button -->
							<div class="d-grid gap-2 d-md-flex justify-content-md-end mt-4">
								<button
									type="button"
									class="btn btn-secondary me-md-2"
									@click="handleCancel"
								>
									Cancel
								</button>
								<button type="submit" class="btn btn-primary">
									<i class="bi bi-save me-1"></i>Save Configuration
								</button>
							</div>
						</form>
					</div>
				</div>

				<!-- Info Alert -->
				<div
					class="alert alert-info mt-3 d-flex align-items-start"
					role="alert"
				>
					<i class="bi bi-info-circle-fill me-2 flex-shrink-0"></i>
					<div>
						<strong>Note:</strong> Changes to feature layer selections will take
						effect immediately after saving. Refreshing the OAuth token will
						require re-authentication with your ArcGIS account.
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Delete Confirmation Modal -->
	<div
		class="modal fade"
		ref="deleteModalElement"
		tabindex="-1"
		aria-labelledby="deleteModalLabel"
		aria-hidden="true"
	>
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header bg-danger text-white">
					<h5 class="modal-title" id="deleteModalLabel">
						<i class="bi bi-exclamation-triangle-fill me-2"></i>Confirm Delete
					</h5>
					<button
						type="button"
						class="btn-close btn-close-white"
						@click="hideDeleteModal"
						aria-label="Close"
					></button>
				</div>
				<div class="modal-body">
					<p class="mb-2">
						Are you sure you want to delete the OAuth token and disable the
						ArcGIS integration?
					</p>
					<p class="text-muted mb-0">
						<strong>This action cannot be undone.</strong> You will need to
						re-authenticate to restore the integration.
					</p>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click="hideDeleteModal"
					>
						Cancel
					</button>
					<button type="button" class="btn btn-danger" @click="handleDelete">
						<i class="bi bi-trash me-1"></i>Delete Token
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { Ref } from "vue";

// Interfaces
interface ArcGISAccount {
	id: string;
	name: string;
}

interface ArcGISOAuth {
	username: string;
	created: string;
	refreshTokenExpires: string;
}

interface ServiceMap {
	arcgisId: string;
	name: string;
}

interface URLs {
	configurationArcGIS: string;
	oAuthRefreshArcGIS: string;
}

// Props (if data comes from parent)
interface Props {
	initialAccount?: ArcGISAccount;
	initialOAuth?: ArcGISOAuth;
	initialServiceMaps?: ServiceMap[];
	initialUrls?: URLs;
}

const props = withDefaults(defineProps<Props>(), {
	initialAccount: () => ({ id: "", name: "" }),
	initialOAuth: () => ({ username: "", created: "", refreshTokenExpires: "" }),
	initialServiceMaps: () => [],
	initialUrls: () => ({
		configurationArcGIS: "/settings/integrations/arcgis",
		oAuthRefreshArcGIS: "/oauth/refresh/arcgis",
	}),
});

// Reactive state
const arcGISAccount: Ref<ArcGISAccount> = ref({ ...props.initialAccount });
const arcGISOAuth: Ref<ArcGISOAuth> = ref({ ...props.initialOAuth });
const serviceMaps: Ref<ServiceMap[]> = ref([...props.initialServiceMaps]);
const urls: Ref<URLs> = ref({ ...props.initialUrls });
const selectedMapService = ref<string>("");

// Modal reference
const deleteModalElement = ref<HTMLElement | null>(null);
let deleteModal: any = null;

// Lifecycle hooks
onMounted(() => {
	// Initialize Bootstrap modal
	if (
		deleteModalElement.value &&
		typeof window !== "undefined" &&
		(window as any).bootstrap
	) {
		deleteModal = new (window as any).bootstrap.Modal(deleteModalElement.value);
	}

	// Set initial selected service if available
	if (serviceMaps.value.length > 0) {
		selectedMapService.value = serviceMaps.value[0].arcgisId;
	}

	// Load data from API
	fetchData();
});

onUnmounted(() => {
	// Clean up modal
	if (deleteModal) {
		deleteModal.dispose();
	}
});

// Methods
const fetchData = async () => {
	try {
		// Replace with your actual API endpoint
		const response = await fetch("/api/settings/arcgis");
		if (response.ok) {
			const data = await response.json();
			arcGISAccount.value = data.account;
			arcGISOAuth.value = data.oauth;
			serviceMaps.value = data.serviceMaps;
			urls.value = data.urls;

			if (data.selectedMapService) {
				selectedMapService.value = data.selectedMapService;
			}
		}
	} catch (error) {
		console.error("Error fetching ArcGIS data:", error);
	}
};

const handleSubmit = async () => {
	try {
		const response = await fetch(urls.value.configurationArcGIS, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				"map-service": selectedMapService.value,
			}),
		});

		if (response.ok) {
			// Handle success (show notification, redirect, etc.)
			console.log("Configuration saved successfully");
		} else {
			// Handle error
			console.error("Failed to save configuration");
		}
	} catch (error) {
		console.error("Error saving configuration:", error);
	}
};

const handleCancel = () => {
	// Reset or navigate back
	window.history.back();
};

const showDeleteModal = () => {
	if (deleteModal) {
		deleteModal.show();
	}
};

const hideDeleteModal = () => {
	if (deleteModal) {
		deleteModal.hide();
	}
};

const handleDelete = async () => {
	try {
		const response = await fetch("/api/oauth/arcgis", {
			method: "DELETE",
		});

		if (response.ok) {
			hideDeleteModal();
			// Handle success (redirect, show notification, etc.)
			console.log("Token deleted successfully");
		} else {
			console.error("Failed to delete token");
		}
	} catch (error) {
		console.error("Error deleting token:", error);
	}
};

const formatTimeRelative = (dateString: string): string => {
	if (!dateString) return "";

	const date = new Date(dateString);
	const now = new Date();
	const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

	if (diffInSeconds < 60) {
		return "just now";
	} else if (diffInSeconds < 3600) {
		const minutes = Math.floor(diffInSeconds / 60);
		return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
	} else if (diffInSeconds < 86400) {
		const hours = Math.floor(diffInSeconds / 3600);
		return `${hours} hour${hours > 1 ? "s" : ""} ago`;
	} else {
		const days = Math.floor(diffInSeconds / 86400);
		return `${days} day${days > 1 ? "s" : ""} ago`;
	}
};
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
