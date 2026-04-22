<template>
	<p>signing out...</p>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { apiClient } from "@/client";
import { router } from "@/route/config";
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
onMounted(() => {
	session.signout().then(() => {
		router.push("/signin");
	});
});
</script>
