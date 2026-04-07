<template>
	<header class="text-center mb-4 pb-3 border-bottom">
		<div class="d-flex align-items-center justify-content-center mb-2">
			<img
				:src="district.url_logo"
				:alt="district.name + ' logo'"
				class="me-2"
				style="height: 40px; width: auto"
			/>
			<h1 class="h5 mb-0">{{ district.name }}</h1>
		</div>
		<div class="text-muted small" v-if="district.phone_office">
			<i class="bi bi-telephone"></i> {{ district.phone_office }}
		</div>
	</header>
</template>

<script setup lang="ts">
import type { District } from "@/type/api";
interface Props {
	district: District;
}
const props = defineProps<Props>();
</script>
