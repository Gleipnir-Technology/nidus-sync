<template>
	<div class="app-container">
		<Sidebar v-if="$route.meta.showSidebar" />
		<MainContent>
			<div v-if="session.loading">Loading...</div>
			<div v-else-if="session.error">Error: {{ session.error }}</div>
			<router-view v-else />
		</MainContent>
	</div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useSessionStore } from "@/store/session";
import { Session } from "@/types";

import Sidebar from "./components/layout/Sidebar.vue";
import MainContent from "./components/layout/MainContent.vue";
import NavigationLink from "@/components/common/NavigationLink.vue";

const session = useSessionStore();
onMounted(() => {
	session.get().then((session: Session) => {
		console.log("session loaded", session);
	});
});
</script>

<style scoped>
.app-container {
	display: flex;
	height: 100vh;
}
</style>
