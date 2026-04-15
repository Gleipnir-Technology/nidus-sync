<style scoped>
.upload-card {
	transition:
		transform 0.2s,
		box-shadow 0.2s;
	cursor: pointer;
	height: 100%;
}
.upload-card:hover {
	transform: translateY(-5px);
	box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}
.card-icon {
	font-size: 3rem;
	margin-bottom: 1rem;
}
</style>
<template>
	<div class="container py-5">
		<div class="row mb-4">
			<div class="col-12 text-center">
				<h1 class="display-5 fw-bold mb-2">Green Pool CSV Data</h1>
				<p class="lead text-muted">
					Select the type of data you want to upload
				</p>
			</div>
		</div>

		<div class="row g-4">
			<!-- Option 1: Green Pool Flyover Data -->
			<div class="col-md-6">
				<div class="card upload-card border-0 shadow-sm">
					<div class="card-body text-center p-5">
						<div class="card-icon text-primary">
							<i class="bi bi-cloud-upload"></i>
						</div>
						<h3 class="card-title h4 mb-3">Green Pool Flyover Data</h3>
						<p class="card-text text-muted mb-4">
							Upload aerial survey data from ABC Data Analytics. This includes
							GPS coordinates, timestamp information, and pool identification
							data.
						</p>
						<div class="mb-3">
							<span
								class="badge bg-primary-subtle text-primary border border-primary me-2"
							>
								<i class="bi bi-building me-1"></i> ABC Data Analytics
							</span>
							<span class="badge bg-info-subtle text-info border border-info">
								<i class="bi bi-filetype-csv me-1"></i> CSV Format
							</span>
						</div>
						<RouterLink to="/_/configuration/upload/pool/flyover">
							<button class="btn btn-primary btn-lg w-100">
								<i class="bi bi-upload me-2"></i>Let's do this
							</button>
						</RouterLink>
					</div>
				</div>
			</div>

			<!-- Option 2: Custom Green Pool Operations Data -->
			<div class="col-md-6">
				<div class="card upload-card border-0 shadow-sm">
					<div class="card-body text-center p-5">
						<div class="card-icon text-success">
							<i class="bi bi-file-earmark-spreadsheet"></i>
						</div>
						<h3 class="card-title h4 mb-3">Custom Operations Data</h3>
						<p class="card-text text-muted mb-4">
							Upload custom green pool operations data. This includes treatment
							records, inspection logs, and maintenance activities in your own
							CSV format.
						</p>
						<div class="mb-3">
							<span
								class="badge bg-success-subtle text-success border border-success me-2"
							>
								<i class="bi bi-gear me-1"></i> Custom Format
							</span>
							<span class="badge bg-info-subtle text-info border border-info">
								<i class="bi bi-filetype-csv me-1"></i> CSV Format
							</span>
						</div>
						<RouterLink to="/_/configuration/upload/pool/custom">
							<button class="btn btn-success btn-lg w-100">
								<i class="bi bi-upload me-2"></i>Pick me
							</button>
						</RouterLink>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Upload Modal -->
	<div
		class="modal fade"
		id="uploadModal"
		tabindex="-1"
		aria-labelledby="uploadModalLabel"
		aria-hidden="true"
	>
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="uploadModalLabel">Upload CSV File</h5>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
					></button>
				</div>
				<div class="modal-body">
					<div class="mb-3">
						<label for="fileInput" class="form-label">Select CSV file:</label>
						<input
							type="file"
							class="form-control"
							id="fileInput"
							accept=".csv"
						/>
					</div>
					<div class="alert alert-secondary mb-0" role="alert">
						<small
							><strong>Upload Type:</strong>
							<span id="uploadType">Not selected</span></small
						>
					</div>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						data-bs-dismiss="modal"
					>
						Cancel
					</button>
					<button type="button" class="btn btn-primary">
						<i class="bi bi-upload me-2"></i>Upload
					</button>
				</div>
			</div>
		</div>
	</div>
</template>
