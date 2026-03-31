<template>
	<div class="user-selector">
		<label v-if="label" :for="selectId" class="form-label">
			{{ label }}
		</label>

		<select
			:id="selectId"
			v-model="selectedUserId"
			class="form-select"
			:class="{ 'is-invalid': error }"
			:disabled="loading || disabled"
			@change="handleChange"
		>
			<option :value="null">{{ placeholder }}</option>
			<option v-for="user in usersStore.all" :key="user.id" :value="user.id">
				{{ user.display_name || user.username || `User ${user.id}` }}
			</option>
		</select>

		<div v-if="loading" class="form-text">
			<span
				class="spinner-border spinner-border-sm me-2"
				role="status"
				aria-hidden="true"
			></span>
			Loading users...
		</div>

		<div v-if="error" class="invalid-feedback d-block">
			Failed to load users. Please try again.
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import { useUserStore } from "@/store/user";
import type { User } from "@/types";

interface Props {
	modelValue?: number | null;
	label?: string;
	placeholder?: string;
	disabled?: boolean;
}

interface Emits {
	(e: "update:modelValue", value: number | null): void;
	(e: "change", user: User | null): void;
}

const props = withDefaults(defineProps<Props>(), {
	modelValue: null,
	label: "",
	placeholder: "Select a user...",
	disabled: false,
});

const emit = defineEmits<Emits>();

const usersStore = useUserStore();
const selectedUserId = ref<number | null>(props.modelValue);
const loading = ref(false);
const error = ref(false);
const selectId = computed(
	() => `user-select-${Math.random().toString(36).substr(2, 9)}`,
);

// Watch for external changes to modelValue
watch(
	() => props.modelValue,
	(newValue) => {
		selectedUserId.value = newValue;
	},
);

const handleChange = () => {
	emit("update:modelValue", selectedUserId.value);

	const selectedUser = selectedUserId.value
		? usersStore.byID(selectedUserId.value)
		: null;

	emit("change", selectedUser || null);
};

onMounted(async () => {
	// Only fetch if users haven't been loaded yet
	if (!usersStore.all) {
		loading.value = true;
		error.value = false;

		try {
			await usersStore.fetchAll();
		} catch (err) {
			console.error("Failed to fetch users:", err);
			error.value = true;
		} finally {
			loading.value = false;
		}
	}
});
</script>

<style scoped>
.user-selector {
	margin-bottom: 1rem;
}
</style>
