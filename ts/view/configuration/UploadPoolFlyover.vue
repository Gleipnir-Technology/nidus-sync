<template>
	<!-- Main Content -->
	<div class="container mt-4 upload-container">
		<h2 class="mb-4">Upload Pool Flyover Data</h2>

		<div class="card mb-4">
			<TableUploadRequirements :requirements="requirements" />
		</div>

		<div class="card">
			<div class="card-header bg-light">
				<h5 class="mb-0">Upload Data</h5>
			</div>
			<div class="card-body">
				<CSVUpload
					upload-url="/api/upload/pool/flyover"
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
import { router } from "@/route/config";

const requirements = ref<UploadRequirement[]>([
	{
		description: "The city portion of the address",
		example: "Visalia",
		field: "City",
		format: "Text",
		is_required: true,
	},
	{
		description: "The condition of the pool",
		example: '"blue", "dry", "false pool", "green", or "murky"',
		field: "Comment",
		format: "Text",
		is_required: true,
	},
	{
		description: "The house number portion of the address",
		example: "123",
		field: "HouseNo",
		format: "Text",
		is_required: true,
	},
	{
		description: "The state portion of the address",
		example: "California",
		field: "State",
		format: "Text",
		is_required: true,
	},
	{
		description: "The street portion of the address",
		example: "Main St",
		field: "Street",
		format: "Text",
		is_required: true,
	},
	{
		description: "The latitude of the target location",
		example: "36.56245379",
		field: "TargetLat",
		format: "Decimal Number",
		is_required: true,
	},
	{
		description: "The longitude of the target location",
		example: "-119.3948222",
		field: "TargetLon",
		format: "Decimal Number",
		is_required: true,
	},
	{
		description: "The postal code (ZIP) portion of the address",
		example: "93681",
		field: "ZIP",
		format: "Text",
		is_required: true,
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
