<style scoped lang="scss">
.avatar {
	width: 48px;
	height: 48px;
	border-radius: 50%;
	object-fit: cover;
}

.status-badge {
	min-width: 70px;
	text-align: center;
}

.badge {
	margin-right: 0.25rem;
}

.bg-warrant {
	background-color: #6c757d;
}

.bg-drone {
	background-color: #0dcaf0;
}

/* Override Bootstrap hover if needed */
.table-hover tbody tr:hover {
	cursor: pointer;
}
</style>
<template>
	<div class="container-fluid p-4">
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h1 class="mb-0">User Management</h1>
			<RouterLink to="/configuration/user/add">
				<button class="btn btn-primary" id="addUserBtn">
					<i class="bi bi-plus-circle me-2"></i>Add New User
				</button>
			</RouterLink>
		</div>

		<div class="card">
			<div class="card-body">
				<div v-if="users" class="table-responsive">
					<table class="table table-striped table-hover">
						<thead class="table-light">
							<tr>
								<th>User</th>
								<th>Role</th>
								<th>Tags</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="user in users" :key="user.id">
								<td>
									<div class="d-flex align-items-center">
										<i v-if="!user.avatar" class="bi bi-avatar avatar"></i>
										<img
											:src="user.avatar"
											:alt="user.display_name"
											class="avatar me-3"
											v-else
										/>
										<div>
											<div class="fw-bold">{{ user.display_name }}</div>
										</div>
									</div>
								</td>
								<td>
									<span class="badge bg-success">{{ user.role }}</span>
								</td>
								<td>
									<span
										v-for="tag in user.tags"
										:key="tag"
										class="badge"
										:class="getTagClass(tag)"
									>
										{{ tag }}
									</span>
								</td>
								<td>
									<RouterLink :to="`/_/configuration/user/${user.id}`">
										<button class="btn btn-sm btn-primary" title="Edit">
											<i class="bi bi-person-x"></i>
										</button>
									</RouterLink>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div v-else>
					<p>loading...</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useUserStore } from "@/store/user";

interface User {
	id: number;
	name: string;
	avatar: string;
	role: string;
	status: "Active" | "Inactive";
	tags: string[];
}

interface URLConfiguration {
	userAdd: string;
}

// Props (if needed from parent component)
// const props = defineProps<{
//   urlConfiguration?: URLConfiguration
// }>()

// Reactive state
const userStore = useUserStore();
const users = computed(() => {
	return userStore.all;
});
const urlConfiguration = ref<URLConfiguration>({
	userAdd: "/configuration/user/add", // Update with your actual route
});

const getTagClass = (tag: string): string => {
	if (tag === "warrant service") return "bg-warrant";
	if (tag === "drone pilot") return "bg-drone";
	return "bg-secondary";
};

const deactivateUser = (userId: number): void => {
	if (users.value == null) {
		return;
	}
	const user = users.value.find((u) => u.id === userId);
	if (!user) {
		return;
	}
	user.is_active = false;
	// Add your deactivation logic here (e.g., API call)
	console.log(`Deactivating user: ${userId}`);
};

// Lifecycle hooks
onMounted(() => {
	// Fetch users from API if needed
	userStore.fetchAll();
});

// Optional: API call example
// const fetchUsers = async (): Promise<void> => {
//   try {
//     const response = await fetch('/api/users')
//     const data = await response.json()
//     users.value = data
//   } catch (error) {
//     console.error('Error fetching users:', error)
//   }
// }
</script>
