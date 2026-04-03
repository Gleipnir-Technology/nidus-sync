<style scoped>
.bi-large {
	height: 48px;
	width: 48px;
}
</style>
<template>
	<!-- Main Content -->
	<main>
		<slot name="header"></slot>

		<!-- Services Section -->
		<section class="py-5">
			<div class="container">
				<h3 class="text-center mb-4">How Can We Help You Today?</h3>
				<div class="row g-4">
					<div class="col-md-4">
						<div class="card service-card h-100">
							<div class="card-body text-center">
								<div class="mb-3">
									<i class="bi bi-large bi-mosquito-color"></i>
								</div>
								<h4 class="card-title">Report Mosquito Nuisance</h4>
								<p class="card-text">
									Report areas with high adult mosquito activity causing
									discomfort or concern.
								</p>
								<RouterLink to="/nuisance" class="btn btn-primary mt-3"
									>Report Problem</RouterLink
								>
							</div>
						</div>
					</div>
					<div class="col-md-4">
						<div class="card service-card h-100">
							<div class="card-body text-center">
								<div class="mb-3">
									<i class="bi bi-large bi-pond-color"></i>
								</div>
								<h4 class="card-title">Report Standing Water</h4>
								<p class="card-text">
									Report any water that has been sitting for several days, where
									mosquitoes can live.
								</p>
								<RouterLink to="/water" class="btn btn-primary mt-3"
									>Report Source</RouterLink
								>
							</div>
						</div>
					</div>
					<div class="col-md-4">
						<div class="card service-card h-100">
							<div class="card-body text-center">
								<div class="mb-3">
									<i class="bi bi-large bi-check-report-color"></i>
								</div>
								<h4 class="card-title">Follow-up or Check Status</h4>
								<p class="card-text">
									Check on a previous request or view current mosquito activity
									in your area.
								</p>
								<RouterLink to="/status" class="btn btn-primary mt-3"
									>Get Status</RouterLink
								>
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>
	</main>
</template>

<script setup lang="ts">
import { ref } from "vue";

interface District {
	name: string;
	url_logo: string;
}
const district = ref<District | null>(null);
</script>
