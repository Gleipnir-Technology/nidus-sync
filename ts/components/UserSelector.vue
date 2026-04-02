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
				<Avatar :user="user" />

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
import { ref, computed, onMounted, watch } from "vue";
import Avatar from "@/components/Avatar.vue";
import { useUserStore } from "@/store/user";
import type { User } from "@/types";

interface Props {
	modelValue?: User | null;
	placeholder?: string;
	minChars?: number;
}

const props = withDefaults(defineProps<Props>(), {
	modelValue: null,
	placeholder: "Search users...",
	minChars: 3,
});

const emit = defineEmits<{
	"update:modelValue": [user: User | null];
}>();

const usersStore = useUserStore();
const searchQuery = ref("");
const showDropdown = ref(false);

onMounted(async () => {
	// Fetch all users if not already loaded
	if (!usersStore.all) {
		await usersStore.fetchAll();
	}

	// Initialize search query with selected user's name if provided
	if (props.modelValue) {
		searchQuery.value = props.modelValue.display_name;
	}
});

// Watch for external changes to modelValue
watch(
	() => props.modelValue,
	(newValue) => {
		if (newValue) {
			searchQuery.value = newValue.display_name;
		} else {
			searchQuery.value = "";
		}
	},
);

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

	// Clear selection if user is typing
	if (props.modelValue && searchQuery.value !== props.modelValue.display_name) {
		emit("update:modelValue", null);
	}
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
	emit("update:modelValue", user);
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
