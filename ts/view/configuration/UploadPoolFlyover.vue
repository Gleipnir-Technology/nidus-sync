<template>
	<!-- Main Content -->
	<div class="container mt-4 upload-container">
		<h2 class="mb-4">Upload Pool Data</h2>

		<div class="card mb-4">
			<TableUploadRequirements />
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
import CSVUpload from "@/components/CSVUpload.vue";
import TableUploadRequirements from "@/components/TableUploadRequirements.vue";
import { router } from "@/router";

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
