<template>
	<div class="container py-4">
		<div class="row">
			<div class="col-12">
				<div class="d-flex justify-content-between align-items-center mb-4">
					<h1>
						<i class="bi bi-geo-alt-fill text-primary me-2"></i> District
						Settings
					</h1>
					<button class="btn btn-primary" @click="saveChanges">
						<i class="bi bi-save me-2"></i>Save Changes
					</button>
				</div>

				<MapServiceArea
					:organization-id="organization.id"
					:tegola="tegolaUrl"
					:xmin="organization.serviceArea.min.x"
					:ymin="organization.serviceArea.min.y"
					:xmax="organization.serviceArea.max.x"
					:ymax="organization.serviceArea.max.y"
				/>

				<div class="row">
					<!-- Basic Information -->
					<div class="col-md-6">
						<div class="card settings-card">
							<div class="card-header bg-light">
								<h5>
									<i class="bi bi-building me-2"></i> Organization Information
								</h5>
							</div>
							<div class="card-body">
								<div class="mb-3">
									<label for="agencyName" class="form-label">
										<i class="bi bi-briefcase me-1"></i> Agency Name
									</label>
									<input
										type="text"
										class="form-control"
										id="name"
										v-model="organization.name"
									/>
								</div>
								<div class="mb-3">
									<label for="website" class="form-label">
										<i class="bi bi-globe me-1"></i> Website
									</label>
									<input
										type="url"
										class="form-control"
										id="website"
										v-model="organization.website"
									/>
								</div>
								<div class="mb-3">
									<label for="generalManager" class="form-label">
										<i class="bi bi-person-badge me-1"></i> General Manager Name
									</label>
									<input
										type="text"
										class="form-control"
										id="generalManager"
										v-model="organization.generalManagerName"
									/>
								</div>
							</div>
						</div>
					</div>

					<!-- Contact Information -->
					<div class="col-md-6">
						<div class="card settings-card">
							<div class="card-header bg-light">
								<h5>
									<i class="bi bi-telephone me-2"></i> Contact Information
								</h5>
							</div>
							<div class="card-body">
								<div class="mb-3">
									<label for="address" class="form-label">
										<i class="bi bi-geo me-1"></i> Street
									</label>
									<input
										type="text"
										class="form-control"
										id="address"
										v-model="organization.officeAddressStreet"
									/>
								</div>
								<div class="row">
									<div class="col-md-6 mb-3">
										<label for="city" class="form-label">
											<i class="bi bi-building me-1"></i> City
										</label>
										<input
											type="text"
											class="form-control"
											id="city"
											v-model="organization.officeAddressCity"
										/>
									</div>
									<div class="col-md-6 mb-3">
										<label for="postalCode" class="form-label">
											<i class="bi bi-mailbox me-1"></i> Postal Code
										</label>
										<input
											type="text"
											class="form-control"
											id="postalCode"
											v-model="organization.officeAddressPostalCode"
										/>
									</div>
								</div>
								<div class="mb-3">
									<label for="phoneNumber" class="form-label">
										<i class="bi bi-telephone me-1"></i> Phone Number
									</label>
									<input
										type="tel"
										class="form-control"
										id="phoneNumber"
										v-model="organization.officePhone"
									/>
								</div>
								<div class="mb-3">
									<label for="faxNumber" class="form-label">
										<i class="bi bi-printer me-1"></i> Fax Number
									</label>
									<input
										type="tel"
										class="form-control"
										id="faxNumber"
										v-model="organization.officeFax"
									/>
								</div>
							</div>
						</div>
					</div>

					<!-- Organization Coverage Information -->
					<div class="col-12">
						<div class="card settings-card">
							<div class="card-header bg-light">
								<h5><i class="bi bi-map me-2"></i> Service Area Coverage</h5>
							</div>
							<div class="card-body">
								<div class="row">
									<div class="col-md-6 mb-3">
										<label for="totalArea" class="form-label">
											<i class="bi bi-rulers me-1"></i> Total Area (square
											meters)
										</label>
										<input
											type="number"
											class="form-control"
											id="totalArea"
											v-model.number="organization.serviceAreaSquareMeters"
										/>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import MapServiceArea from "../../components/MapServiceArea.vue";

interface ServiceAreaBounds {
	min: {
		x: number;
		y: number;
	};
	max: {
		x: number;
		y: number;
	};
}

interface Organization {
	id: string;
	name: string;
	website: string;
	generalManagerName: string;
	officeAddressStreet: string;
	officeAddressCity: string;
	officeAddressPostalCode: string;
	officePhone: string;
	officeFax: string;
	serviceArea: ServiceAreaBounds;
	serviceAreaSquareMeters: number | null;
}

const organization = ref<Organization>({
	id: "",
	name: "",
	website: "",
	generalManagerName: "",
	officeAddressStreet: "",
	officeAddressCity: "",
	officeAddressPostalCode: "",
	officePhone: "",
	officeFax: "",
	serviceArea: {
		min: { x: 0, y: 0 },
		max: { x: 0, y: 0 },
	},
	serviceAreaSquareMeters: null,
});

const tegolaUrl = ref<string>("");

const fetchOrganizationData = async (): Promise<void> => {
	try {
		// Replace with your actual API endpoint
		const response = await fetch("/api/organization/settings");
		if (!response.ok) {
			throw new Error("Failed to fetch organization data");
		}
		const data = await response.json();
		organization.value = {
			id: data.id,
			name: data.name,
			website: data.website || "",
			generalManagerName: data.generalManagerName || "",
			officeAddressStreet: data.officeAddressStreet || "",
			officeAddressCity: data.officeAddressCity || "",
			officeAddressPostalCode: data.officeAddressPostalCode || "",
			officePhone: data.officePhone || "",
			officeFax: data.officeFax || "",
			serviceArea: data.serviceArea || {
				min: { x: 0, y: 0 },
				max: { x: 0, y: 0 },
			},
			serviceAreaSquareMeters: data.serviceAreaSquareMeters || null,
		};
		tegolaUrl.value = data.tegolaUrl || "";
	} catch (error) {
		console.error("Error fetching organization data:", error);
	}
};

const saveChanges = async (): Promise<void> => {
	try {
		// Replace with your actual API endpoint
		const response = await fetch("/api/organization/settings", {
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(organization.value),
		});

		if (!response.ok) {
			throw new Error("Failed to save changes");
		}

		// Show success notification
		alert("Changes saved successfully!");
	} catch (error) {
		console.error("Error saving changes:", error);
		alert("Failed to save changes. Please try again.");
	}
};

onMounted(() => {
	fetchOrganizationData();
});
</script>

<style scoped>
.settings-card {
	box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
	margin-bottom: 30px;
}
</style>
