<template>
	<p>A flyover pool</p>
	<div v-if="session.organization">
		<MapProxiedArcgisTile
			:location="location"
			:markers="markers"
			:organizationId="session.organization.id"
			:tegola="session.urls?.tegola ?? ''"
			:urlTiles="session.urls?.tile ?? ''"
		/>
	</div>
	<div v-else>
		<p>Loading...</p>
	</div>
</template>

<script setup lang="ts">
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";
import { Location } from "@/type/api";
import { Marker } from "@/types";
import { useSessionStore } from "@/store/session";

interface Props {
	location: Location;
	markers: Marker[];
}
const props = defineProps<Props>();
const session = useSessionStore();
</script>
