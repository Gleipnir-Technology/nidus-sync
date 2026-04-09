<style scoped>
.avatar {
	width: 48px;
	height: 48px;
	object-fit: cover;
}
.bi-avatar {
	height: 48px;
	width: 48px;
}
</style>
<template>
	<img
		v-if="user.avatar"
		:src="user.avatar"
		:alt="user.display_name"
		class="rounded-circle me-2 avatar"
		width="32"
		height="32"
	/>
	<i v-else class="bi bi-avatar"></i>
</template>

<script setup lang="ts">
import { User } from "@/type/api";

interface Props {
	user: User;
}
const props = defineProps<Props>();
</script>
