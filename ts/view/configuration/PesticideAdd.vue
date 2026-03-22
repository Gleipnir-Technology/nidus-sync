# VueJS Single-File Component (TypeScript) ```vue
<template>
	<div class="container py-4">
		<!-- Breadcrumb -->
		<nav aria-label="breadcrumb" class="mb-4">
			<ol class="breadcrumb">
				<li class="breadcrumb-item"><a href="#">Settings</a></li>
				<li class="breadcrumb-item">
					<a href="pesticide-config.html">Pesticide</a>
				</li>
				<li class="breadcrumb-item active" aria-current="page">
					{{ pesticide.name }}
				</li>
			</ol>
		</nav>

		<!-- Main Content -->
		<div class="card shadow-sm mb-4">
			<div class="card-body">
				<!-- Product Header -->
				<div class="d-flex justify-content-between align-items-start mb-4">
					<div>
						<h1 class="mb-2">{{ pesticide.name }}</h1>
						<p class="text-muted mb-0">
							{{ pesticide.description }}
						</p>
					</div>
					<span class="tag tag-enabled" v-if="pesticide.enabled">
						<i class="bi bi-check-circle-fill"></i> Enabled
					</span>
					<span class="tag tag-optional" v-else>
						<i class="bi bi-x-circle-fill"></i> Disabled
					</span>
				</div>

				<!-- General Information -->
				<div class="mb-4">
					<h2 class="section-heading">General Information</h2>
					<div class="row g-3">
						<div class="col-md-6 col-lg-4">
							<div class="info-label">Formulation</div>
							<div>{{ pesticide.formulation }}</div>
						</div>
						<div class="col-md-6 col-lg-4">
							<div class="info-label">EPA Registration Number</div>
							<div>{{ pesticide.epaNumber }}</div>
						</div>
						<div class="col-md-6 col-lg-4">
							<div class="info-label">Active Ingredients</div>
							<div>
								<div
									v-for="ingredient in pesticide.activeIngredients"
									:key="ingredient.name"
								>
									{{ ingredient.name }} ({{ ingredient.percentage }}%)
								</div>
							</div>
						</div>
						<div class="col-md-6 col-lg-4">
							<div class="info-label">Biological Targeting</div>
							<div class="mt-1">
								<span
									v-for="stage in instarStages"
									:key="stage.code"
									class="target-icon"
									:class="stage.active ? 'target-active' : 'target-inactive'"
									:title="stage.label"
								>
									{{ stage.code }}
								</span>
							</div>
						</div>
						<div class="col-md-6 col-lg-4">
							<div class="info-label">Application Rates</div>
							<div>
								Low: {{ pesticide.applicationRates.low }}<br />
								High: {{ pesticide.applicationRates.high }}
							</div>
						</div>
						<div class="col-md-6 col-lg-4">
							<div class="info-label">Residual</div>
							<div>{{ pesticide.residual }}</div>
						</div>
					</div>
				</div>

				<!-- Usage Notes -->
				<div class="alert alert-info mb-4" v-if="pesticide.usageNotes">
					<div class="d-flex">
						<div class="me-3">
							<i class="bi bi-info-circle-fill fs-4"></i>
						</div>
						<div>
							<h5 class="alert-heading">Key Usage Notes</h5>
							<p class="mb-0">
								{{ pesticide.usageNotes }}
							</p>
						</div>
					</div>
				</div>

				<!-- PPE Requirements -->
				<div class="mb-4">
					<h2 class="section-heading">PPE Requirements</h2>
					<div>
						<span
							v-for="ppe in pesticide.ppeRequirements"
							:key="ppe.name"
							class="tag"
							:class="ppe.optional ? 'tag-optional' : 'tag-ppe'"
						>
							<i :class="`bi ${ppe.icon}`"></i>
							{{ ppe.name }}
							<span v-if="ppe.optional"> (Optional)</span>
						</span>
					</div>
				</div>

				<!-- Equipment Supported -->
				<div class="mb-4">
					<h2 class="section-heading">Equipment Supported</h2>
					<div>
						<span
							v-for="equipment in pesticide.equipmentSupported"
							:key="equipment.name"
							class="tag tag-equipment"
						>
							<i :class="`bi ${equipment.icon}`"></i>
							{{ equipment.name }}
						</span>
					</div>
				</div>

				<!-- Suitability -->
				<div class="mb-4">
					<h2 class="section-heading">Suitability</h2>
					<div class="row g-3">
						<div
							v-for="suit in pesticide.suitability"
							:key="suit.type"
							class="col-md-6 col-lg-3"
						>
							<div class="info-label">{{ suit.type }}</div>
							<div>
								<span
									class="badge"
									:class="getSuitabilityBadgeClass(suit.level)"
								>
									{{ suit.label }}
								</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Actions -->
				<div class="d-flex justify-content-between mt-5 pt-3 border-top">
					<button
						class="btn btn-outline-danger"
						@click="handleRemoveFromInventory"
					>
						<i class="bi bi-trash me-2"></i> Remove from Inventory
					</button>
					<button class="btn btn-success" @click="handleAddToInventory">
						<i class="bi bi-plus-circle me-2"></i> Add to Allowed Inventory
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";

interface ActiveIngredient {
	name: string;
	percentage: number;
}

interface ApplicationRates {
	low: string;
	high: string;
}

interface BiologicalTarget {
	code: string;
	label: string;
	active: boolean;
}

interface PPERequirement {
	name: string;
	icon: string;
	optional: boolean;
}

interface Equipment {
	name: string;
	icon: string;
}

interface Suitability {
	type: string;
	level: "recommended" | "ok" | "none" | "warning";
	label: string;
}

interface Pesticide {
	name: string;
	description: string;
	enabled: boolean;
	formulation: string;
	epaNumber: string;
	activeIngredients: ActiveIngredient[];
	biologicalTargets: BiologicalTarget[];
	applicationRates: ApplicationRates;
	residual: string;
	usageNotes: string;
	ppeRequirements: PPERequirement[];
	equipmentSupported: Equipment[];
	suitability: Suitability[];
}

// Sample data - replace with API call or props
const pesticide = ref<Pesticide>({
	name: "VectoMax FG",
	description:
		"Biological larvicide granules combining Bacillus thuringiensis subspecies israelensis and Bacillus sphaericus for extended residual control of mosquito larvae.",
	enabled: true,
	formulation: "Granule",
	epaNumber: "73049-429",
	activeIngredients: [
		{ name: "Bacillus thuringiensis subspecies israelensis", percentage: 2.7 },
		{ name: "Bacillus sphaericus", percentage: 4.5 },
	],
	biologicalTargets: [
		{ code: "I1", label: "Instar Stage 1", active: true },
		{ code: "I2", label: "Instar Stage 2", active: true },
		{ code: "I3", label: "Instar Stage 3", active: true },
		{ code: "I4", label: "Instar Stage 4", active: true },
		{ code: "P", label: "Pupae", active: false },
	],
	applicationRates: {
		low: "5 lbs/acre",
		high: "20 lbs/acre",
	},
	residual: "Up to 30 days (environmental conditions dependent)",
	usageNotes:
		"Apply evenly across water surface. Use higher rate when L4 present or when organic load is high. Avoid application in ponds with fish unless approved by a supervisor.",
	ppeRequirements: [
		{ name: "Gloves", icon: "bi-hand-thumbs-up", optional: false },
		{ name: "Eye Protection", icon: "bi-eyeglasses", optional: false },
		{ name: "Respirator", icon: "bi-mask", optional: true },
	],
	equipmentSupported: [
		{ name: "Backpack Spreader", icon: "bi-backpack" },
		{ name: "Hand Spreader", icon: "bi-hand-index-thumb" },
		{ name: "Truck Granule Unit", icon: "bi-truck" },
	],
	suitability: [
		{ type: "Pools", level: "recommended", label: "Recommended" },
		{ type: "Vegetation", level: "ok", label: "OK" },
		{ type: "High Organics", level: "ok", label: "OK" },
		{ type: "Organic Crop Restriction", level: "none", label: "None" },
	],
});

const instarStages = computed(() => pesticide.value.biologicalTargets);

const getSuitabilityBadgeClass = (level: string): string => {
	const classes: Record<string, string> = {
		recommended: "bg-success",
		ok: "bg-info text-dark",
		none: "bg-secondary",
		warning: "bg-warning text-dark",
	};
	return classes[level] || "bg-secondary";
};

const handleRemoveFromInventory = (): void => {
	// Implement remove logic
	console.log("Remove from inventory");
	// Example: emit event or call API
	// emit('remove', pesticide.value);
};

const handleAddToInventory = (): void => {
	// Implement add logic
	console.log("Add to inventory");
	// Example: emit event or call API
	// emit('add', pesticide.value);
};
</script>

<style scoped>
.section-heading {
	font-size: 1.2rem;
	font-weight: 600;
	margin-bottom: 1rem;
	padding-bottom: 0.5rem;
	border-bottom: 1px solid #dee2e6;
}

.info-label {
	font-weight: 600;
	color: #495057;
}

.target-icon {
	display: inline-block;
	width: 30px;
	height: 30px;
	text-align: center;
	line-height: 30px;
	border-radius: 50%;
	font-size: 14px;
	font-weight: bold;
	margin-right: 4px;
	color: white;
}

.target-active {
	background-color: #0d6efd;
}

.target-inactive {
	background-color: #dee2e6;
	color: #6c757d;
}

.tag {
	display: inline-flex;
	align-items: center;
	padding: 0.4rem 0.8rem;
	margin: 0.2rem;
	border-radius: 30px;
	background-color: #f8f9fa;
	border: 1px solid #dee2e6;
	font-size: 0.9rem;
}

.tag i {
	margin-right: 0.5rem;
}

.tag-enabled {
	background-color: #d1e7dd;
	color: #0f5132;
	border-color: #a3cfbb;
}

.tag-ppe {
	background-color: #e2e3e5;
	color: #41464b;
	border-color: #d3d6d8;
}

.tag-equipment {
	background-color: #cff4fc;
	color: #055160;
	border-color: #9eeaf9;
}

.tag-suitability {
	background-color: #fff3cd;
	color: #664d03;
	border-color: #ffecb5;
}

.tag-optional {
	background-color: #f8f9fa;
	color: #6c757d;
	border-color: #dee2e6;
}
</style>
