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
							v-model="poolLocation.longitude"
							:class="{
								'border-warning':
									poolLocation.longitude !==
									selectedTask.pool?.location.longitude,
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
							v-model="poolLocation.latitude"
							:class="{
								'border-warning':
									poolLocation.latitude !==
									selectedTask.pool?.location?.latitude,
							}"
						/>
					</div>
				</div>

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Pool Condition:</label>
					<div class="col-sm-9">
						<select
							class="form-select"
							v-model="poolCondition"
							:class="{
								'border-warning':
									poolCondition !== selectedTask.pool?.condition,
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

				<div class="row mb-3">
					<label class="col-sm-3 col-form-label fw-bold">Owner Contact:</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="siteOwner.name"
							:class="{
								'border-warning':
									siteOwner.name !== selectedTask.pool?.site.owner?.name,
							}"
						/>
					</div>
				</div>

				<div class="row mb-4">
					<label class="col-sm-3 col-form-label fw-bold">
						Resident Contact:
					</label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							v-model="siteResident.name"
							:class="{
								'border-warning':
									siteResident.name !== selectedTask.pool?.site.resident?.name,
							}"
						/>
					</div>
				</div>
			</form>
		</div>

		<!-- Map Components -->
		<div class="map-container" v-if="session.organization">
			<MapMultipoint
				ref="mapMultipoint"
				id="map"
				:bounds="mapBounds"
				:markers="mapMarkers"
				:organizationId="session.organization.id"
				:tegola="session.urls?.tegola ?? ''"
			></MapMultipoint>
		</div>
		<div v-else>
			<p>loading...</p>
		</div>

		<div class="map-container" v-if="session.organization && selectedTask.pool">
			<MapProxiedArcgisTile
				:location="selectedTask.pool?.location"
				:markers="[]"
				:organizationId="session.organization.id"
				:tegola="session.urls?.tegola ?? ''"
				:urlTiles="session.urls?.tile ?? ''"
				@map-click="doPoolLocation"
			></MapProxiedArcgisTile>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import MapMultipoint from "@/components/MapMultipoint.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { formatAddress } from "@/format";
import { useSessionStore } from "@/store/session";
import {
	Bounds,
	Contact,
	MapClickEvent,
	Marker,
	Pool,
	ReviewTask,
	User,
} from "@/types";
import type { Location } from "@/type/api";

interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapMarkers: Marker[];
	newPoolCondition: string;
	newPoolLocation: Location;
	selectedTask?: ReviewTask;
}
const props = defineProps<Props>();
const poolCondition = ref<string>("unknown");
const poolLocation = ref<Location>({
	latitude: 0,
	longitude: 0,
});
const siteOwner = ref<Contact>({
	has_email: false,
	has_phone: false,
	name: "",
});
const siteResident = ref<Contact>({
	has_email: false,
	has_phone: false,
	name: "",
});
const session = useSessionStore();
function doPoolLocation(event: MapClickEvent) {
	console.log("pool location", event);
}
</script>
