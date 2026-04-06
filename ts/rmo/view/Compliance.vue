<style scoped>
body {
	background-color: #f8f9fa;
	min-height: 100vh;
	display: flex;
	flex-direction: column;
}

body > .container-fluid {
	flex: 1;
}

.progress-bar {
	background-color: #0d6efd;
	transition: width 0.3s ease;
}
</style>
<template>
	<template v-if="district">
		<Intro :district="district" />
	</template>
	<template v-else>
		<p>loading {{ slug }}...</p>
	</template>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { computedAsync } from "@vueuse/core";

import { useDistrictStore } from "@/rmo/store/district";
import Intro from "@/rmo/content/compliance/Intro.vue";
import type { District } from "@/type/api";
interface Props {
	slug: string;
}

const districtStore = useDistrictStore();

const props = defineProps<Props>();
const district = computedAsync(async (): Promise<District | undefined> => {
	const districts = await districtStore.get();
	return districts.find((district: District) => district.slug == props.slug);
});
</script>
