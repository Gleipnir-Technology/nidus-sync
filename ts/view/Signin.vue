<style scoped>
.login-container {
	max-width: 900px;
	margin: 0 auto;
}
.login-box {
	box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
	border-radius: 10px;
	overflow: hidden;
}
.login-form-section {
	padding: 40px;
}
.product-info-section {
	padding: 40px;
	background-color: #f8f9fa;
}
.login-header {
	margin-bottom: 25px;
}
</style>

<template>
	<div
		class="container min-vh-100 d-flex align-items-center justify-content-center py-5"
	>
		<div class="login-container">
			<div class="row login-box g-0">
				<!-- Left side: Login Form -->
				<div class="col-md-6 login-form-section">
					<div class="login-header">
						<h2>Welcome Back</h2>
						<p class="text-muted">Please enter your credentials</p>
					</div>

					<input type="hidden" name="next" value="none" />
					<div class="mb-3">
						<label for="username" class="form-label">Username</label>
						<input
							type="text"
							class="form-control"
							name="username"
							required
							v-model="username"
						/>
					</div>

					<div class="mb-3">
						<label for="password" class="form-label">Password</label>
						<input
							type="password"
							class="form-control"
							name="password"
							v-model="password"
							required
						/>
					</div>

					<!--
							<div class="alert alert-danger" role="alert">
								The credentials you provided weren't recognized.
							</div>
						-->

					<div class="d-grid gap-2">
						<ButtonLoading
							@click="doLogin()"
							:loading="loading"
							text="Login"
							variant="primary"
						/>
					</div>

					<div class="mt-3 text-center">
						<p>Don't have an account? <a href="/signup">Sign up</a></p>
						<a href="forgot-password.html">Forgot password?</a>
					</div>

					<div class="mt-3 text-center" v-if="error">
						<div class="alert alert-danger">{{ error }}</div>
					</div>
				</div>

				<!-- Right side: Product Information -->
				<div class="col-md-6 product-info-section">
					<div>
						<img src="/static/img/nidus-logo-256-transparent.png" />
						<h2>Nidus Sync</h2>
						<p class="lead mb-4">
							All your field data, sync'd to all your techs
						</p>

						<div class="mb-4">
							<p>Something intelligent and intriguing</p>
						</div>

						<div class="mb-4">
							<h5>Key Features</h5>
							<ul>
								<li>Works with <b>Fieldseeker</b></li>
								<li>Works <i>with</i> Fieldseeker</li>
								<li><b>Works</b> with Fieldseeker</li>
							</ul>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { apiClient } from "@/client";
import ButtonLoading from "@/components/common/ButtonLoading.vue";
import { router } from "@/router";

const error = ref<string>("");
const loading = ref<boolean>(false);
const password = ref<string>("");
const username = ref<string>("");
async function doLogin() {
	loading.value = true;
	try {
		const resp = await apiClient.JSONPost("/api/signin", {
			password: password.value,
			username: username.value,
		});
		router.push("/");
	} catch (e) {
		console.log("login failed", e);
		error.value = `Login failed: ${e}`;
	} finally {
		loading.value = false;
	}
}
</script>
