<style scoped>
.app-container {
	display: flex;
	height: 100vh;
}
</style>
<template>
	<div class="app-container">
		<template v-if="session.loading">Loading...</template>
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
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
import Sidebar from "@/components/layout/Sidebar.vue";
import MainContent from "@/components/layout/MainContent.vue";
</script>
