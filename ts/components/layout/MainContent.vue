<template>
	<div id="content">
		<div v-if="session.loading">Loading...</div>
		<div v-else-if="session.error">Error: {{ session.error }}</div>
		<slot />
	</div>
</template>

<script setup lang="ts">
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
</script>

<style scoped></style>
