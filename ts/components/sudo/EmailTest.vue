<template>
	<!-- Email Testing -->
	<div class="card mb-4">
		<div class="card-header bg-warning text-dark">
			<i class="bi bi-envelope"></i> Email Testing
		</div>
		<div class="card-body">
			<form action="/sudo/email" method="POST">
				<div class="mb-3">
					<label for="emailFrom" class="form-label">From Account</label>
					<select class="form-select" id="emailFrom" name="emailFrom">
						<option value="RMO address"></option>
						<option value="Sync address"></option>
					</select>
				</div>
				<div class="mb-3">
					<label for="emailTo" class="form-label">Recipient Email</label>
					<input
						type="email"
						class="form-control"
						id="emailTo"
						name="emailTo"
						placeholder="user@example.com"
					/>
				</div>
				<div class="mb-3">
					<label for="emailSubject" class="form-label">Subject</label>
					<input
						type="text"
						class="form-control"
						id="emailSubject"
						name="emailSubject"
						placeholder="Email Subject"
					/>
				</div>
				<div class="mb-3">
					<label for="emailBody" class="form-label">Message</label>
					<textarea
						class="form-control"
						id="emailBody"
						name="emailBody"
						placeholder="Enter your email message here"
						rows="3"
					></textarea>
				</div>
				<button type="submit" class="btn btn-warning text-dark">
					Send Email <i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>
