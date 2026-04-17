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
	<div class="card" v-show="!!selectedSite">
		<table class="table">
			<tbody>
				<tr>
					<td><b>Address</b></td>
					<td>
						{{ formatAddress(selectedSite?.address) }}
						{{ selectedSite?.address.region }}
						{{ selectedSite?.address.postal_code }}
					</td>
				</tr>
				<tr>
					<td><b>Owner</b></td>
					<td>{{ selectedSite?.owner?.name }}</td>
				</tr>
				<tr>
					<td><b>Parcel APN</b></td>
					<td>{{ selectedSite?.parcel?.apn }}</td>
				</tr>
				<tr>
					<td><b>Parcel Description</b></td>
					<td>{{ selectedSite?.parcel?.description }}</td>
				</tr>
				<tr>
					<td>
						Features<br />
						<ul>
							<li v-for="(feature, index) in selectedSite?.features">
								{{ feature.type }}
							</li>
						</ul>
					</td>
				</tr>
				<tr>
					<td>
						<table>
							<tbody>
								<tr v-for="(lead, index) in selectedSite?.leads">
									<td>{{ lead.type }}</td>
									<td>
										<ul>
											<li
												v-for="(crr, index2) in lead.compliance_report_requests"
											>
												<a
													:href="mailerLink(crr)"
													target="_blank"
													rel="noopener noreferrer"
													>Compliance Report Request {{ crr.public_id }}</a
												>
											</li>
										</ul>
									</td>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
	<div class="map-container" v-if="session.organization">
		<MapLocator
			:markers="mapMarkers"
			:useSatellite="true"
			v-model="mapCamera"
		></MapLocator>
	</div>
	<div class="map-container" v-if="session.organization && session.urls">
		<MapProxiedArcgisTile
			:markers="mapMarkers"
			:organizationId="session.organization.id"
			:tegola="session.urls!.tegola"
			:urlTiles="session.urls!.tile"
			v-model="_mapFlyoverCamera"
		></MapProxiedArcgisTile>
	</div>
</template>
<script setup lang="ts">
import { ref, watch } from "vue";
import MapLocator from "@/components/MapLocator.vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { formatAddress } from "@/format";
import { useSessionStore } from "@/store/session";
import { ComplianceReportRequest, Site } from "@/type/api";
import { Camera } from "@/type/map";
import type { Marker } from "@/types";

interface Emits {
	(e: "update:modelValue", value: any): void;
}
interface Props {
	mapFlyoverCamera: Camera;
	mapMarkers: Marker[];
	selectedSite: Site | undefined;
}
const _mapFlyoverCamera = ref<Camera>(new Camera());
const mapCamera = ref<Camera>(new Camera());
const props = defineProps<Props>();
const session = useSessionStore();

function mailerLink(crr: ComplianceReportRequest): string {
	return `/mailer/mode-3/${crr.public_id}/preview`;
}
watch(
	() => props.mapFlyoverCamera,
	(newMapFlyoverCamera: Camera) => {
		console.log("map flyover camera update", newMapFlyoverCamera);
		_mapFlyoverCamera.value = newMapFlyoverCamera;
	},
);
</script>
