<style scoped>
.tech-photo {
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
				<div class="table-responsive">
					<table class="table table-striped table-hover">
						<thead class="table-light">
							<tr>
								<th>User</th>
								<th>Role</th>
								<th>Status</th>
								<th>Tags</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="user in users" :key="user.id">
								<td>
									<div class="d-flex align-items-center">
										<img
											:src="user.avatar"
											:alt="user.name"
											class="tech-photo me-3"
										/>
										<div>
											<div class="fw-bold">{{ user.name }}</div>
										</div>
									</div>
								</td>
								<td>
									<span class="badge bg-success">{{ user.role }}</span>
								</td>
								<td>
									<span
										class="badge status-badge"
										:class="getStatusClass(user.status)"
									>
										{{ user.status }}
									</span>
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
									<button
										class="btn btn-sm btn-warning"
										title="Deactivate"
										@click="deactivateUser(user.id)"
									>
										<i class="bi bi-person-x"></i>
									</button>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

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
const users = ref<User[]>([
	{
		id: 1,
		name: "John Davis",
		avatar: "https://randomuser.me/api/portraits/men/32.jpg",
		role: "Tech I",
		status: "Active",
		tags: [],
	},
	{
		id: 2,
		name: "Sarah Johnson",
		avatar: "https://randomuser.me/api/portraits/women/65.jpg",
		role: "Tech III",
		status: "Active",
		tags: ["warrant service", "drone pilot"],
	},
	{
		id: 3,
		name: "Michael Chen",
		avatar: "https://randomuser.me/api/portraits/men/44.jpg",
		role: "Tech I",
		status: "Active",
		tags: ["drone pilot"],
	},
]);

const urlConfiguration = ref<URLConfiguration>({
	userAdd: "/configuration/user/add", // Update with your actual route
});

// Methods
const getStatusClass = (status: string): string => {
	return status === "Active" ? "bg-success" : "bg-secondary";
};

const getTagClass = (tag: string): string => {
	if (tag === "warrant service") return "bg-warrant";
	if (tag === "drone pilot") return "bg-drone";
	return "bg-secondary";
};

const deactivateUser = (userId: number): void => {
	const user = users.value.find((u) => u.id === userId);
	if (user) {
		user.status = "Inactive";
		// Add your deactivation logic here (e.g., API call)
		console.log(`Deactivating user: ${userId}`);
	}
};

// Lifecycle hooks
onMounted(() => {
	// Fetch users from API if needed
	// fetchUsers()
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
