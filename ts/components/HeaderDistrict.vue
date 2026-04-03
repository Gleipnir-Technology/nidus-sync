<style scoped>
.district-logo {
	display: block;
	margin-left: auto;
	margin-right: auto;
	max-height: 88px;
	width: auto;
}
</style>
<template>
	<header
		class="navbar navbar-expand-lg navbar-dark text-light bg-primary shadow-sm"
	>
		<div class="container">
			<div class="d-flex align-items-center">
				<template v-if="district">
					<a class="navbar-brand d-flex align-items-center" href="/">
						<img
							:src="district.url_logo"
							style="height: 48px"
							alt="District logo"
						/>
					</a>
					<h1 class="mb-0 ms-3">{{ district.name }}</h1>
				</template>
				<template v-else>
					<h1>Loading...</h1>
				</template>
			</div>
		</div>
	</header>
</template>
<script setup lang="ts">
import type { District } from "@/rmo/type";
interface Props {
	district?: District;
}
const props = defineProps<Props>();
</script>
