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
				<form
					action="/configuration/upload/pool/flyover"
					method="POST"
					enctype="multipart/form-data"
				>
					<div class="upload-area">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="48"
							height="48"
							fill="currentColor"
							class="bi bi-cloud-arrow-up text-primary mb-3"
							viewBox="0 0 16 16"
						>
							<path
								fill-rule="evenodd"
								d="M7.646 5.146a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1-.708.708L8.5 6.707V10.5a.5.5 0 0 1-1 0V6.707L6.354 7.854a.5.5 0 1 1-.708-.708l2-2z"
							/>
							<path
								d="M4.406 3.342A5.53 5.53 0 0 1 8 2c2.69 0 4.923 2 5.166 4.579C14.758 6.804 16 8.137 16 9.773 16 11.569 14.502 13 12.687 13H3.781C1.708 13 0 11.366 0 9.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383zm.653.757c-.757.653-1.153 1.44-1.153 2.056v.448l-.445.049C2.064 6.805 1 7.952 1 9.318 1 10.785 2.23 12 3.781 12h8.906C13.98 12 15 10.988 15 9.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 4.825 10.328 3 8 3a4.53 4.53 0 0 0-2.941 1.1z"
							/>
						</svg>
						<h5>Select your CSV file</h5>
						<p class="text-muted">
							Drag and drop a file here or click to browse
						</p>
						<input
							type="file"
							class="form-control"
							id="csvFile"
							name="csvfile"
							accept=".csv"
						/>
					</div>

					<div class="d-grid gap-2 text-center">
						<button class="btn btn-primary" type="submit">
							Upload and Continue
						</button>
					</div>
				</form>
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
