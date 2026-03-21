<template>
	<a
		:href="to"
		data-bs-toggle="tooltip"
		data-bs-placement="right"
		:title="label"
	>
		<div class="menu-icon"><i :class="`bi bi-${icon}`"></i></div>
		<span class="menu-text ms-2">{{ label }}</span>
		<span
			v-if="notificationCount != null && notificationCount > 0"
			class="position-absolute translate-middle badge rounded-pill bg-primary"
		>
			<span>{{ notificationCount > 99 ? "99+" : notificationCount }}</span>
			<span class="visually-hidden">unread notifications</span>
		</span>
	</a>
</template>

<script setup lang="ts">
defineProps<{
	to: string;
	icon: string;
	label: string;
	notificationCount?: number | null;
}>();
</script>
