<template>
	<p>signing out...</p>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { apiClient } from "@/client";
import { router } from "@/router";
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
onMounted(() => {
	session.signout().then(() => {
		apiClient.setAuthenticated(false);
		router.push("/signin");
	});
});
</script>
