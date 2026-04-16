<template>
	<p>signing out...</p>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { router } from "@/router";
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
onMounted(() => {
	session.signout().then(() => {
		router.push("/signin");
	});
});
</script>
