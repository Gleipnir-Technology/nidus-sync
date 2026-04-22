<style scoped>
.app-container {
	display: flex;
	height: 100vh;
}
</style>
<template>
	<div class="app-container">
		<template v-if="session.isLoading">Loading...</template>
		<template v-else-if="session.error">Error: {{ session.error }}</template>
		<template v-else>
			<Sidebar />
			<MainContent>
				<router-view v-slot="{ Component }">
					<component :is="Component" />
				</router-view>
			</MainContent>
		</template>
	</div>
</template>
<script setup lang="ts">
import { onMounted } from "vue";

import Sidebar from "@/components/layout/Sidebar.vue";
import MainContent from "@/components/layout/MainContent.vue";
import { Session } from "@/type/api";
import { router } from "@/route/config";
import { useRoutes } from "@/route/use";
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
onMounted(() => {
	session
		.get()
		.then((session: Session) => {
			console.log("session loaded by Authenticated", session);
		})
		.catch((e) => {
			console.log(
				"root session not loaded, user is not authenticated",
				router.currentRoute.value.fullPath,
			);
			router.push(`/signin?next=${router.currentRoute.value.fullPath}`);
		});
});
</script>
