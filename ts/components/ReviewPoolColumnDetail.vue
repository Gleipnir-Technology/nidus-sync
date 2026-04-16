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
</style>
<template>
	<div class="mb-4">
		<!-- No Selection State -->
		<div
			v-show="!selectedTask"
			class="h-100 align-items-center justify-content-center text-muted"
		>
			<div class="text-center">
				<i class="bi bi-cursor-fill" style="font-size: 48px"></i>
				<p class="mt-2">Select an entry from the list to review</p>
			</div>
		</div>
		<div v-show="selectedTask !== undefined">
			<h4 class="mb-3">Entry #{{ selectedTask?.id ?? "" }} Details</h4>

			<div class="row mb-3">
				<label class="col-sm-3 col-form-label fw-bold">Address:</label>
				<div class="col-sm-9">
					<input
						type="text"
						class="form-control"
						:value="modelValue.address"
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
						v-model="modelValue.location.longitude"
						:class="{
							'border-warning':
								modelValue.location.longitude !==
								selectedTask?.pool?.location?.longitude,
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
						v-model="modelValue.location.latitude"
						:class="{
							'border-warning':
								modelValue.location.latitude !==
								selectedTask?.pool?.location?.latitude,
						}"
					/>
				</div>
			</div>

			<div class="row mb-3">
				<label class="col-sm-3 col-form-label fw-bold">Pool Condition:</label>
				<div class="col-sm-9">
					<select
						class="form-select"
						v-model="modelValue.condition"
						:class="{
							'border-warning':
								modelValue.condition !== selectedTask?.pool?.condition,
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
						v-model="modelValue.owner"
						:class="{
							'border-warning':
								modelValue.owner !==
								(selectedTask?.pool?.site.owner?.name ?? ''),
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
								siteResident.name !==
								(selectedTask?.pool?.site.resident?.name ?? ''),
						}"
					/>
				</div>
			</div>
		</div>
	</div>

	<!-- Map Components -->
	<div class="map-container" v-if="session.organization">
		<MapLocator
			@click="doPoolLocation"
			:markers="mapMarkers"
			:useSatellite="true"
			v-model="mapCamera"
		></MapLocator>
	</div>
	<div v-else>
		<p>loading...</p>
	</div>

	<div class="map-container" v-if="session.organization && session.urls">
		<MapProxiedArcgisTile
			@click="doPoolLocation"
			:markers="mapMarkers"
			:organizationId="session.organization.id"
			:tegola="session.urls!.tegola"
			:urlTiles="session.urls!.tile"
			v-model="_mapFlyoverCamera"
		></MapProxiedArcgisTile>
	</div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from "vue";
import MapLocator from "@/components/MapLocator.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { useSessionStore } from "@/store/session";
import type { MapClickEvent, Marker } from "@/types";
import { Bounds, Contact, Pool, ReviewTask, User } from "@/type/api";
import type { Location } from "@/type/api";
import { Camera } from "@/type/map";

interface Emits {
	(e: "update:modelValue", value: ReviewTaskPoolForm): void;
}
export interface ReviewTaskPoolForm {
	address: string;
	condition: string;
	location: Location;
	owner: string;
	resident: string;
}
interface Props {
	loading: boolean;
	mapBounds?: Bounds;
	mapFlyoverCamera: Camera;
	mapMarkers: Marker[];
	modelValue: ReviewTaskPoolForm;
	selectedTask: ReviewTask | undefined;
}
const emit = defineEmits<Emits>();
const mapCamera = ref<Camera>(new Camera());
const _mapFlyoverCamera = ref<Camera>(new Camera());
const props = defineProps<Props>();
const siteResident = ref<Contact>(new Contact());
const session = useSessionStore();
function doPoolLocation(event: MapClickEvent) {
	emit("update:modelValue", {
		address: props.modelValue.address,
		condition: props.modelValue.condition,
		location: event.location,
		owner: props.modelValue.owner,
		resident: props.modelValue.resident,
	});
}
watch(
	() => props.mapFlyoverCamera,
	(newMapFlyoverCamera: Camera) => {
		console.log("map flyover camera update", newMapFlyoverCamera);
		_mapFlyoverCamera.value = newMapFlyoverCamera;
	},
);
</script>
