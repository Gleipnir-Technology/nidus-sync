<style scoped>
.btn {
	display: inline-flex;
	align-items: center;
	justify-content: center;
}

.btn:disabled {
	cursor: not-allowed;
}
</style>
<template>
	<button
		:class="buttonClasses"
		:disabled="disabled || loading"
		@click="handleClick"
	>
		<!-- Loading Spinner -->
		<span
			v-if="loading"
			class="spinner-border spinner-border-sm me-2"
			role="status"
			aria-hidden="true"
		></span>

		<!-- Icon (only show when not loading) -->
		<i v-if="icon && !loading" :class="iconClasses"></i>

		<!-- Button Text -->
		<span v-if="text">{{ text }}</span>

		<!-- Slot for additional content -->
		<slot></slot>
	</button>
</template>

<script setup>
import { computed } from "vue";

// Define props
const props = defineProps({
	text: {
		type: String,
		default: "",
	},
	icon: {
		type: String,
		default: "",
	},
	variant: {
		type: String,
		default: "primary",
		validator: (value) =>
			[
				"primary",
				"secondary",
				"success",
				"danger",
				"warning",
				"info",
				"light",
				"dark",
				"link",
				"outline-primary",
				"outline-secondary",
				"outline-success",
				"outline-danger",
				"outline-warning",
				"outline-info",
				"outline-light",
				"outline-dark",
			].includes(value),
	},
	size: {
		type: String,
		default: "",
		validator: (value) => ["", "sm", "lg"].includes(value),
	},
	loading: {
		type: Boolean,
		default: false,
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	block: {
		type: Boolean,
		default: false,
	},
});

// Define emits
const emit = defineEmits(["click"]);

// Computed classes for button
const buttonClasses = computed(() => {
	return [
		"btn",
		`btn-${props.variant}`,
		{
			[`btn-${props.size}`]: props.size,
			"w-100": props.block,
			disabled: props.loading,
		},
	];
});

// Computed classes for icon
const iconClasses = computed(() => {
	return [
		props.icon,
		{ "me-2": props.text }, // Add margin if there's text
	];
});

// Handle click event
const handleClick = (event) => {
	if (!props.loading && !props.disabled) {
		emit("click", event);
	}
};
</script>
