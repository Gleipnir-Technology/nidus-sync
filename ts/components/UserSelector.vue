<template>
	<div class="user-autocomplete position-relative">
		<input
			type="text"
			class="form-control"
			v-model="searchQuery"
			@input="onInput"
			@focus="onFocus"
			@blur="onBlur"
			:placeholder="placeholder"
			autocomplete="off"
		/>

		<div
			v-if="showDropdown && filteredUsers.length > 0"
			class="dropdown-menu show w-100"
			style="max-height: 300px; overflow-y: auto"
		>
			<a
				v-for="user in filteredUsers"
				:key="user.id"
				href="#"
				class="dropdown-item d-flex align-items-center"
				@mousedown.prevent="selectUser(user)"
			>
				<img
					v-if="user.avatar"
					:src="user.avatar"
					:alt="user.display_name"
					class="rounded-circle me-2"
					width="32"
					height="32"
				/>
				<span v-else class="badge bg-secondary me-2">{{ user.initials }}</span>

				<div class="flex-grow-1">
					<div v-html="highlightMatch(user.display_name)"></div>
					<small
						class="text-muted"
						v-html="highlightMatch(user.username)"
					></small>
				</div>
			</a>
		</div>

		<div
			v-if="
				showDropdown &&
				searchQuery.length >= minChars &&
				filteredUsers.length === 0
			"
			class="dropdown-menu show w-100"
		>
			<div class="dropdown-item text-muted">No users found</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useUserStore } from "@/store/user";
import type { User } from "@/types";

interface Props {
	placeholder?: string;
	minChars?: number;
}

const props = withDefaults(defineProps<Props>(), {
	placeholder: "Search users...",
	minChars: 3,
});

const emit = defineEmits<{
	select: [user: User];
	input: [query: string];
}>();

const usersStore = useUserStore();
const searchQuery = ref("");
const showDropdown = ref(false);

onMounted(async () => {
	// Fetch all users if not already loaded
	if (!usersStore.all) {
		await usersStore.fetchAll();
	}
});

const filteredUsers = computed(() => {
	if (searchQuery.value.length < props.minChars || !usersStore.all) {
		return [];
	}

	const query = searchQuery.value.toLowerCase();

	return usersStore.all
		.filter((user: User) => {
			const displayName = user.display_name.toLowerCase();
			const username = user.username.toLowerCase();
			return displayName.includes(query) || username.includes(query);
		})
		.slice(0, 10); // Limit to 10 results
});

function onInput() {
	showDropdown.value = searchQuery.value.length >= props.minChars;
	emit("input", searchQuery.value);
}

function onFocus() {
	if (searchQuery.value.length >= props.minChars) {
		showDropdown.value = true;
	}
}

function onBlur() {
	// Delay to allow click event on dropdown items
	setTimeout(() => {
		showDropdown.value = false;
	}, 200);
}

function selectUser(user: User) {
	searchQuery.value = user.display_name;
	showDropdown.value = false;
	emit("select", user);
}

function highlightMatch(text: string): string {
	if (!searchQuery.value || searchQuery.value.length < props.minChars) {
		return escapeHtml(text);
	}

	const query = escapeHtml(searchQuery.value);
	const escapedText = escapeHtml(text);
	const regex = new RegExp(`(${query})`, "gi");

	return escapedText.replace(regex, "<strong>$1</strong>");
}

function escapeHtml(text: string): string {
	const div = document.createElement("div");
	div.textContent = text;
	return div.innerHTML;
}
</script>

<style scoped>
.user-autocomplete {
	position: relative;
}

.dropdown-menu {
	position: absolute;
	top: 100%;
	left: 0;
	z-index: 1000;
}

.dropdown-item {
	cursor: pointer;
}

.dropdown-item:hover {
	background-color: #f8f9fa;
}
</style>
