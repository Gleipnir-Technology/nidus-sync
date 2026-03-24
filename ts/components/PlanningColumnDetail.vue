<style scoped>
.map-container {
	height: 400px;
	margin-bottom: 1rem;
}
</style>
<template>
	<div class="card shadow-sm mb-3">
		<div class="card-header bg-white pane-header">
			Active Investigation Workbench
		</div>
		<div class="card-body">
			<div class="map-container">
				<MapMultipoint
					id="map"
					:markers="markers"
					:organization-id="user.organization.id"
					:tegola="user.urls.tegola"
					:xmin="user?.organization.service_area?.xmin ?? 0"
					:ymin="user?.organization.service_area?.ymin ?? 0"
					:xmax="user?.organization.service_area?.xmax ?? 0"
					:ymax="user?.organization.service_area?.ymax ?? 0"
				></MapMultipoint>
			</div>

			<div class="row g-3">
				<div class="col-md-12">
					<div class="card border">
						<div class="card-body">
							<div class="fw-semibold">Selected Signals</div>
							<div class="text-muted small">
								{{ selectedSignals.length }} Signal{{
									selectedSignals.length !== 1 ? "s" : ""
								}}
								Selected
							</div>

							<div
								v-if="selectedSignals.length === 0"
								class="text-muted small mt-2 fst-italic"
							>
								Click signals from the left panel to select them
							</div>

							<div
								class="mt-2"
								v-show="selectedSignals.length > 0"
							>
									<div v-for="signal in selectedSignals" :key="signal.id">
										<PlanningColumnDetailEntry :signal="signal"/>
									</div>
							</div>

							<button
								v-show="selectedSignals.length > 0"
								@click="clearSelection()"
								class="btn btn-sm btn-outline-secondary mt-2 w-100"
							>
								Clear Selection
							</button>
						</div>
					</div>
				</div>
			</div>

			<div v-show="showMapTile" class="map-container">
				<MapProxiedArcgisTile
					ref="mapTile"
					class="map"
					:organization-id="user.organization.id"
					:tegola="user.urls.tegola"
					:url-tiles="user.urls.tile"
					:latitude="selectedSignalLocation()?.latitude ?? 0.0"
					:longitude="selectedSignalLocation()?.longitude ?? 0.0"
					@map-click="updateSignalLocation"
				>
				</MapProxiedArcgisTile>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import MapMultipoint from "./MapMultipoint.vue";
import MapProxiedArcgisTile from "./MapProxiedArcgisTile.vue";
import { shortAddress } from "../format";
import TimeRelative from "./TimeRelative.vue";
import PlanningColumnDetailEntry from "./PlanningColumnDetailEntry.vue";
import { useUserStore } from "../store/user";

interface Props {
	markers: Marker[];
	selectedSignals: Array<Signal>;
}
const props = defineProps<Props>();
const user = useUserStore();
const configureMapTile = () => {
	if (!mapTile.value) return;

	mapTile.value.on("load", () => {
		mapTile.value.addLayer({
			id: "parcel",
			minzoom: 14,
			paint: {
				"line-color": "#0f0",
			},
			source: "tegola",
			"source-layer": "parcel",
			type: "line",
		});
		mapTile.value.addLayer({
			id: "signal-point",
			paint: {
				"circle-color": "#0D6EfD",
				"circle-radius": 7,
				"circle-stroke-width": 2,
				"circle-stroke-color": "#024AB6",
			},
			source: "tegola",
			"source-layer": "signal-point",
			type: "circle",
		});
	});
};

const selectedSignalLocation = () => {
	const first_pool = props.selectedSignals
		.values()
		.reduce((accumulator, current) => {
			if (accumulator == null && current.type === "flyover pool") {
				return current;
			}
			return accumulator;
		}, null);
	return first_pool?.location;
};
const showMapTile = () => {
	return props.selectedSignals.value
		.values()
		.reduce(
			(accumulator, current) => accumulator || current.type === "flyover pool",
			false,
		);
};
const updateSignalLocation = (event) => {
	const signalId = event.detail.signalId;
	console.log("map click", signalId, event.detail);

	const map = event.detail.map;
	const loc = {
		latitude: event.detail.lat,
		longitude: event.detail.lng,
	};

	map.SetMarkers([loc]);
	poolLocations.value[signalId] = loc;
};
</script>
