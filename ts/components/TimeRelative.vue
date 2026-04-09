<template>
	<span class="time-relative">{{ relativeTime }}</span>
</template>

<script lang="ts">
import { defineComponent } from "vue";

/**
 * TimeRelative component that displays relative time
 * Usage: <TimeRelative time="2024-01-01T12:00:00Z" />
 */
export default defineComponent({
	name: "TimeRelative",

	props: {
		time: {
			type: Date,
		},
	},

	data() {
		return {
			relativeTime: "" as string,
			intervalId: null as number | null,
		};
	},

	watch: {
		time: {
			immediate: true,
			handler() {
				this.updateTime();
			},
		},
	},

	mounted() {
		this.updateTime();
		// Update every 60 seconds
		this.intervalId = window.setInterval(() => this.updateTime(), 60000);
	},

	unmounted() {
		if (this.intervalId !== null) {
			clearInterval(this.intervalId);
		}
	},

	methods: {
		updateTime(): void {
			if (this.time) {
				this.relativeTime = this.formatRelativeTime(this.time);
			} else {
				this.relativeTime = "";
			}
		},

		formatRelativeTime(timestamp: Date): string {
			const now = new Date();
			const date = new Date(timestamp);
			const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

			// Time units in seconds
			const minute = 60;
			const hour = minute * 60;
			const day = hour * 24;
			const week = day * 7;
			const month = day * 30;
			const year = day * 365;

			if (diffInSeconds < minute) {
				return "just now";
			} else if (diffInSeconds < hour) {
				const minutes = Math.floor(diffInSeconds / minute);
				return `${minutes} ${minutes === 1 ? "min" : "min"} ago`;
			} else if (diffInSeconds < day) {
				const hours = Math.floor(diffInSeconds / hour);
				return `${hours} ${hours === 1 ? "hour" : "hours"} ago`;
			} else if (diffInSeconds < week) {
				const days = Math.floor(diffInSeconds / day);
				return `${days} ${days === 1 ? "day" : "days"} ago`;
			} else if (diffInSeconds < month) {
				const weeks = Math.floor(diffInSeconds / week);
				return `${weeks} ${weeks === 1 ? "week" : "weeks"} ago`;
			} else if (diffInSeconds < year) {
				const months = Math.floor(diffInSeconds / month);
				return `${months} ${months === 1 ? "month" : "months"} ago`;
			} else {
				const years = Math.floor(diffInSeconds / year);
				return `${years} ${years === 1 ? "year" : "years"} ago`;
			}
		},
	},
});
</script>

<style scoped>
.time-relative {
	/* Add your styles here */
}
</style>
