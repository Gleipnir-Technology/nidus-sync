<template>
	<div id="app">
		<RouterView />
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
//import { useHead } from "@vueuse/head";
import { router } from "@/rmo/router";
import { useStoreDistrict } from "@/rmo/store/district";
import type { District } from "@/type/api";

const district = useStoreDistrict();
const count = ref<number>(0);
const message = ref<string>("hey");

const increment = (): void => {
	count.value++;
};

onMounted(() => {
	district
		.list()
		.then((districts: District[]) => {
			console.log("got districts");
		})
		.catch((e: Error) => {
			console.error("Failed to get districts", e);
		});
});
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
