<template>
	<div>
		<h1>About</h1>
		<p>This is an example Vue 3 + TypeScript application.</p>

		<div class="info-card">
			<h2>Features</h2>
			<ul>
				<li>Vue 3 Composition API</li>
				<li>TypeScript support</li>
				<li>Vue Router for navigation</li>
				<li>Composables for shared state</li>
				<li>Single File Components</li>
			</ul>
		</div>

		<div class="info-card">
			<h2>Tech Stack</h2>
			<div class="tech-list">
				<span class="tech-badge">Vue 3</span>
				<span class="tech-badge">TypeScript</span>
				<span class="tech-badge">Vue Router</span>
				<span class="tech-badge">Vite</span>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
// No special logic needed for this component
</script>

<style scoped>
.info-card {
	margin-top: 24px;
	padding: 24px;
	background: white;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.info-card h2 {
	margin-top: 0;
	color: #2c3e50;
}

.info-card ul {
	line-height: 1.8;
}

.tech-list {
	display: flex;
	gap: 12px;
	flex-wrap: wrap;
}

.tech-badge {
	padding: 8px 16px;
	background-color: #3498db;
	color: white;
	border-radius: 20px;
	font-size: 14px;
	font-weight: 500;
}
</style>
