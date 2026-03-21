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
			<li>
				<a
					href="/communication"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Communication"
				>
					<div class="menu-icon"><i class="bi bi-messaging"></i></div>
					<span class="menu-text ms-2">Communication</span>
					<span
						x-show="notification_counts.communication > 0"
						x-cloak
						class="position-absolute translate-middle badge rounded-pill bg-primary"
					>
						<span
							x-text="notification_counts.communication > 99 ? '99+' : notification_counts.communication"
						></span>
						<span class="visually-hidden">unread notifications</span>
					</span>
				</a>
			</li>
			<li>
				<a
					href="/planning"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Planning"
				>
					<div class="menu-icon"><i class="bi bi-strategy"></i></div>
					<span class="menu-text ms-2">Planning</span>
				</a>
			</li>
			<li>
				<a
					href="Operations"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Operations"
				>
					<div class="menu-icon"><i class="bi bi-assign"></i></div>
					<span class="menu-text ms-2">Operations</span>
				</a>
			</li>
			<li>
				<a
					href="/review"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Review"
				>
					<div class="menu-icon"><i class="bi bi-review"></i></div>
					<span class="menu-text ms-2">Review</span>
					<span
						x-show="notification_counts.review > 0"
						x-cloak
						class="position-absolute translate-middle badge rounded-pill bg-primary"
					>
						<span
							x-text="notification_counts.review > 99 ? '99+' : notification_counts.review"
						></span>
						<span class="visually-hidden">unread notifications</span>
					</span>
				</a>
			</li>
			<li>
				<a
					href="/configuration"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Configuration"
				>
					<div class="menu-icon"><i class="bi bi-settings"></i></div>
					<span class="menu-text ms-2">Configuration</span>
				</a>
			</li>
			<li>
				<a
					href="/sudo"
					data-bs-toggle="tooltip"
					data-bs-placement="right"
					title="Sudo"
				>
					<div class="menu-icon"><i class="bi bi-god"></i></div>
					<span class="menu-text ms-2">Sudo</span>
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
