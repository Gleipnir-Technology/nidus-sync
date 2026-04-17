<template>
	<router-view />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useSessionStore } from "@/store/session";
import { Session } from "@/type/api";
import { router } from "@/router";

const session = useSessionStore();
onMounted(() => {
	session
		.get()
		.then((session: Session) => {
			console.log("session loaded by AppSync", session);
		})
		.catch((e) => {
			console.log("root session not loaded", e);
			router.push("/signin");
		});
});
</script>
