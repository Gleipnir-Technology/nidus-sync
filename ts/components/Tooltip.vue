<template>
	<span
		ref="tooltipElement"
		data-bs-toggle="tooltip"
		:data-bs-placement="placement"
		:title="title"
	>
		<slot></slot>
	</span>
</template>

<script>
import { Tooltip } from "bootstrap";

export default {
	name: "BsTooltip",

	props: {
		title: {
			type: String,
			required: true,
		},
		placement: {
			type: String,
			default: "top",
			validator(value) {
				return ["top", "bottom", "left", "right"].includes(value);
			},
		},
	},

	data() {
		return {
			tooltip: null,
		};
	},

	mounted() {
		this.tooltip = new Tooltip(this.$refs.tooltipElement);
	},

	beforeUnmount() {
		if (this.tooltip) {
			this.tooltip.dispose();
		}
	},

	watch: {
		title(newTitle) {
			if (this.tooltip) {
				this.tooltip.dispose();
				this.tooltip = new Tooltip(this.$refs.tooltipElement);
			}
		},
	},
};
</script>
