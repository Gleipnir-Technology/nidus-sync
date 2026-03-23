<template>
	<div class="app-container">
		<Sidebar v-if="$route.meta.showSidebar" />
		<MainContent>
			<div v-if="userStore.loading">Loading...</div>
			<div v-else-if="userStore.error">Error: {{ userStore.error }}</div>
			<router-view v-else />
		</MainContent>
	</div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useUserStore } from "@/store/user";

import Sidebar from "./components/layout/Sidebar.vue";
import MainContent from "./components/layout/MainContent.vue";
import NavigationLink from "@/components/common/NavigationLink.vue";

const userStore = useUserStore();
onMounted(() => {
	userStore.fetchUser();
});
</script>

<style scoped>
.app-container {
	display: flex;
	height: 100vh;
}
</style>
