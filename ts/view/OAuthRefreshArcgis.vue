<template>
	<div
		class="container min-vh-100 d-flex align-items-center justify-content-center py-5"
	>
		<div class="connect-container">
			<div class="connect-box">
				<div class="connect-header">
					<h1>Connect Your ArcGIS Account</h1>
					<p class="text-muted">Link your data to get started</p>
				</div>

				<div class="connect-content">
					<p>
						To provide you with the best experience, we need to connect to your
						ArcGIS account. This allows us to securely access and visualize your
						spatial data within our platform.
					</p>

					<div class="steps-container">
						<h4>What to expect:</h4>

						<div v-for="step in steps" :key="step.number" class="step">
							<h5>{{ step.number }}. {{ step.title }}</h5>
							<p>{{ step.description }}</p>
						</div>
					</div>

					<div class="alert alert-info">
						<strong>Note:</strong> You'll need an active ArcGIS Online account
						or ArcGIS Enterprise account to proceed. If you don't have one, you
						can
						<a
							href="https://www.arcgis.com/home/signin.html"
							target="_blank"
							rel="noopener noreferrer"
						>
							create an ArcGIS account here </a
						>.
					</div>

					<p>By connecting your ArcGIS account, you'll be able to:</p>
					<ul>
						<li v-for="benefit in benefits" :key="benefit">
							{{ benefit }}
						</li>
					</ul>

					<div class="text-center connect-btn">
						<a :href="oauthUrl" class="btn btn-primary btn-lg">
							Connect to ArcGIS
						</a>
						<p class="mt-2 text-muted">
							<small
								>You can disconnect your account at any time in settings</small
							>
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";

interface Step {
	number: number;
	title: string;
	description: string;
}

// Define steps
const steps = ref<Step[]>([
	{
		number: 1,
		title: "Secure Authentication",
		description:
			'When you click the "Connect to ArcGIS" button below, you\'ll be redirected to the official ArcGIS login page. This connection is secure and uses OAuth 2.0 protocol.',
	},
	{
		number: 2,
		title: "Grant Permissions",
		description:
			"After logging in with your ArcGIS credentials, you'll be asked to approve permissions for our application to access your data. We only request access to what's needed for the platform to function.",
	},
	{
		number: 3,
		title: "Return to Platform",
		description:
			"Once authentication is complete, you'll be automatically redirected back to our platform where your data will be available to work with.",
	},
]);

// Define benefits
const benefits = ref<string[]>([
	"Access and visualize your spatial data",
	"Perform advanced analysis using our integrated tools",
	"Share results with team members securely",
	"Keep your data synchronized across platforms",
]);

// OAuth URL
const oauthUrl = computed(() => "/arcgis/oauth/begin");
</script>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
	name: "ArcGISConnect",
});
</script>

<style scoped>
.connect-container {
	max-width: 800px;
	margin: 0 auto;
}

.connect-box {
	box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
	border-radius: 10px;
	padding: 40px;
	background-color: #fff;
}

.connect-header {
	margin-bottom: 25px;
	text-align: center;
}

.steps-container {
	margin: 30px 0;
}

.step {
	margin-bottom: 20px;
	padding: 15px;
	border-left: 3px solid #0d6efd;
	background-color: #f8f9fa;
}

.connect-btn {
	margin-top: 30px;
}
</style>
