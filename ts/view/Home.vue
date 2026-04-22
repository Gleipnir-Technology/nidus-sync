<template>
	<p>loading home...</p>
</template>
<script setup lang="ts">
import { onMounted } from "vue";
import { useSessionStore } from "@/store/session";
import { Session } from "@/type/api";
import { router } from "@/route/config";

const session = useSessionStore();
onMounted(() => {
	session
		.get()
		.then((session: Session) => {
			console.log("hit home with a valid session, going to dash");
			router.push("/_/dash");
		})
		.catch((e) => {
			console.log("hit home with no session, going to signin");
			router.push("/signin");
		});
});
</script>
