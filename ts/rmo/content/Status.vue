<template>
	<div class="container my-4">
		<!-- Search Box -->
		<div class="card search-box mb-4">
			<div class="card-body">
				<form class="row g-3 align-items-center" action="#" id="lookup-form">
					<div class="col-md-9">
						<address-or-report-input
							name="address-or-report"
							placeholder="Enter a report ID, address, neighborhood, or zip code"
						></address-or-report-input>
					</div>
					<div class="col-md-3">
						<span
							data-bs-toggle="tooltip"
							id="lookup-tooltip"
							title="You can look up a report once you type in the full report ID. Start typing and I'll suggest complete IDs"
						>
							<button
								type="submit"
								class="btn btn-primary btn-lg w-100 disabled"
								disabled
								id="lookup"
							>
								Lookup Report by ID
							</button>
						</span>
					</div>
					<div class="col-12">
						<div class="form-check custom-circle-checkbox">
							<input
								class="form-check-input"
								type="checkbox"
								id="checkboxNuisance"
								data-color="danger"
								checked
							/>
							<label class="form-check-label" for="checkboxNuisance"
								>Mosquito Nuisance</label
							>
						</div>

						<div class="form-check custom-circle-checkbox">
							<input
								class="form-check-input"
								type="checkbox"
								id="checkboxWater"
								data-color="success"
								checked
							/>
							<label class="form-check-label" for="checkboxWater"
								>Standing Water</label
							>
						</div>
					</div>
				</form>
			</div>
		</div>

		<!-- Map Section -->
		<div class="card mb-4">
			<div class="card-header bg-info text-white">
				<h5 class="mb-0"><i class="bi bi-pin-map-fill me-2"></i>Reports Map</h5>
			</div>
			<div class="card-body p-0">
				<div class="map-container">
					<map-multipoint
						id="map"
						xmax="-118.0"
						ymax="37.0"
						xmin="-119.0"
						ymin="36.0"
						tegola="{{ .URL.Tegola }}"
						zoom="9"
					></map-multipoint>
				</div>
			</div>
		</div>

		<!-- Results Section -->
		<div class="card">
			<div
				class="card-header bg-primary text-white d-flex justify-content-between align-items-center"
			>
				<h5 class="mb-0">
					<i class="bi bi-geo-fill me-2"></i>Reports Near You
				</h5>
				<span class="badge bg-light text-dark" id="report-count"
					>- Reports Found</span
				>
			</div>
			<div class="card-body p-0">
				<div class="table-responsive">
					<table-report />
				</div>
			</div>
			<!--
		<div class="card-footer">
			<nav aria-label="Page navigation">
				<ul class="pagination justify-content-center mb-0">
					<li class="page-item disabled">
						<a class="page-link" href="#" tabindex="-1" aria-disabled="true">Previous</a>
					</li>
					<li class="page-item active"><a class="page-link" href="#">1</a></li>
					<li class="page-item"><a class="page-link" href="#">2</a></li>
					<li class="page-item"><a class="page-link" href="#">3</a></li>
					<li class="page-item">
						<a class="page-link" href="#">Next</a>
					</li>
				</ul>
			</nav>
		</div>
		-->
		</div>
	</div>
</template>
