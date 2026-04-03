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
	<template v-if="district">
		<section class="py-5 bg-primary text-white">
			<div class="container">
				<div class="row justify-content-center">
					<div class="col-lg-10">
						<h2 class="text-center mb-4">Report a Mosquito Problem</h2>
						<p class="lead text-center">
							Submit a report to help reduce mosquito activity in your
							neighborhood.
						</p>
						<p class="lead text-center">
							Report Mosquitoes Online works with local mosquito control
							agencies to receive public reports.
						</p>
						<p class="lead text-center">
							For this area, mosquito control services are provided by
						</p>
						<h3 class="text-center">{{ district.name }}</h3>
						<img class="district-logo" :src="district.url_logo" />
					</div>
				</div>
			</div>
		</section>
	</template>
	<template v-else>
		<p>loading district...</p>
	</template>
</template>
<script setup lang="ts">
import type { District } from "@/rmo/type";
interface Props {
	district?: District;
}
const props = defineProps<Props>();
</script>
