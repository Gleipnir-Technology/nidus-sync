<template>
	<!-- RCS Testing -->
	<div class="card mb-4">
		<div class="card-header bg-info text-white">
			<i class="bi bi-chat-square-text"></i> RCS Testing
		</div>
		<div class="card-body">
			<form>
				<div class="mb-3">
					<label for="rcsPhone" class="form-label"
						>Recipient Phone Number</label
					>
					<input
						type="tel"
						class="form-control"
						id="rcsPhone"
						placeholder="+1 (555) 123-4567"
					/>
				</div>
				<div class="mb-3">
					<label for="rcsMessage" class="form-label">Message</label>
					<textarea
						class="form-control"
						id="rcsMessage"
						rows="3"
						placeholder="Enter your RCS message here"
					></textarea>
				</div>
				<div class="mb-3">
					<label class="form-label">Interactive Options</label>
					<div class="form-check">
						<input
							class="form-check-input"
							type="checkbox"
							id="rcsIncludeButtons"
						/>
						<label class="form-check-label" for="rcsIncludeButtons"
							>Include Action Buttons</label
						>
					</div>
				</div>
				<button type="submit" class="btn btn-info">
					Send RCS <i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>
