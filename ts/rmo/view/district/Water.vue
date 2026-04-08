<style scoped></style>
<template>
	<Water :slug="slug">
		<template #header>
			<!-- Introduction Section -->
			<HeaderDistrict :district="district" />
		</template>
	</Water>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { computedAsync } from "@vueuse/core";
import Water from "@/rmo/content/Water.vue";
import type { District } from "@/type/api";
import { useStoreDistrict } from "@/rmo/store/district";
import HeaderDistrict from "@/components/HeaderDistrict.vue";

interface Props {
	slug: string;
}
const props = defineProps<Props>();
const districtStore = useStoreDistrict();

const district = computedAsync(async (): Promise<District | undefined> => {
	const districts = await districtStore.list();
	return districts.find((district: District) => district.slug == props.slug);
});
</script>
