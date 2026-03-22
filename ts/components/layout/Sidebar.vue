<template>
	<div id="sidebar" :class="{ collapsed: isCollapsed }">
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
				<NavigationLink to="/intelligence" icon="brain" label="Intelligence" />
			</li>
			<li>
				<NavigationLink
					to="/communication"
					icon="messaging"
					label="Communication"
					:notificationCount="user?.notification_counts?.communication ?? 0"
				/>
			</li>
			<li>
				<NavigationLink to="/planning" icon="strategy" label="Planning" />
			</li>
			<li>
				<NavigationLink to="/operations" icon="assign" label="Operations" />
			</li>
			<li>
				<NavigationLink
					to="/review"
					icon="review"
					label="Review"
					:notificationCount="user?.notification_counts?.review ?? 0"
				/>
			</li>
			<li>
				<NavigationLink
					to="/configuration"
					icon="assign"
					label="Configuration"
				/>
			</li>
			<li>
				<NavigationLink to="/sudo" icon="god" label="Sudo" />
			</li>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from "vue";
import { Tooltip, Popover } from "bootstrap";
import NavigationLink from "../common/NavigationLink.vue";
import { useUserStore } from "../../store/user";

// Reactive state
const isCollapsed = ref(false);

const user = useUserStore();

// Bootstrap tooltip instances
let tooltipInstances = [];
let sseUnsubscribe = null;

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

	initializeBootstrap();
	setTooltipsForSidebar();
});

onBeforeUnmount(() => {
	// Cleanup Bootstrap tooltips
	tooltipInstances.forEach((tooltip) => tooltip.dispose());

	// Unsubscribe from SSE
	if (sseUnsubscribe) {
		sseUnsubscribe();
	}
});
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
