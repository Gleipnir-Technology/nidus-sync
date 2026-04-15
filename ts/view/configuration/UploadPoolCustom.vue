<template>
	<!-- Main Content -->
	<div class="container mt-4 upload-container">
		<h2 class="mb-4">Upload Pool Custom Data</h2>

		<div class="card mb-4">
			<TableUploadRequirements :requirements="requirements" />
		</div>

		<div class="card">
			<div class="card-header bg-light">
				<h5 class="mb-0">Upload Data</h5>
			</div>
			<div class="card-body">
				<CSVUpload
					upload-url="/api/upload/pool/custom"
					@doError="onError"
					@doFileSelected="onFileSelected"
					@doSuccess="onUploadSuccess"
				/>
			</div>
		</div>

		<div class="text-muted text-center mt-4">
			<small
				>Need assistance? Contact
				<a href="mailto:support@example.com">support@example.com</a></small
			>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import CSVUpload from "@/components/CSVUpload.vue";
import TableUploadRequirements, {
	UploadRequirement,
} from "@/components/TableUploadRequirements.vue";
import { router } from "@/router";

const requirements = ref<UploadRequirement[]>([
	{
		field: "Street Address",
		description: "Street number and name of the address of the pool",
		format: "Text",
		example: "123 Main St.",
		is_required: true,
	},
	{
		field: "City",
		description: "The city portion of the pool's address",
		format: "Text",
		example: "Visalia",
		is_required: true,
	},
	{
		field: "Notes",
		description:
			" Any notes from the district to include with the pool record ",
		format: "Text",
		example: '"Collects rain water when empty"',
		is_required: false,
	},
	{
		field: "Postal Code",
		description: "Postal (Zip) Code of the pool's address",
		format: "numbers and optional hypen",
		example: "81234 or 91234-5678",
		is_required: true,
	},
	{
		field: "Pool Condition",
		description: "The condition of the pool when it was last inspected",
		format: "Text",
		example: '"blue", "dry", "false pool", "green", or "murky"',
		is_required: false,
	},
	{
		field: "Property Owner Name",
		description: "Name of the person or entity that owns the property",
		format: "Text",
		example: "No",
		is_required: false,
	},
	{
		field: "Property Owner Phone",
		description: "Phone number of the person or entity that owns the property",
		format: "e164",
		example:
			'"+14155552671" or "1-(901)-555-1234" or "9015551234" or "1901-555-12-34"',
		is_required: false,
	},
	{
		field: "Resident Owned",
		description:
			" Whether or not the current resident of the property is also the owner ",
		format: "Yes, No, or empty",
		example: '"Yes" or "No" or ""',
		is_required: false,
	},
	{
		field: "Resident Phone",
		description: "Phone number of the resident",
		format: "e164",
		example:
			'"+14155552671" or "1-(901)-555-1234" or "9015551234" or "1901-555-12-34"',
		is_required: false,
	},
	{
		field: "Tags",
		description:
			" Any additional columns in the file will be treated as tags and attached to the record ",
		format: "Text",
		example: '"Hostile" or "Unresponsive" or "Dog"',
		is_required: false,
	},
]);

function onError(err: Error) {
	console.error("CSV upload error", err);
}
function onFileSelected(file: File) {
	console.log("file selected", file);
}
function onUploadSuccess(data: any) {
	console.log("upload success", data);
	router.push("/_" + data.uri);
}
</script>
