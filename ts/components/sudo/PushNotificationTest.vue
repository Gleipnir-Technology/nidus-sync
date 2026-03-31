<template>
	<!-- Push Notification Testing -->
	<div class="card mb-5">
		<div class="card-header bg-danger text-white">
			<i class="bi bi-bell"></i> Push Notification Testing
		</div>
		<div class="card-body">
			<form>
				<div class="row">
					<div class="col-md-6">
						<div class="mb-3">
							<label for="pushUser" class="form-label"
								>User ID or Device Token</label
							>
							<input
								type="text"
								class="form-control"
								id="pushUser"
								placeholder="User ID or Device Token"
							/>
						</div>
					</div>
					<div class="col-md-6">
						<div class="mb-3">
							<label for="pushTitle" class="form-label"
								>Notification Title</label
							>
							<input
								type="text"
								class="form-control"
								id="pushTitle"
								placeholder="Notification Title"
							/>
						</div>
					</div>
				</div>
				<div class="mb-3">
					<label for="pushBody" class="form-label">Notification Body</label>
					<textarea
						class="form-control"
						id="pushBody"
						rows="2"
						placeholder="Enter your notification message here"
					></textarea>
				</div>
				<div class="mb-3">
					<label for="pushData" class="form-label"
						>Additional Data (JSON)</label
					>
					<textarea
						class="form-control"
						id="pushData"
						rows="2"
						placeholder='{"key": "value"}'
					></textarea>
				</div>
				<button type="submit" class="btn btn-danger">
					Send Push Notification <i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>
