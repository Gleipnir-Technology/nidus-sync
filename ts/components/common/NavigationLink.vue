<template>
	<router-link :to="to" class="nav-link">
		{{ label }}
	</router-link>
</template>

<script setup lang="ts">
defineProps<{
	to: string;
	label: string;
}>();
</script>

<style scoped>
.nav-link {
	color: #ecf0f1;
	text-decoration: none;
	padding: 10px;
	border-radius: 4px;
	transition: background-color 0.2s;
}

.nav-link:hover {
	background-color: #34495e;
}

.nav-link.router-link-active {
	background-color: #3498db;
}
</style>
