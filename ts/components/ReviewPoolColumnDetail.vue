<style scoped>
.map-container {
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
	height: 500px;
	margin-bottom: 20px;
	margin-top: 20px;
	align-items: center;
	justify-content: center;
	/* Prevent touch scrolling issues */
	touch-action: pan-y pinch-zoom;
}
#map {
	width: 100%;
	height: 100%;
}
</style>
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
									selectedTask.pool?.location?.longitude,
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
			<MapLocator :markers="markers" v-model="mapCamera"></MapLocator>
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
import { computed, ref } from "vue";
import MapLocator from "@/components/MapLocator.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { formatAddress } from "@/format";
import { useSessionStore } from "@/store/session";
import type { MapClickEvent, Marker } from "@/types";
import { Bounds, Contact, Pool, ReviewTask, User } from "@/type/api";
import type { Location } from "@/type/api";
import { Camera } from "@/type/map";

interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapMarkers: Marker[];
	newPoolCondition: string;
	newPoolLocation: Location;
	selectedTask?: ReviewTask;
}
const mapCamera = ref<Camera>(new Camera());
const props = defineProps<Props>();
const poolCondition = ref<string>("unknown");
const poolLocation = ref<Location>({
	latitude: 0,
	longitude: 0,
});
const siteOwner = ref<Contact>(new Contact());
const siteResident = ref<Contact>(new Contact());
const session = useSessionStore();
const markers = computed((): Marker[] => {
	if (!poolLocation.value) {
		return [];
	}
	if (
		poolLocation.value.latitude == 0.0 &&
		poolLocation.value.longitude == 0.0
	) {
		return [];
	}
	const marker = {
		color: "#FF0000",
		draggable: false,
		id: "x",
		location: poolLocation.value,
	};
	return [marker];
});
function doPoolLocation(event: MapClickEvent) {
	console.log("pool location", event);
}
</script>
