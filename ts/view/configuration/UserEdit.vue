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
					<div v-if="userChanges" class="card-body">
						<!-- Avatar Section -->
						<div class="row mb-4">
							<div class="col-md-12">
								<label class="form-label fw-bold">Avatar</label>
								<div class="d-flex align-items-center">
									<div class="position-relative">
										<i
											v-if="
												userChanges.avatar == null || userChanges.avatar == ''
											"
											class="bi bi-avatar"
										></i>
										<img
											v-else
											:src="userChanges.avatar"
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
									v-model="userChanges.display_name"
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
									v-model="userChanges.username"
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
								<select
									v-if="userChanges && userChanges.role == 'root'"
									disabled
									id="role"
									class="form-select"
								>
									<option value="root">Root</option>
								</select>
								<select
									id="role"
									v-else
									v-model="userChanges.role"
									class="form-select"
								>
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
										type="checkbox"
										class="btn-check"
										id="isActive"
										v-model="userChanges.is_active"
										value="active"
									/>
									<label class="btn btn-outline-success" for="isActive">
										<i class="bi bi-check-circle"></i> Active
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
										v-for="tag in userChanges.tags"
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
										v-if="userChanges.tags.length === 0"
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
											:disabled="userChanges.tags.includes(tag)"
										>
											{{ tag }}
										</option>
									</select>
									<button
										class="btn btn-outline-primary"
										type="button"
										@click="addTag"
										:disabled="
											!selectedTag || userChanges.tags.includes(selectedTag)
										"
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
										:disabled="isSaving"
										type="button"
										@click="cancelChanges"
									>
										Cancel
									</button>
									<button
										class="btn btn-primary"
										:disabled="isSaving"
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
import { onMounted, ref, toRaw } from "vue";
import { useSessionStore } from "@/store/session";
import { useUserStore } from "@/store/user";
import { User } from "@/type/api";

interface Props {
	id: number;
}
interface Option {
	label: string;
	value: string;
}
export interface UserChanges {
	avatar: string;
	display_name: string;
	is_active: boolean;
	role: string;
	tags: string[];
	username: string;
}
const fileInput = ref<HTMLInputElement | null>(null);
const isSaving = ref<boolean>(false);
const props = defineProps<Props>();
const selectedFile = ref<File | null>(null);
const selectedTag = ref<string>("");
const userStore = useUserStore();
const session = useSessionStore();
const userChanges = ref<UserChanges>();

const optionRoles: Option[] = [
	{ value: "account-owner", label: "Account Owner" },
	{ value: "manager", label: "Manager" },
	{ value: "tech1", label: "Tech 1" },
	{ value: "tech2", label: "Tech 2" },
	{ value: "tech3", label: "Tech 3" },
];

const availableTags: string[] = ["warrant", "drone pilot"];

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
			if (userChanges.value == undefined) {
				console.log("can't update userChanges, it's undefined");
				return;
			}
			userChanges.value.avatar = e.target?.result as string;
		};
		reader.readAsDataURL(file);
	}
};

const removeAvatar = () => {
	if (userChanges.value == undefined) {
		return;
	}
	userChanges.value.avatar = "";
	if (fileInput.value) {
		fileInput.value.value = "";
	}
};

const addTag = () => {
	if (userChanges.value == null) {
		return;
	}
	if (
		selectedTag.value &&
		!userChanges.value.tags.includes(selectedTag.value)
	) {
		userChanges.value.tags.push(selectedTag.value);
		selectedTag.value = "";
	}
};

const removeTag = (tag: string) => {
	if (userChanges.value == null) {
		return;
	}
	const index = userChanges.value.tags.indexOf(tag);
	if (index > -1) {
		userChanges.value.tags.splice(index, 1);
	}
};

interface UserRequestPut {
	avatar?: string | null;
	display_name?: string;
	is_active?: boolean;
	role?: string;
	tags?: string[];
}
const saveChanges = async () => {
	const uc = userChanges.value;
	if (!uc) {
		console.log("empty user changes");
		return;
	}
	console.log("Saving user changes");
	isSaving.value = true;
	const all_users = await userStore.withAll();
	const u = all_users.find((u: User) => u.id == props.id);
	if (!u) {
		console.log("no matching user");
		isSaving.value = false;
		return;
	}
	let payload: UserRequestPut = {};
	if (uc.avatar != u.avatar) {
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
				isSaving.value = false;
				return;
			}
		} else if (!uc.avatar) {
			payload.avatar = null;
		}
	}
	if (uc.display_name != u.display_name) {
		payload.display_name = uc.display_name;
	}
	if (uc.is_active != u.is_active) {
		payload.is_active = uc.is_active;
	}
	if (uc.role != u.role) {
		payload.role = uc.role;
	}
	if (uc.tags != u.tags) {
		payload.tags = uc.tags;
	}
	if (Object.keys(payload).length === 0) {
		console.log("refusing to make empty changes");
		isSaving.value = false;
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
	isSaving.value = false;
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
				userChanges.value = {
					avatar: u.avatar,
					display_name: u.display_name,
					is_active: u.is_active,
					role: u.role,
					tags: structuredClone(toRaw(u.tags)),
					username: u.username,
				};
				console.log("User set to", u);
			}
		}
	});
});
</script>
