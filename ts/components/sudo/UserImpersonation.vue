<template>
	<!-- User Impersonation Section -->
	<h2 id="user-impersonation">
		<i class="bi bi-person-badge"></i> User Impersonation
	</h2>
	<div class="card">
		<div class="card-header bg-dark text-white">
			<i class="bi bi-people"></i> Impersonate User
		</div>
		<div class="card-body">
			<template v-if="session.impersonating">
				<h1>You're impersonating</h1>
				<p>{{ session.impersonating }}</p>
				<button class="btn btn-primary" @click="doImpersonationEnd">
					End Impersonation
				</button>
			</template>
			<template v-else>
				<div class="row mb-3">
					<div class="col-md-6">
						<label for="userSearch" class="form-label">Search Users</label>
						<UserSelector
							v-model="selectedUser"
							label="Choose a user"
							placeholder="Select a user..."
						/>
					</div>
					<div class="col-md-6">
						<label for="userRole" class="form-label">Filter by Role</label>
						<select class="form-select" id="userRole">
							<option value="">All Roles</option>
							<option value="admin">Admin</option>
							<option value="user">Standard User</option>
							<option value="support">Support</option>
							<option value="premium">Premium User</option>
						</select>
					</div>
				</div>
				<div class="row mb-3">
					<button class="btn btn-danger" @click="doImpersonationStart">
						Impersonate
					</button>
				</div>
			</template>
		</div>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useSessionStore } from "@/store/session";
import { useUserStore } from "@/store/user";
import UserSelector from "@/components/UserSelector.vue";
import type { Session, User } from "@/type/api";

const session = useSessionStore();
const user = useUserStore();

const impersonatedUser = ref<User | null>(null);
const isImpersonating = ref<boolean>(false);
const selectedUser = ref<User | null>(null);

const doImpersonationEnd = async () => {
	const url = session.urls!.api.impersonation;
	const response = await fetch(url, {
		method: "DELETE",
	});
	if (!response.ok) {
		throw new Error(`Failed to end impersonation: ${response.statusText}`);
	}
	const new_session = await session.fetchSession();
	console.log("session is now", new_session);
};
const doImpersonationStart = async () => {
	if (!selectedUser.value) {
		console.log("Can't impersonate, null user");
		return;
	}
	console.log("doing impersonation of user", selectedUser.value);
	const body = {
		id: selectedUser.value.id,
	};
	const url = session.urls!.api.impersonation;
	const response = await fetch(url, {
		body: JSON.stringify(body),
		headers: {
			"Content-Type": "application/json",
		},
		method: "POST",
	});
	if (!response.ok) {
		throw new Error(`Upload failed: ${response.statusText}`);
	}
	const result = await response.json();
	console.log("impersonation", result);
	const new_session = await session.fetchSession();
	console.log("session is now", new_session);
};
/*
onMounted(() => {
	session.get().then((session: Session) => {
		if (session.impersonating) {
			isImpersonating.value = true;
			user.byURI(session.impersonating).then((user: User | null) => {
				impersonatedUser.value = user;
			});
		} else {
			isImpersonating.value = false;
			impersonatedUser.value = null;
		}
	});
});
*/
</script>
