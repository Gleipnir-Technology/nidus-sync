<template>
	<aside class="sidebar">
		<h2>My App</h2>
		<nav>
			<NavigationLink to="/" label="Home" />
			<NavigationLink to="/about" label="About" />
		</nav>
	</aside>
</template>

<script setup lang="ts">
import NavigationLink from "../common/NavigationLink.vue";
</script>

<style scoped>
.sidebar {
	width: 250px;
	background-color: #2c3e50;
	color: white;
	padding: 20px;
}

nav {
	margin-top: 30px;
	display: flex;
	flex-direction: column;
	gap: 10px;
}
</style>
