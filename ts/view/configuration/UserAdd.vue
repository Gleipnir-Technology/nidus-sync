<template>
	<div class="container py-5">
		<div class="row justify-content-center">
			<div class="col-lg-8">
				<div class="card shadow">
					<div class="card-header bg-white">
						<div class="d-flex justify-content-between align-items-center">
							<h3 class="mb-0">Add New User</h3>
						</div>
					</div>
					<div class="card-body">
						<form
							id="addUserForm"
							class="needs-validation"
							novalidate
							@submit.prevent="handleSubmit"
						>
							<!-- Full Name -->
							<div class="mb-3">
								<label for="fullName" class="form-label required-field">
									Full Name
								</label>
								<input
									type="text"
									class="form-control"
									id="fullName"
									v-model="formData.fullName"
									:class="{ 'is-invalid': isSubmitted && !formData.fullName }"
									required
								/>
								<div class="invalid-feedback">
									Please provide the user's full name.
								</div>
							</div>

							<!-- Email Address -->
							<div class="mb-3">
								<label for="emailAddress" class="form-label required-field">
									Email Address
								</label>
								<input
									type="email"
									class="form-control"
									id="emailAddress"
									v-model="formData.emailAddress"
									:class="{ 'is-invalid': isSubmitted && !isValidEmail }"
									required
								/>
								<div class="invalid-feedback">
									Please provide a valid email address.
								</div>
								<div class="form-text">
									An invitation will be sent to this email address.
								</div>
							</div>

							<!-- Username -->
							<div class="mb-3">
								<label for="username" class="form-label required-field">
									Username
								</label>
								<input
									type="text"
									class="form-control"
									id="username"
									v-model="formData.username"
									:class="{ 'is-invalid': isSubmitted && !formData.username }"
									required
								/>
								<div class="invalid-feedback">Please provide a username.</div>
								<div class="form-text">
									Username must be unique and contain only letters, numbers, and
									underscores.
								</div>
							</div>

							<div class="row">
								<!-- Role -->
								<div class="col-md-6 mb-3">
									<label for="userRole" class="form-label required-field">
										Role
									</label>
									<select
										class="form-select"
										id="userRole"
										v-model="formData.userRole"
										:class="{ 'is-invalid': isSubmitted && !formData.userRole }"
										required
									>
										<option value="" disabled>Select a role</option>
										<option value="lead">Lead</option>
										<option value="technician">Technician</option>
										<option value="administrator">Administrator</option>
									</select>
									<div class="invalid-feedback">Please select a role.</div>
								</div>

								<!-- Initial Status -->
								<div class="col-md-6 mb-3">
									<label for="initialStatus" class="form-label">
										Initial Status
									</label>
									<select
										class="form-select"
										id="initialStatus"
										v-model="formData.initialStatus"
									>
										<option value="invited">Invited</option>
										<option value="active">Active</option>
									</select>
								</div>
							</div>

							<!-- Permissions -->
							<div class="mb-4">
								<label class="form-label d-block">Permissions</label>
								<div class="form-check form-switch">
									<input
										class="form-check-input switch-lg"
										type="checkbox"
										id="serveWarrants"
										v-model="formData.serveWarrants"
									/>
									<label class="form-check-label" for="serveWarrants">
										Can serve warrants
									</label>
								</div>
							</div>

							<!-- Send welcome email checkbox -->
							<div class="mb-4">
								<div class="form-check">
									<input
										class="form-check-input"
										type="checkbox"
										id="sendWelcomeEmail"
										v-model="formData.sendWelcomeEmail"
									/>
									<label class="form-check-label" for="sendWelcomeEmail">
										Send welcome email with login instructions
									</label>
								</div>
							</div>

							<hr />

							<!-- Form actions -->
							<div class="d-flex justify-content-end gap-2">
								<button
									type="button"
									class="btn btn-secondary"
									@click="handleCancel"
								>
									Cancel
								</button>
								<button
									type="submit"
									class="btn btn-primary"
									:disabled="isLoading"
								>
									<i class="bi bi-person-plus me-1"></i>
									{{ isLoading ? "Adding..." : "Add User" }}
								</button>
							</div>
						</form>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";

interface UserFormData {
	fullName: string;
	emailAddress: string;
	username: string;
	userRole: string;
	initialStatus: string;
	serveWarrants: boolean;
	sendWelcomeEmail: boolean;
}

interface Props {
	cancelUrl?: string;
}

const props = withDefaults(defineProps<Props>(), {
	cancelUrl: "/configuration/user",
});

const emit = defineEmits<{
	submit: [formData: UserFormData];
	cancel: [];
}>();

const router = useRouter();

const formData = ref<UserFormData>({
	fullName: "",
	emailAddress: "",
	username: "",
	userRole: "",
	initialStatus: "invited",
	serveWarrants: false,
	sendWelcomeEmail: true,
});

const isSubmitted = ref(false);
const isLoading = ref(false);

const isValidEmail = computed(() => {
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return emailRegex.test(formData.value.emailAddress);
});

const isFormValid = computed(() => {
	return (
		formData.value.fullName.trim() !== "" &&
		isValidEmail.value &&
		formData.value.username.trim() !== "" &&
		formData.value.userRole !== ""
	);
});

const handleSubmit = async () => {
	isSubmitted.value = true;

	if (!isFormValid.value) {
		return;
	}

	try {
		isLoading.value = true;
		emit("submit", formData.value);

		// Example API call (uncomment and adjust as needed):
		// await api.post('/api/users', formData.value);
		// router.push(props.cancelUrl);
	} catch (error) {
		console.error("Error adding user:", error);
		// Handle error (e.g., show toast notification)
	} finally {
		isLoading.value = false;
	}
};

const handleCancel = () => {
	emit("cancel");
	router.push(props.cancelUrl);
};
</script>

<style scoped>
.form-check-input.switch-lg {
	width: 3em;
	height: 1.5em;
}

.required-field::after {
	content: " *";
	color: red;
}
</style>
