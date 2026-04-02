<style scoped lang="scss">
#content {
	transition: all 0.3s;
	margin-left: 250px;
	padding: 10px;
	width: calc(100% - 250px);
}

#content.expanded {
	margin-left: 70px;
	width: calc(100% - 70px);
}

.logo-container {
	display: flex;
	justify-content: center;
	align-items: center;
	transition: all 0.3s ease;
}

.logo {
	max-width: 100%;
	height: auto;
	transition: all 0.3s ease;
}

#sidebar {
	background-color: $off-white;
	min-height: 100vh;
	transition: all 0.3s;
	width: 250px;
	position: fixed;
	z-index: 1000;
	padding: 20px;
}

#sidebar.collapsed {
	width: 70px;
	padding: 20px 10px;
}
/* Logo style when sidebar is collapsed */
#sidebar.collapsed .logo-container {
	width: 100%;
}

#sidebar.collapsed .logo-img {
	max-width: 40px; /* smaller size for collapsed state */
}
#sidebar.impersonating {
	background-color: $danger;
}
#sidebar.collapsed .menu-text {
	opacity: 0;
	visibility: hidden;
	width: 0;
}

#sidebar.collapsed .sidebar-header h4 {
	opacity: 0;
	visibility: hidden;
}

#sidebar.collapsed .sidebar-menu .menu-icon {
	min-width: 100%;
	font-size: 1.5rem;
}

#sidebarToggle {
	position: absolute;
	left: calc(250px - 15px);
	top: 50%;
	transform: translateY(-50%);
	z-index: 1050;
	width: 30px;
	height: 30px;
	border-radius: 50%;
	border: 1px solid #dee2e6;
	display: flex;
	align-items: center;
	transition: left 0.3s;
	padding: 0;
}
#sidebarToggle i {
	transition: transform 0.3s;
}

#sidebar.collapsed > #sidebarToggle {
	left: calc(70px - 15px);
}

#sidebar > #sidebarToggle i {
	position: relative;
	left: 5px;
}

#sidebar.collapsed > #sidebarToggle i {
	transform: rotate(180deg);
}
.sidebar-header {
	padding-bottom: 20px;
	border-bottom: 1px solid $off-black;
	margin-bottom: 20px;
	overflow: hidden;
	white-space: nowrap;
	display: flex;
	justify-content: center; /* Center for the logo */
}

.sidebar-menu {
	list-style: none;
	padding: 0;
}

.sidebar-menu li {
	padding: 10px 0;
}

.sidebar-menu li a {
	text-decoration: none;
	color: $off-black;
	display: flex;
	align-items: center;
	overflow: hidden;
	white-space: nowrap;
}

.sidebar-menu li a:hover {
	color: $primary;
}

.sidebar-menu .menu-icon {
	font-size: 1.2rem;
	min-width: 30px;
	display: flex;
	justify-content: center;
}
.sidebar-menu .menu-icon svg {
	width: 1.5em;
	height: 1.5em;
}
.sidebar-menu .menu-text {
	transition: opacity 0.3s;
}
</style>
<template>
	<div
		id="sidebar"
		:class="{ collapsed: isCollapsed, impersonating: isImpersonating }"
	>
		<div class="sidebar-header">
			<div class="logo-container">
				<img class="logo" src="/static/img/nidus-logo-256-transparent.png" />
			</div>
		</div>

		<button id="sidebarToggle" class="btn btn-sm p-0" @click="toggleSidebar">
			<i id="sidebarToggleIcon" class="bi bi-chevron-left"></i>
		</button>

		<ul class="sidebar-menu">
			<li>
				<NavigationLink to="/" icon="house" label="Home" />
			</li>
			<li>
				<NavigationLink
					to="/_/intelligence"
					icon="brain"
					label="Intelligence"
				/>
			</li>
			<li>
				<NavigationLink
					to="/_/communication"
					icon="messaging"
					label="Communication"
					:notificationCount="session.notification_counts?.communication ?? 0"
				/>
			</li>
			<li>
				<NavigationLink to="/_/planning" icon="strategy" label="Planning" />
			</li>
			<li>
				<NavigationLink to="/_/operations" icon="assign" label="Operations" />
			</li>
			<li>
				<NavigationLink
					to="/_/review"
					icon="review"
					label="Review"
					:notificationCount="session.notification_counts?.review ?? 0"
				/>
			</li>
			<li>
				<NavigationLink
					to="/_/configuration"
					icon="assign"
					label="Configuration"
				/>
			</li>
			<li>
				<NavigationLink to="/_/sudo" icon="god" label="Sudo" />
			</li>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from "vue";
import { Tooltip, Popover } from "bootstrap";
import NavigationLink from "@/components/common/NavigationLink.vue";
import { SSEManager } from "@/SSEManager";
import { useSessionStore } from "@/store/session";

// Reactive state
const isCollapsed = ref(false);
const isImpersonating = ref(false);

const session = useSessionStore();

// Bootstrap tooltip instances
let tooltipInstances: Tooltip[] = [];

// Initialize Bootstrap components
const initializeBootstrap = () => {
	// Initialize popovers
	const popoverElements = document.querySelectorAll(
		'[data-bs-toggle="popover"]',
	);
	popoverElements.forEach((el) => new Popover(el));
	console.log("Initialized", popoverElements.length, "popovers");

	// Initialize tooltips
	const tooltipElements = document.querySelectorAll(
		'[data-bs-toggle="tooltip"]',
	);
	tooltipInstances = Array.from(tooltipElements).map((el) => new Tooltip(el));
	console.log("Initialized", tooltipElements.length, "tooltips");
};

// Restore sidebar state from localStorage
const restoreLocalStorage = () => {
	const expanded = localStorage.getItem("sidebar.expanded");
	if (expanded === "false") {
		isCollapsed.value = true;
		document.getElementById("content")?.classList.add("expanded");
	} else {
		isCollapsed.value = false;
		document.getElementById("content")?.classList.remove("expanded");
		localStorage.setItem("sidebar.expanded", "true");
	}
};

// Toggle sidebar collapsed state
const toggleSidebar = () => {
	isCollapsed.value = !isCollapsed.value;
	document.getElementById("content")?.classList.toggle("expanded");
	setTooltipsForSidebar();
	localStorage.setItem("sidebar.expanded", (!isCollapsed.value).toString());
};

// Enable/disable tooltips based on sidebar state
const setTooltipsForSidebar = () => {
	const sidebarTooltips = document.querySelectorAll(
		'#sidebar [data-bs-toggle="tooltip"]',
	);
	const isExpanded = document
		.getElementById("content")
		?.classList.contains("expanded");

	sidebarTooltips.forEach((el) => {
		const tooltip = Tooltip.getOrCreateInstance(el);
		if (isExpanded) {
			tooltip.enable();
		} else {
			tooltip.disable();
		}
	});
};

// Lifecycle hooks
onMounted(async () => {
	restoreLocalStorage();

	await nextTick();

	const s = await session.get();
	isImpersonating.value = !!s.impersonating;
	initializeBootstrap();
	setTooltipsForSidebar();
});

onBeforeUnmount(() => {
	// Cleanup Bootstrap tooltips
	tooltipInstances.forEach((tooltip) => tooltip.dispose());
});
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
