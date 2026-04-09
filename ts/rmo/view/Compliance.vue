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
		<router-view v-slot="{ Component }">
			<component
				:is="Component"
				:district="district"
				@doComments="doComments"
				@doImages="doImages"
				@doLocator="doLocator"
			/>
		</router-view>
	</template>
	<template v-else>
		<p>loading {{ slug }}...</p>
	</template>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { computedAsync } from "@vueuse/core";

import type { Image } from "@/components/ImageUpload.vue";
import { useStoreDistrict } from "@/rmo/store/district";
import Intro from "@/rmo/content/compliance/Intro.vue";
import type { District } from "@/type/api";
import { Locator } from "@/type/map";

interface Props {
	slug: string;
}

const districtStore = useStoreDistrict();

const props = defineProps<Props>();
const district = computedAsync(async (): Promise<District | undefined> => {
	const districts = await districtStore.list();
	return districts.find((district: District) => district.slug == props.slug);
});
function doComments(comments: string) {
	console.log("comments", comments);
}
function doImages(images: Image[]) {
	console.log("images", images);
}
function doLocator(locator: Locator | null) {
	console.log("locator", locator);
}
</script>
