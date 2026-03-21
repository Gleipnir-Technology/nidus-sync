<template>
	<main class="main-content">
		<slot />
	</main>
</template>

<script setup lang="ts">
// No imports needed for this simple component
</script>

<style scoped>
.main-content {
	flex: 1;
	padding: 40px;
	overflow-y: auto;
	background-color: #f5f5f5;
}
</style>
