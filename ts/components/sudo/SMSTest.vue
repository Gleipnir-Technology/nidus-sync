<template>
	<div class="card mb-4">
		<div class="card-header bg-primary text-white">
			<i class="bi bi-chat-dots"></i> SMS Testing
		</div>
		<div class="card-body">
			<form action="/sudo/sms" method="POST">
				<div class="mb-3">
					<label for="smsPhone" class="form-label"
						>Recipient Phone Number</label
					>
					<input
						type="tel"
						class="form-control"
						id="smsPhone"
						name="smsPhone"
						placeholder="+1 (555) 123-4567"
					/>
				</div>
				<div class="mb-3">
					<label for="smsMessage" class="form-label">Message</label>
					<textarea
						class="form-control"
						id="smsMessage"
						name="smsMessage"
						rows="3"
						maxlength="130"
						placeholder="Enter your SMS message here"
					></textarea>
				</div>
				<button type="submit" class="btn btn-primary">
					Send SMS <i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>
