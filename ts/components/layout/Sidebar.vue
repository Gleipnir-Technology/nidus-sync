<template>
	<div id="sidebar" x-data="$store.user">
		<div class="sidebar-header">
			<div class="logo-container">
				<img class="logo" src="/static/img/nidus-logo-256-transparent.png" />
			</div>
		</div>
		<button id="sidebarToggle" class="btn btn-sm p-0">
			<i id="sidebarToggleIcon" class="bi bi-chevron-left"></i>
		</button>

		<ul class="sidebar-menu">
			<li>
				<a
					href="/root"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Home"
				>
					<div class="menu-icon"><i class="bi bi-house"></i></div>
					<span class="menu-text ms-2">Home</span>
				</a>
			</li>
			<li>
				<a
					href="/intelligence"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Intelligence"
				>
					<div class="menu-icon"><i class="bi bi-brain"></i></div>
					<span class="menu-text ms-2">Intelligence</span>
				</a>
			</li>
		</ul>
	</div>
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
