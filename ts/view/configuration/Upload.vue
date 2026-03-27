<style scoped>
.upload-card {
	transition: transform 0.2s;
	margin-bottom: 30px;
}
.upload-card:hover {
	transform: translateY(-5px);
	box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}
.card-icon {
	font-size: 2.5rem;
	margin-bottom: 15px;
	color: #198754;
}
.badge {
	--bs-bg-opacity: 1;
}
.badge.committed {
	background-color: $success;
}
.badge.committing {
	background-color: $success;
}
.badge.discarded {
	background-color: gray;
}
.badge.error {
	background-color: $danger;
}
.badge.parsed {
	background-color: $secondary;
}
.badge.uploaded {
	background-color: $info;
}
</style>
<template>
	<div class="container mb-5">
		<div class="row">
			<!-- Green Pool Management -->
			<div class="col-md-4">
				<div class="card h-100 upload-card border-success">
					<div class="card-body text-center">
						<i class="bi bi-water card-icon"></i>
						<h5 class="card-title">Green Pool Management</h5>
						<p class="card-text">
							Upload spreadsheets with addresses and contact information of
							unmaintained pools that may breed mosquitoes.
						</p>
						<RouterLink to="/configuration/upload/pool">
							<button class="btn btn-primary">
								<i class="bi bi-upload me-2"></i>Upload Green Pool Data</button
							>
						</RouterLink>
					</div>
					<div class="card-footer bg-white text-muted">
						<small><i class="bi bi-clock"></i> Last import: 02/15/2023</small>
					</div>
				</div>
			</div>

			<!-- Employee Information -->
			<div class="col-md-4">
				<div class="card h-100 upload-card border-primary">
					<div class="card-body text-center">
						<i class="bi bi-people-fill card-icon" style="color: #0d6efd"></i>
						<h5 class="card-title">Employee Information</h5>
						<p class="card-text">
							Import employee data including names, contact information, and
							responsibilities for system user creation.
						</p>
						<!--
						<p class="text-muted small">Supported formats: .xlsx, .csv</p>
						<div class="mb-3">
							<label for="employeeFile" class="form-label"
								>Select file to import</label
							>
							<input class="form-control" type="file" id="employeeFile" />
						</div>
						<button type="button" class="btn btn-primary disabled" disabled>
							<i class="bi bi-upload me-2"></i>Upload Employee Data
						</button>
						-->
					</div>
					<div class="card-footer bg-white text-muted">
						<small><i class="bi bi-clock"></i> Last import: 03/01/2023</small>
					</div>
				</div>
			</div>

			<!-- Field Notebooks -->
			<div class="col-md-4">
				<div class="card h-100 upload-card border-warning">
					<div class="card-body text-center">
						<i class="bi bi-journal-text card-icon" style="color: #fd7e14"></i>
						<h5 class="card-title">Field Notebooks</h5>
						<p class="card-text">
							Upload scanned technician field notebooks to digitize information
							about breeding sources they've identified.
						</p>
						<!--
						<p class="text-muted small">Supported formats: .pdf, .jpg, .png</p>
						<div class="mb-3">
							<label for="notebookFile" class="form-label"
								>Select file to import</label
							>
							<input class="form-control" type="file" id="notebookFile" />
						</div>
						<button type="button" class="btn btn-warning text-white">
							<i class="bi bi-upload me-2"></i>Upload Notebook Data
						</button>
						-->
					</div>
					<div class="card-footer bg-white text-muted">
						<small><i class="bi bi-clock"></i> Last import: 03/15/2023</small>
					</div>
				</div>
			</div>
		</div>

		<!-- Import History Section -->
		<div class="row mt-5">
			<div class="col-12">
				<h3><i class="bi bi-clock-history"></i> Recent Import History</h3>
				<table class="table table-striped table-hover">
					<thead class="table-dark">
						<tr>
							<th scope="col">Date/Time</th>
							<th scope="col">Import Type</th>
							<th scope="col">Filename</th>
							<th scope="col">Status</th>
							<th scope="col">Records</th>
							<th scope="col">Actions</th>
						</tr>
					</thead>
					<tbody>
							<tr v-for="upload in uploads">
								<td><TimeRelative :time="upload.created"/></td>
								<td>{{upload.type}}</td>
								<td>{{upload.filename}}</td>
								<td>
									<span class="badge" :class="upload.status"
										>{{upload.status}}</span
									>
								</td>
								<td>{{upload.record_count}} entries</td>
								<td>
									<a
										class="btn btn-sm btn-outline-primary"
										:href="`/configuration/upload/${upload.id}`"
										>View</a
									>
								</td>
							</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
</template>
