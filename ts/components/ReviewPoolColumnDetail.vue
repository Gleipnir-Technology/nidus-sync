<template>
	<!-- No Selection State -->
	<div
		v-if="!selectedTask"
		class="h-100 d-flex align-items-center justify-content-center text-muted"
	>
		<div class="text-center">
			<i class="bi bi-cursor-fill" style="font-size: 48px"></i>
			<p class="mt-2">Select an entry from the list to review</p>
		</div>
	</div>

	<!-- Selected Task Details -->
	<div v-else>
		<div class="mb-4">
			<h4 class="mb-3">Entry #{{ selectedTask.id }} Details</h4>

			<form @submit.prevent>
				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Address:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							:value="formatAddress(selectedTask.address)"
							readonly
						/>
					</div>
				</div>

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Longitude:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="selectedTaskChanges.location.longitude"
							:class="{
								'border-warning':
									selectedTaskChanges.location.longitude !==
									selectedTask.location.longitude,
							}"
						/>
					</div>
				</div>

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Latitude:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="selectedTaskChanges.location.latitude"
							:class="{
								'border-warning':
									selectedTaskChanges.location?.latitude !==
									selectedTask.location?.latitude,
							}"
						/>
					</div>
				</div>

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Pool Condition:</label>
					<div class="col-sm-9">
						<select
							class="form-select"
							v-model="selectedTaskChanges.pool.condition"
							:class="{
								'border-warning':
									selectedTaskChanges.pool.condition !==
									selectedTask.pool.condition,
							}"
						>
							<option value="">-- Select --</option>
							<option value="blue">Blue</option>
							<option value="dry">Dry</option>
							<option value="false pool">False Pool</option>
							<option value="unknown">Unknown</option>
							<option value="green">Green</option>
							<option value="murky">Murky</option>
						</select>
					</div>
				</div>

				<div v-if="selectedTaskChanges.pool.ownerContact" class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Owner Contact:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="selectedTaskChanges.pool.owner_contact"
							:class="{
								'border-warning':
									selectedTaskChanges.pool.owner_contact !==
									selectedTask.pool.owner_contact,
							}"
						/>
					</div>
				</div>

				<div v-if="selectedTaskChanges.pool.resident_contact" class="row mb-4">
					<label class="col-sm-3 col-form-label fw-bold">
						Resident Contact:
					</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="selectedTaskChanges.pool.resident_contact"
							:class="{
								'border-warning':
									selectedTaskChanges.pool.resident_contact !==
									selectedTask.pool.resident_contact,
							}"
						/>
					</div>
				</div>
			</form>
		</div>

		<!-- Map Components -->
		<div class="map-container">
			<MapMultipoint
				ref="mapMultipoint"
				id="map"
				:bounds="mapBounds"
				:markers="mapMarkers"
				:organizationId="user.organization.id"
				:tegola="user.urls.tegola"
				:xmin="user.organization.service_area?.min.x ?? 0"
				:ymin="user.organization.service_area?.min.y ?? 0"
				:xmax="user.organization.service_area?.max.x ?? 0"
				:ymax="user.organization.service_area?.max.y ?? 0"
			></MapMultipoint>
		</div>

		<div class="map-container">
			<MapProxiedArcgisTile
				ref="mapTile"
				class="map"
				:location="selectedTask.location"
				:organization-id="user.organization.id"
				:tegola="user.urls.tegola"
				:urlTiles="user.urls.tile"
				@map-click="doPoolLocation"
			></MapProxiedArcgisTile>
		</div>
	</div>
</template>
<script setup lang="ts">
import MapMultipoint from "@/components/MapMultipoint.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { formatAddress } from "@/format";
import ReviewTask from "@/types";

interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapMarkers: Marker[];
	selectedTaskChanges: ReviewTask;
	selectedTask?: ReviewTask;
	user: User | null;
}
const props = defineProps<Props>();
function doPoolLocation(lat, lng) {
	console.log("pool location", lat, lng);
}
</script>
