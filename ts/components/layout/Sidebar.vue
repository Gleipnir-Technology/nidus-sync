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
						v-show="notificationCounts.communication > 0"
						class="position-absolute translate-middle badge rounded-pill bg-primary"
					>
						<span>{{
							notificationCounts.communication > 99
								? "99+"
								: notificationCounts.communication
						}}</span>
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
						v-show="notificationCounts.review > 0"
						class="position-absolute translate-middle badge rounded-pill bg-primary"
					>
						<span>{{
							notificationCounts.review > 99 ? "99+" : notificationCounts.review
						}}</span>
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

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from "vue";
import { Tooltip, Popover } from "bootstrap";

// Reactive state
const isCollapsed = ref(false);
const notificationCounts = reactive({
	communication: 0,
	review: 0,
});

const userData = reactive({});

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

// Fetch user state from API
const updateUserState = async () => {
	try {
		const response = await fetch("/api/user/self");
		const data = await response.json();

		// Update reactive data
		Object.keys(data).forEach((key) => {
			if (key === "notification_counts") {
				Object.assign(notificationCounts, data[key]);
			} else {
				userData[key] = data[key];
			}
		});
	} catch (error) {
		console.error("Failed to update user state:", error);
	}
};

// Lifecycle hooks
onMounted(async () => {
	restoreLocalStorage();

	await nextTick();

	initializeBootstrap();
	setTooltipsForSidebar();

	// Subscribe to SSE events (assuming SSEManager is globally available)
	if (window.SSEManager) {
		sseUnsubscribe = window.SSEManager.subscribe("*", (e) => {
			if (e.type !== "heartbeat") {
				updateUserState();
			}
		});
	}

	// Initial user state fetch
	updateUserState();
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
