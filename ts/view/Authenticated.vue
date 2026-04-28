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
			<UpdateNotification :updateDate="updateDate" />
		</template>
	</div>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";

import Sidebar from "@/components/layout/Sidebar.vue";
import MainContent from "@/components/layout/MainContent.vue";
import UpdateNotification from "@/components/UpdateNotification.vue";
import { SSEManager, type SSEMessageStatus } from "@/SSEManager";
import { Session } from "@/type/api";
import { router } from "@/route/config";
import { useRoutes } from "@/route/use";
import { useSessionStore } from "@/store/session";

const revision = ref<string>("");
const session = useSessionStore();
const updateDate = ref<Date | null>(null);
onMounted(() => {
	SSEManager.subscribeStatus((msg: SSEMessageStatus) => {
		if (msg.status == "connected") {
			if (revision.value == "") {
				revision.value = msg.revision;
			} else {
				updateDate.value = new Date();
			}
		}
		console.log("status update:", msg);
	});
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
