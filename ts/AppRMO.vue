<template>
	<div id="app">
		<h1>{{ message }}</h1>
		<p>Count: {{ count }}</p>
		<button @click="increment">Increment</button>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useHead } from "@vueuse/head";

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
