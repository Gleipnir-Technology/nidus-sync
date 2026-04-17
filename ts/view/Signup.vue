<style scoped>
.register-container {
	max-width: 900px;
	margin: 0 auto;
}
.register-box {
	box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
	border-radius: 10px;
	overflow: hidden;
}
.register-form-section {
	padding: 40px;
}
.register-info-section {
	padding: 40px;
	background-color: #f8f9fa;
}
.register-header {
	margin-bottom: 25px;
}
.logo-area {
	text-align: center;
	margin-bottom: 30px;
}
</style>
<template>
	<div
		class="container min-vh-100 d-flex align-items-center justify-content-center py-5"
	>
		<div class="register-container">
			<div class="row register-box g-0">
				<!-- Left side: Registration Form -->
				<div class="col-md-6 register-form-section">
					<div class="register-header">
						<h2>Create an Account</h2>
						<p class="text-muted">Join us today to get started</p>
					</div>

					<!-- Error message display -->
					<div v-if="errorMessage" class="alert alert-danger" role="alert">
						{{ errorMessage }}
					</div>

					<div class="mb-3">
						<label for="name" class="form-label">Display Name</label>
						<input
							type="text"
							class="form-control"
							:disalbel="isLoading"
							name="name"
							required
							v-model="formData.name"
						/>
					</div>

					<div class="mb-3">
						<label for="username" class="form-label">Username</label>
						<input
							class="form-control"
							:disabled="isLoading"
							name="username"
							required
							type="username"
							v-model="formData.username"
						/>
					</div>

					<div class="mb-3">
						<label for="password" class="form-label">Password</label>
						<input
							type="password"
							class="form-control"
							name="password"
							required
							v-model="formData.password"
						/>
					</div>

					<div class="mb-3 form-check">
						<input
							type="checkbox"
							class="form-check-input"
							:disabled="isLoading"
							name="terms"
							required
							v-model="formData.terms"
						/>
						<label class="form-check-label" for="terms"
							>I agree to the <a href="#">Terms of Service</a> and
							<a href="#">Privacy Policy</a></label
						>
					</div>

					<div class="d-grid gap-2">
						<ButtonLoading
							@click="handleRegister()"
							text="Register"
							variant="primary"
							:loading="isLoading"
							:disabled="
								!(formData.name && formData.username && formData.password)
							"
						/>
					</div>

					<div class="mt-3 text-center">
						<p>
							Already have an account?
							<RouterLink to="/signin">Sign in</RouterLink>
						</p>
					</div>
				</div>

				<!-- Right side: Account Information -->
				<div class="col-md-6 register-info-section">
					<div>
						<div class="logo-area">
							<img src="/static/img/nidus-logo-256-transparent.png" />
						</div>

						<div class="mb-4">
							<h5>Who should register?</h5>
							<p>
								This platform is designed for professionals who need to manage
								projects and collaborate with team members. Whether you're a
								freelancer, small business owner, or part of a larger
								organization, our tools can help streamline your workflow.
							</p>
						</div>

						<div class="mb-4">
							<h5>What happens after registration?</h5>
							<p>
								After you register with your email, you'll receive a
								confirmation message with instructions to complete your account
								setup. You'll then have access to all features and can customize
								your workspace based on your specific needs.
							</p>
						</div>

						<div class="mb-3">
							<small class="text-muted"
								>For any questions about account types or registration, please
								contact our support team at support@yourproduct.com</small
							>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

import ButtonLoading from "@/components/common/ButtonLoading.vue";
const router = useRouter();

// Form data
const formData = ref({
	name: "",
	username: "",
	password: "",
	terms: false,
});

// UI state
const isLoading = ref(false);
const errorMessage = ref("");

// Handle form submission
const handleRegister = async () => {
	console.log("registering", formData.value.terms);
	// Clear previous errors
	errorMessage.value = "";

	// Client-side validation
	if (
		!formData.value.name ||
		!formData.value.username ||
		!formData.value.password
	) {
		errorMessage.value = "Please fill in all required fields";
		return;
	}

	if (!formData.value.terms) {
		errorMessage.value =
			"You must agree to the Terms of Service and Privacy Policy";
		return;
	}

	isLoading.value = true;

	try {
		const response = await fetch("/api/signup", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				username: formData.value.username,
				name: formData.value.name,
				password: formData.value.password,
				terms: formData.value.terms,
			}),
		});

		if (!response.ok) {
			const errorData = await response
				.json()
				.catch(() => ({ error: "Registration failed" }));
			throw new Error(
				errorData.error || `HTTP error! status: ${response.status}`,
			);
		}

		// Redirect to the path returned by the backend (default: "/")
		router.push("/_/dash");
	} catch (error) {
		console.error("Registration error:", error);
		errorMessage.value =
			error instanceof Error
				? error.message
				: "An unexpected error occurred during registration";
	} finally {
		isLoading.value = false;
	}
};
</script>
