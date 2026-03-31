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
			<div class="row mb-3">
				<div class="col-md-6">
					<label for="userSearch" class="form-label">Search Users</label>
					<UserSelector
						v-model="selectedUserId"
						label="Choose a user"
						placeholder="Select a user..."
						@change="onUserChange"
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

			<div class="table-responsive">
				<table class="table table-hover">
					<thead>
						<tr>
							<th>User ID</th>
							<th>Name</th>
							<th>Email</th>
							<th>Role</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>1001</td>
							<td>John Doe</td>
							<td>john.doe@example.com</td>
							<td><span class="badge bg-primary">Admin</span></td>
							<td>
								<button class="btn btn-sm btn-primary">
									Impersonate <i class="bi bi-box-arrow-in-right"></i>
								</button>
							</td>
						</tr>
						<tr>
							<td>1002</td>
							<td>Jane Smith</td>
							<td>jane.smith@example.com</td>
							<td><span class="badge bg-info">Support</span></td>
							<td>
								<button class="btn btn-sm btn-primary">
									Impersonate <i class="bi bi-box-arrow-in-right"></i>
								</button>
							</td>
						</tr>
						<tr>
							<td>1003</td>
							<td>Robert Johnson</td>
							<td>robert@example.com</td>
							<td><span class="badge bg-success">Premium User</span></td>
							<td>
								<button class="btn btn-sm btn-primary">
									Impersonate <i class="bi bi-box-arrow-in-right"></i>
								</button>
							</td>
						</tr>
						<tr>
							<td>1004</td>
							<td>Maria Garcia</td>
							<td>maria@example.com</td>
							<td><span class="badge bg-secondary">Standard User</span></td>
							<td>
								<button class="btn btn-sm btn-primary">
									Impersonate <i class="bi bi-box-arrow-in-right"></i>
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>

			<div
				class="alert alert-warning mt-3 mb-0"
				id="impersonationStatus"
				style="display: none"
			>
				<i class="bi bi-exclamation-triangle"></i> You are currently
				impersonating <strong>John Doe</strong>
				<button class="btn btn-sm btn-warning float-end">
					Exit Impersonation <i class="bi bi-box-arrow-left"></i>
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useSessionStore } from "@/store/session";
import type { User } from "@/types";

const session = useSessionStore();

const selectedUserId = ref<number | null>(null);

const onUserChange = (user: User | null) => {
	console.log("Selected user:", user);
};
</script>
