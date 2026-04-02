<style scoped>
.avatar-preview {
	width: 120px;
	height: 120px;
	object-fit: cover;
	border: 3px solid #dee2e6;
}
.bi-avatar {
	height: 128px;
	width: 128px;
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
										<i
											v-if="user.avatar == null && avatar == ''"
											class="bi bi-avatar"
										></i>
										<img
											v-else
											:src="user.avatar || avatar"
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
									id="display_name"
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
								<select id="role" v-model="user.role" class="form-select">
									<option value="">Select a role</option>
									<option
										v-for="option in optionRoles"
										:key="option.value"
										:value="option.value"
									>
										{{ option.label }}
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
										v-model="user.active"
										value="active"
									/>
									<label class="btn btn-outline-success" for="statusActive">
										<i class="bi bi-check-circle"></i> Active
									</label>

									<input
										type="radio"
										class="btn-check"
										id="statusInactive"
										v-model="user.active"
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
import { useSessionStore } from "@/store/session";
import { useUserStore } from "@/store/user";
import { User } from "@/types";

interface Props {
	id: number;
}
interface Option {
	label: string;
	value: string;
}
const avatar = ref<string>("");
const fileInput = ref<HTMLInputElement | null>(null);
const props = defineProps<Props>();
const selectedFile = ref<File | null>(null);
const selectedTag = ref<string>("");
const userStore = useUserStore();
const session = useSessionStore();
const user = ref<User | null>(null);

const optionRoles: Option[] = [
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
		selectedFile.value = file;
		const reader = new FileReader();
		reader.onload = (e) => {
			avatar.value = e.target?.result as string;
		};
		reader.readAsDataURL(file);
	}
};

const removeAvatar = () => {
	avatar.value = "";
	if (fileInput.value) {
		fileInput.value.value = "";
	}
};

const addTag = () => {
	if (user.value == null) {
		return;
	}
	if (selectedTag.value && !user.value.tags.includes(selectedTag.value)) {
		user.value.tags.push(selectedTag.value);
		selectedTag.value = "";
	}
};

const removeTag = (tag: string) => {
	if (user.value == null) {
		return;
	}
	const index = user.value.tags.indexOf(tag);
	if (index > -1) {
		user.value.tags.splice(index, 1);
	}
};

interface UserRequestPut {
	avatar: string | null;
}
const saveChanges = async () => {
	console.log("Saving user changes");
	let payload: UserRequestPut = {
		avatar: "",
	};
	if (selectedFile.value) {
		try {
			const formData = new FormData();
			formData.append("file", selectedFile.value);

			const url = session.urls?.api.avatar;
			if (!url) {
				console.log("empty avatar url");
				return;
			}
			const response = await fetch(url, {
				body: formData,
				method: "POST",
			});
			if (!response.ok) {
				throw new Error(`Upload failed: ${response.statusText}`);
			}

			const data = await response.json();
			payload.avatar = data.uri;
		} catch (error) {
			console.error("Failed to upload avatar", error);
		}
	}
	const u = user.value;
	if (!u) {
		console.log("empty user");
		return;
	}
	const url = u.uri;
	const response = await fetch(url, {
		method: "PUT",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(payload),
	});
	if (!response.ok) {
		const errorData = await response.json();
		throw new Error(
			errorData.message || `HTTP error! status: ${response.status}`,
		);
	}
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
	userStore.withAll().then((users) => {
		console.log("got users. looking for match", users, props.id);
		for (const u of users) {
			if (u.id == props.id) {
				user.value = u;
				console.log("User set to", u);
			}
		}
	});
});
</script>
