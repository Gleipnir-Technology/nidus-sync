export function SetupSidebar() {
	var popoverTriggerList = [].slice.call(
		document.querySelectorAll('[data-bs-toggle="popover"]'),
	);
	var popoverList = popoverTriggerList.map(function (popoverTriggerEl) {
		return new bootstrap.Popover(popoverTriggerEl);
	});
	console.log("Initialized ", popoverTriggerList.length, " popovers");

	var tooltipTriggerList = [].slice.call(
		document.querySelectorAll('[data-bs-toggle="tooltip"]'),
	);
	var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
		let t = new bootstrap.Tooltip(tooltipTriggerEl);
		return t;
	});
	console.log("Initialized ", tooltipTriggerList.length, " tooltips");
	restoreLocalStorage();
	setTooltipsForSidebar();
	SSEManager.subscribe("*", function (e) {
		if (e.type != "heartbeat") {
			updateUserState();
		}
	});
	document.getElementById("sidebarToggle").addEventListener("click", () => {
		const sidebar = document.getElementById("sidebar");
		sidebar.classList.toggle("collapsed");
		document.getElementById("content").classList.toggle("expanded");
		setTooltipsForSidebar();
		localStorage.setItem(
			"sidebar.expanded",
			(!sidebar.classList.contains("collapsed")).toString(),
		);
	});
	updateUserState();
}
function restoreLocalStorage() {
	const expanded = localStorage.getItem("sidebar.expanded");
	if (expanded == "false") {
		document.getElementById("sidebar").classList.add("collapsed");
		document.getElementById("content").classList.add("expanded");
	} else {
		document.getElementById("sidebar").classList.remove("collapsed");
		document.getElementById("content").classList.remove("expanded");
		localStorage.setItem("sidebar.expanded", "true");
	}
}
function setTooltipsForSidebar() {
	const sidebarTooltips = document.querySelectorAll(
		'#sidebar [data-bs-toggle="tooltip"]',
	);
	const isExpanded = document
		.getElementById("content")
		.classList.contains("expanded");
	sidebarTooltips.forEach((t) => {
		const tooltip = bootstrap.Tooltip.getOrCreateInstance(t);
		if (isExpanded) {
			tooltip.enable();
		} else {
			tooltip.disable();
		}
	});
}
async function updateUserState() {
	const response = await fetch("/api/user/self");
	const data = await response.json();
	Object.keys(data).forEach((key) => {
		store_user[key] = data[key];
	});
}
