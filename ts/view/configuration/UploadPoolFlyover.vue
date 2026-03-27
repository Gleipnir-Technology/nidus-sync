<template>
	<!-- Main Content -->
	<div class="container mt-4 upload-container">
		<h2 class="mb-4">Upload Pool Data</h2>

		<div class="card mb-4">
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
						<tr>
							<td class="required-field">City</td>
							<td>The city portion of the address</td>
							<td>Text</td>
							<td>Visalia</td>
						</tr>
						<tr>
							<td class="required-field">Comment</td>
							<td>The condition of the pool</td>
							<td>Text</td>
							<td>"blue", "dry", "false pool", "green", or "murky"</td>
						</tr>
						<tr>
							<td class="required-field">HouseNo</td>
							<td>The house number portion of the address</td>
							<td>Text</td>
							<td>123</td>
						</tr>
						<tr>
							<td class="required-field">State</td>
							<td>The state portion of the address</td>
							<td>Text</td>
							<td>California</td>
						</tr>
						<tr>
							<td class="required-field">Street</td>
							<td>The street portion of the address</td>
							<td>Text</td>
							<td>Main St</td>
						</tr>
						<tr>
							<td class="required-field">TargetLat</td>
							<td>The latitude of the target location</td>
							<td>Decimal Number</td>
							<td>36.56245379</td>
						</tr>
						<tr>
							<td class="required-field">TargetLon</td>
							<td>The longitude of the target location</td>
							<td>Decimal Number</td>
							<td>-119.3948222</td>
						</tr>
						<tr>
							<td class="required-field">ZIP</td>
							<td>The postal code (ZIP) portion of the address</td>
							<td>Text</td>
							<td>93681</td>
						</tr>
					</tbody>
				</table>
			</div>
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
import { router } from "@/router";

function onError(err) {
	console.error("CSV upload error", err);
}
function onFileSelected(file) {
	console.log("file selected", file);
}
function onUploadSuccess(data) {
	console.log("upload success", data);
	router.push("/_" + data.uri);
}
</script>
