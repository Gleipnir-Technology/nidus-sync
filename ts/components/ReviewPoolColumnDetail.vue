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
							v-model="form.longitude"
							:class="{
								'border-warning': form.longitude !== originalValues.longitude,
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
							v-model="form.latitude"
							:class="{
								'border-warning': form.latitude !== originalValues.latitude,
							}"
						/>
					</div>
				</div>

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Pool Condition:</label>
					<div class="col-sm-9">
						<select
							class="form-select"
							v-model="form.condition"
							:class="{
								'border-warning': form.condition !== originalValues.condition,
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

				<div v-if="form.ownerContact" class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Owner Contact:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="form.ownerContact"
							:class="{
								'border-warning':
									form.ownerContact !== originalValues.ownerContact,
							}"
						/>
					</div>
				</div>

				<div v-if="form.residentContact" class="row mb-4">
					<label class="col-sm-3 col-form-label fw-bold">
						Resident Contact:
					</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="form.residentContact"
							:class="{
								'border-warning':
									form.residentContact !== originalValues.residentContact,
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
				:organization-id="organizationId"
				:tegola="tegolaUrl"
				:xmin="serviceArea.xmin"
				:ymin="serviceArea.ymin"
				:xmax="serviceArea.xmax"
				:ymax="serviceArea.ymax"
			></MapMultipoint>
		</div>

		<div class="map-container">
			<MapProxiedArcgisTile
				ref="mapTile"
				class="map"
				:organization-id="organizationId"
				:tegola="tegolaUrl"
				:tiles-url="tilesUrl"
				:latitude="selectedTask.location.latitude"
				:longitude="selectedTask.location.longitude"
				@map-click="updatePoolLocation"
			></MapProxiedArcgisTile>
		</div>
	</div>
</template>
<script setup lang="ts">
import MapMultipoint from "@/components/MapMultipoint.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import ReviewTask from "@/types";

interface Props {
	selectedTask?: ReviewTask;
}
const props = defineProps<Props>();
</script>
