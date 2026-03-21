<template>
	<div class="app-container">
		<Sidebar />
		<MainContent>
			<router-view />
		</MainContent>
	</div>
</template>

<script setup lang="ts">
import Sidebar from "./components/layout/Sidebar.vue";
import MainContent from "./components/layout/MainContent.vue";
</script>

<style scoped>
.app-container {
	display: flex;
	height: 100vh;
}
</style>
