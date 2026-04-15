<template>
	<div class="card-header bg-light">
		<h5 class="mb-0">CSV Upload Requirements</h5>
	</div>
	<div class="card-body">
		<p>
			Your CSV file must contain the following columns in any order. Please
			ensure your data matches the required format.
		</p>

		<table class="table table-bordered schema-table">
			<thead class="table-light">
				<tr>
					<th>Field</th>
					<th>Description</th>
					<th>Format</th>
					<th>Example</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="(req, index) in requirements">
					<td :class="req.is_required ? 'required-field' : ''">
						{{ req.field }}
					</td>
					<td>{{ req.description }}</td>
					<td>{{ req.format }}</td>
					<td>{{ req.example }}</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<script setup lang="ts">
export interface UploadRequirement {
	description: string;
	example: string;
	field: string;
	format: string;
	is_required: boolean;
}
interface Props {
	requirements: UploadRequirement[];
}
const props = defineProps<Props>();
</script>
