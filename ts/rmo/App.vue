<template>
	<div id="app">
		<RouterView />
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useHead } from "@vueuse/head";
import { router } from "@/rmo/router";

const count = ref<number>(0);
const message = ref<string>("hey");

const increment = (): void => {
	count.value++;
};

// Reactive head management
/*
useHead({
	title: computed(() => `Count: ${count.value} - My Vue App`),
	link: [
		{
			rel: "icon",
			type: "image/x-icon",
			href: "/favicon.ico",
		},
	],
});
*/
</script>
