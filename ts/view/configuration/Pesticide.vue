<template>
	<div class="container-fluid p-4">
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h1 class="mb-0">Pesticide Products Configuration</h1>
			<RouterLink to="/configuration/pesticide/add"
				><button class="btn btn-primary" id="addProductBtn">
					<i class="bi bi-plus-circle me-2"></i>Add New Product
				</button>
			</RouterLink>
		</div>

		<div class="card">
			<div class="card-body">
				<div class="table-responsive">
					<table class="table table-striped table-hover align-middle">
						<thead class="table-light">
							<tr>
								<th>Product</th>
								<th>Formulation</th>
								<th>Targets</th>
								<th>Residual (days)</th>
								<th>Low Rate</th>
								<th>Max Rate</th>
								<th>Pools</th>
								<th>Info</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="product in products" :key="product.id">
								<td>
									<strong>{{ product.name }}</strong>
								</td>
								<td>{{ product.formulation }}</td>
								<td>
									<span
										v-for="target in targets"
										:key="target.code"
										class="target-icon"
										:class="
											product.activeTargets.includes(target.code)
												? 'target-active'
												: 'target-inactive'
										"
										:title="target.title"
									>
										{{ target.code }}
									</span>
								</td>
								<td>{{ product.residualDays }}</td>
								<td>{{ product.lowRate }}</td>
								<td>{{ product.maxRate }}</td>
								<td>
									<span
										class="badge"
										:class="getPoolBadgeClass(product.poolStatus)"
									>
										{{ product.poolStatus }}
									</span>
								</td>
								<td>
									<a
										:href="`product-details.html?id=${product.id}`"
										class="btn btn-sm btn-info"
										title="Product Information"
									>
										<i class="bi bi-info-circle"></i>
									</a>
								</td>
								<td>
									<button
										class="btn btn-sm btn-primary"
										title="Edit"
										@click="editProduct(product.id)"
									>
										<i class="bi bi-pencil"></i>
									</button>
									<button
										class="btn btn-sm btn-danger"
										title="Delete"
										@click="deleteProduct(product.id)"
									>
										<i class="bi bi-trash"></i>
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
import { ref, computed } from "vue";

interface Target {
	code: string;
	title: string;
}

interface Product {
	id: string;
	name: string;
	formulation: string;
	activeTargets: string[];
	residualDays: number;
	lowRate: string;
	maxRate: string;
	poolStatus: "Recommended" | "Allowed" | "Prohibited" | "Discouraged";
}

// Configuration URL (you can replace this with a prop or composable)
const configurationUrl = ref("/configuration/pesticide/add");

// Define target stages
const targets: Target[] = [
	{ code: "I1", title: "Instar Stage 1" },
	{ code: "I2", title: "Instar Stage 2" },
	{ code: "I3", title: "Instar Stage 3" },
	{ code: "I4", title: "Instar Stage 4" },
	{ code: "P", title: "Pupae" },
];

// Sample pesticide data
const products = ref<Product[]>([
	{
		id: "bva-oil",
		name: "BVA Oil",
		formulation: "Liquid",
		activeTargets: ["P"],
		residualDays: 1,
		lowRate: "0.5 gal/acre",
		maxRate: "5 gal/acre",
		poolStatus: "Recommended",
	},
	{
		id: "vectomax-fg",
		name: "VectoMax FG",
		formulation: "Granule",
		activeTargets: ["I1", "I2", "I3", "I4"],
		residualDays: 30,
		lowRate: "5 lbs/acre",
		maxRate: "20 lbs/acre",
		poolStatus: "Recommended",
	},
	{
		id: "censor",
		name: "Censor",
		formulation: "Liquid",
		activeTargets: ["I1", "I2", "I3", "I4"],
		residualDays: 21,
		lowRate: "0.75 gal/acre",
		maxRate: "2.5 gal/acre",
		poolStatus: "Allowed",
	},
	{
		id: "aquabac-xt",
		name: "AquaBac XT",
		formulation: "Liquid",
		activeTargets: ["I1", "I2", "I3"],
		residualDays: 14,
		lowRate: "0.25 gal/acre",
		maxRate: "2 gal/acre",
		poolStatus: "Prohibited",
	},
	{
		id: "natular-g30",
		name: "Natular G30",
		formulation: "Granule",
		activeTargets: ["I1", "I2", "I3", "I4"],
		residualDays: 30,
		lowRate: "5 lbs/acre",
		maxRate: "12 lbs/acre",
		poolStatus: "Discouraged",
	},
]);

// Helper function to get badge class based on pool status
const getPoolBadgeClass = (status: Product["poolStatus"]): string => {
	const statusMap: Record<Product["poolStatus"], string> = {
		Recommended: "bg-success",
		Allowed: "bg-warning text-dark",
		Prohibited: "bg-danger",
		Discouraged: "bg-secondary",
	};
	return statusMap[status] || "bg-secondary";
};

// Action handlers
const editProduct = (productId: string): void => {
	console.log("Edit product:", productId);
	// Implement edit logic or navigation
};

const deleteProduct = (productId: string): void => {
	console.log("Delete product:", productId);
	// Implement delete logic with confirmation
	if (confirm("Are you sure you want to delete this product?")) {
		const index = products.value.findIndex((p) => p.id === productId);
		if (index !== -1) {
			products.value.splice(index, 1);
		}
	}
};
</script>

<style scoped>
.target-icon {
	display: inline-block;
	width: 24px;
	height: 24px;
	text-align: center;
	line-height: 24px;
	border-radius: 50%;
	font-size: 12px;
	font-weight: bold;
	margin-right: 2px;
	color: white;
}

.target-active {
	background-color: #0d6efd;
}

.target-inactive {
	background-color: #dee2e6;
	color: #6c757d;
}

.table-responsive {
	overflow-x: auto;
}

th {
	white-space: nowrap;
}
</style>
