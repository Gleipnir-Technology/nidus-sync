<template>
	<div class="container-fluid px-3 py-3">
		<HeaderCompliance :district="district" />
		<ProgressBarCompliance :step="1" />

		<!-- Main Content -->
		<main>
			<h2 class="h4 mb-3">Action requested for this property</h2>

			<div class="alert alert-warning" role="alert">
				<p class="mb-2">
					The {{ district.name }} has identified a possible mosquito breeding
					source at this property during a recent inspection.
				</p>
				<p class="mb-0">
					Please confirm the property address and provide any photos, access
					details, and contact information to help us review and resolve this
					issue quickly.
				</p>
			</div>

			<div class="d-grid mt-4">
				<RouterLink
					class="btn btn-primary btn-lg"
					:to="routes.ComplianceAddress(props.publicID)"
				>
					Get Started</RouterLink
				>
			</div>
		</main>
	</div>
</template>
<script setup lang="ts">
import type { District } from "@/type/api";
import HeaderCompliance from "@/rmo/components/HeaderCompliance.vue";
import ProgressBarCompliance from "@/rmo/components/ProgressBarCompliance.vue";
import { useRoutes } from "@/rmo/route/use";
interface Props {
	district: District;
	publicID: string;
}
const props = defineProps<Props>();
const routes = useRoutes();
</script>
