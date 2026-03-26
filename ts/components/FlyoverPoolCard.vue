<template>
	<p>A flyover pool</p>
	<MapProxiedArcgisTile
						id="tile-map"
						:latitude="pool.location.latitude"
						:longitude="pool.location.longitude"
						:markers="tileMapMarkers"
						:organizationId="user.organization.id"
						:tegola="user.urls.tegola" />
</template>

<script setup lang="ts">
import { useUserStore } from "../store/user";
import MapProxiedArcgisTile from "@/components/MapProxiedArcgisTile.vue";

interface Props {
	pool: Pool;
}
const props = defineProps<Props>();
const user = useUserStore();
</script>
