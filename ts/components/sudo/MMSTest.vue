<template>
	<!-- MMS Testing -->
	<div class="card mb-4">
		<div class="card-header bg-success text-white">
			<i class="bi bi-image"></i> MMS Testing
		</div>
		<div class="card-body">
			<form>
				<div class="mb-3">
					<label for="mmsPhone" class="form-label"
						>Recipient Phone Number</label
					>
					<input
						type="tel"
						class="form-control"
						id="mmsPhone"
						placeholder="+1 (555) 123-4567"
					/>
				</div>
				<div class="mb-3">
					<label for="mmsMessage" class="form-label">Message</label>
					<textarea
						class="form-control"
						id="mmsMessage"
						rows="2"
						placeholder="Enter your MMS message here"
					></textarea>
				</div>
				<div class="mb-3">
					<label for="mmsImage" class="form-label">Attach Image</label>
					<input
						class="form-control"
						type="file"
						id="mmsImage"
						accept="image/*"
					/>
				</div>
				<button type="submit" class="btn btn-success">
					Send MMS <i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>
