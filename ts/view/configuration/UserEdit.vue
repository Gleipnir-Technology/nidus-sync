<style scoped>
.avatar-preview {
	width: 120px;
	height: 120px;
	object-fit: cover;
	border: 3px solid #dee2e6;
}

.btn-close-white {
	opacity: 0.8;
}

.btn-close-white:hover {
	opacity: 1;
}

.card {
	border-radius: 0.5rem;
}

.badge {
	font-size: 0.9rem;
	padding: 0.5rem 0.75rem;
	display: inline-flex;
	align-items: center;
}

pre {
	background-color: #f8f9fa;
	padding: 1rem;
	border-radius: 0.375rem;
	font-size: 0.875rem;
}
</style>
<template>
	<div class="container mt-4">
		<div class="row">
			<div class="col-lg-8 mx-auto">
				<div class="card shadow-sm">
					<div class="card-header bg-primary text-white">
						<h4 class="mb-0">User Configuration</h4>
					</div>
					<div v-if="user" class="card-body">
						<!-- Avatar Section -->
						<div class="row mb-4">
							<div class="col-md-12">
								<label class="form-label fw-bold">Avatar</label>
								<div class="d-flex align-items-center">
									<div class="position-relative">
										<img
											:src="user.avatar || defaultAvatar"
											alt="User Avatar"
											class="rounded-circle avatar-preview"
										/>
										<button
											class="btn btn-sm btn-danger position-absolute bottom-0 start-0 rounded-circle"
											@click="removeAvatar"
											type="button"
										>
											<i class="bi bi-trash"></i>
										</button>
										<button
											class="btn btn-sm btn-primary position-absolute bottom-0 end-0 rounded-circle"
											@click="triggerFileInput"
											type="button"
										>
											<i class="bi bi-camera"></i>
										</button>
									</div>
									<div class="ms-3">
										<input
											ref="fileInput"
											type="file"
											class="d-none"
											accept="image/*"
											@change="handleAvatarChange"
										/>
									</div>
								</div>
							</div>
						</div>

						<!-- Display Name -->
						<div class="row mb-3">
							<div class="col-md-6">
								<label for="displayName" class="form-label fw-bold">
									Display Name
								</label>
								<input
									id="displayName"
									v-model="user.display_name"
									type="text"
									class="form-control"
									placeholder="Enter display name"
								/>
							</div>

							<!-- Username (Read-only) -->
							<div class="col-md-6">
								<label for="username" class="form-label fw-bold">
									Username
								</label>
								<input
									id="username"
									v-model="user.username"
									type="text"
									class="form-control"
									readonly
									disabled
								/>
							</div>
						</div>

						<!-- User Role -->
						<div class="row mb-3">
							<div class="col-md-6">
								<label for="userRole" class="form-label fw-bold">
									User Role
								</label>
								<select id="userRole" v-model="user.role" class="form-select">
									<option value="">Select a role</option>
									<option
										v-for="role in availableRoles"
										:key="role.value"
										:value="role.value"
									>
										{{ role.label }}
									</option>
								</select>
							</div>

							<!-- User Status -->
							<div class="col-md-6">
								<label class="form-label fw-bold">User Status</label>
								<div class="btn-group w-100" role="group">
									<input
										type="radio"
										class="btn-check"
										id="statusActive"
										v-model="user.status"
										value="active"
									/>
									<label class="btn btn-outline-success" for="statusActive">
										<i class="bi bi-check-circle"></i> Active
									</label>

									<input
										type="radio"
										class="btn-check"
										id="statusInactive"
										v-model="user.status"
										value="inactive"
									/>
									<label class="btn btn-outline-secondary" for="statusInactive">
										<i class="bi bi-x-circle"></i> Inactive
									</label>
								</div>
							</div>
						</div>

						<!-- User Tags -->
						<div class="row mb-4">
							<div class="col-md-12">
								<label class="form-label fw-bold">User Tags</label>
								<div class="mb-2">
									<span
										v-for="tag in user.tags"
										:key="tag"
										class="badge bg-info text-dark me-2 mb-2"
									>
										{{ tag }}
										<button
											type="button"
											class="btn-close btn-close-white ms-2"
											@click="removeTag(tag)"
											style="font-size: 0.6rem"
										></button>
									</span>
									<span
										v-if="user.tags.length === 0"
										class="text-muted fst-italic"
									>
										No tags added
									</span>
								</div>
								<div class="input-group">
									<select v-model="selectedTag" class="form-select">
										<option value="">Select a tag</option>
										<option
											v-for="tag in availableTags"
											:key="tag"
											:value="tag"
											:disabled="user.tags.includes(tag)"
										>
											{{ tag }}
										</option>
									</select>
									<button
										class="btn btn-outline-primary"
										type="button"
										@click="addTag"
										:disabled="!selectedTag || user.tags.includes(selectedTag)"
									>
										<i class="bi bi-plus-lg"></i> Add Tag
									</button>
								</div>
							</div>
						</div>

						<!-- Action Buttons -->
						<div class="row">
							<div class="col-md-12">
								<hr />
								<div class="d-flex justify-content-end">
									<button
										class="btn btn-secondary me-2"
										type="button"
										@click="cancelChanges"
									>
										Cancel
									</button>
									<button
										class="btn btn-primary"
										type="button"
										@click="saveChanges"
									>
										<i class="bi bi-save"></i> Save Changes
									</button>
								</div>
							</div>
						</div>
					</div>
					<div v-else>
						<p>loading</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, defineComponent, onMounted, ref, reactive } from "vue";
import { useUsersStore } from "@/store/users";

interface User {
	avatar: string;
	displayName: string;
	username: string;
	role: string;
	status: "active" | "inactive";
	tags: string[];
}

interface Role {
	value: string;
	label: string;
}

interface Props {
	id: int;
}

const fileInput = ref<HTMLInputElement | null>(null);
const props = defineProps<Props>();
const selectedTag = ref<string>("");
const usersStore = useUsersStore();
const user = ref<User | null>(null);

const defaultAvatar =
	"https://via.placeholder.com/150/cccccc/666666?text=No+Avatar";

const availableRoles: Role[] = [
	{ value: "account-owner", label: "Account Owner" },
	{ value: "manager", label: "Manager" },
	{ value: "tech1", label: "Tech 1" },
	{ value: "tech2", label: "Tech 2" },
	{ value: "tech3", label: "Tech 3" },
];

const availableTags: string[] = [
	"warrant",
	"drone pilot",
	"certified",
	"supervisor",
	"field ops",
];

const triggerFileInput = () => {
	fileInput.value?.click();
};

const handleAvatarChange = (event: Event) => {
	const target = event.target as HTMLInputElement;
	const file = target.files?.[0];

	if (file) {
		const reader = new FileReader();
		reader.onload = (e) => {
			user.avatar = e.target?.result as string;
		};
		reader.readAsDataURL(file);
	}
};

const removeAvatar = () => {
	user.avatar = "";
	if (fileInput.value) {
		fileInput.value.value = "";
	}
};

const addTag = () => {
	if (selectedTag.value && !user.tags.includes(selectedTag.value)) {
		user.tags.push(selectedTag.value);
		selectedTag.value = "";
	}
};

const removeTag = (tag: string) => {
	const index = user.tags.indexOf(tag);
	if (index > -1) {
		user.tags.splice(index, 1);
	}
};

const saveChanges = () => {
	// Implement save logic here
	console.log("Saving user configuration:", user);
	alert("User configuration saved successfully!");
};

const cancelChanges = () => {
	// Implement cancel/reset logic here
	console.log("Canceling changes");
	if (
		confirm(
			"Are you sure you want to cancel? All unsaved changes will be lost.",
		)
	) {
		// Reset to original values or navigate away
		window.history.back();
	}
};
onMounted(() => {
	usersStore.fetchAll().then((users) => {
		for (const u of users) {
			if (u.id == props.id) {
				user.value = u;
				console.log("User set to", u);
			}
		}
	});
});
</script>
