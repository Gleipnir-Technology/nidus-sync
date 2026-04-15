<template>
	<p>A flyover pool</p>
	<div v-if="session.organization && session.urls">
		<MapProxiedArcgisTile
			:location="location"
			:markers="markers"
			:organizationId="session.organization.id"
			:tegola="session.urls!.tegola"
			:urlTiles="session.urls!.tile"
			v-model="cameraFlyover"
		/>
	</div>
	<div v-else>
		<p>Loading...</p>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { useSessionStore } from "@/store/session";
import { Marker } from "@/types";
import { Location } from "@/type/api";
import { Camera } from "@/type/map";

interface Props {
	location: Location;
	markers: Marker[];
}
const cameraFlyover = ref<Camera>(new Camera());
const props = defineProps<Props>();
const session = useSessionStore();
</script>
