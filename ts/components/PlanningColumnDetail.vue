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
			<div class="map-container" v-if="session.organization">
				<MapMultipoint
					id="map"
					:bounds="session.organization.service_area"
					:markers="markers"
					:organizationId="session.organization.id"
					:tegola="session.urls?.tegola ?? ''"
				></MapMultipoint>
			</div>
			<div v-else>
				<p>loading...</p>
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

							<div class="mt-2" v-show="selectedSignals.length > 0">
								<div v-for="signal in selectedSignals" :key="signal.id">
									<PlanningColumnDetailEntry :signal="signal" />
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>

			<div v-show="showMapTile" class="map-container">
				<MapProxiedArcgisTile
					:location="selectedSignalLocation()"
					:markers="[]"
					:organizationId="session.organization?.id ?? 0"
					:tegola="session.urls?.tegola ?? ''"
					:urlTiles="session.urls?.tile ?? ''"
					@map-click="updateSignalLocation"
				>
				</MapProxiedArcgisTile>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import MapMultipoint from "@/components/MapMultipoint.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import PlanningColumnDetailEntry from "@/components/PlanningColumnDetailEntry.vue";
import TimeRelative from "@/components/TimeRelative.vue";
import { shortAddress } from "@/format";
import { useSessionStore } from "@/store/session";
import { Location, MapClickEvent, Marker, Signal } from "@/types";

interface Props {
	markers: Marker[];
	selectedSignals: Array<Signal>;
}
const props = defineProps<Props>();
const session = useSessionStore();
const configureMapTile = () => {
	/*
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
*/
};

const selectedSignalLocation = (): Location => {
	const first_pool = props.selectedSignals.reduce(
		(accumulator: Signal | null, current: Signal) => {
			if (accumulator == null && current.type === "flyover pool") {
				return current;
			}
			return accumulator;
		},
		null as Signal | null,
	);
	const loc = first_pool?.location;
	return (
		loc || {
			lat: 0,
			lng: 0,
		}
	);
};
const showMapTile = () => {
	const hasPool = props.selectedSignals.reduce(
		(accumulator: boolean | null, current: Signal) => {
			return accumulator || current.type === "flyover pool";
		},
		false,
	);
	return selectedSignalLocation() && hasPool;
};
const updateSignalLocation = (event: MapClickEvent) => {
	console.log("map click", event.location);
	//const signalId = event.detail.signalId;

	//const map = event.map;
	//const loc = event.location;

	//map.SetMarkers([loc]);
	//poolLocations.value[signalId] = loc;
};
</script>
