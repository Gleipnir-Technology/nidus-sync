<template>
	<!-- SSE Testing -->
	<div class="card mb-5">
		<div class="card-header bg-danger text-white">
			<i class="bi bi-bell"></i> Server-sent event testing
		</div>
		<div class="card-body">
			<form action="/sudo/sse" method="POST">
				<div class="row">
					<div class="col-md-6">
						<div class="mb-3">
							<label for="organizationID" class="form-label"
								>Organization ID</label
							>
							<input
								class="form-control"
								id="organization-id"
								name="organizationID"
								placeholder="Organization ID"
								type="text"
								:value="session"
							/>
						</div>
						<div class="mb-3">
							<label for="type" class="form-label">Type</label>
							<select class="form-select" id="type" name="type">
								<option value="created">Created</option>
								<option value="deleted">Deleted</option>
								<option value="heartbeat">Heartbeat</option>
								<option value="sudo">Sudo</option>
								<option value="updated">Updated</option>
							</select>
						</div>
					</div>
					<div class="col-md-6">
						<div class="mb-3">
							<label for="resource" class="form-label">Resource</label>
							<input
								class="form-control"
								id="resource"
								name="resource"
								type="text"
								value="rmo:nuisance"
							/>
						</div>
						<div class="mb-3">
							<label for="uriPath" class="form-label">URI path</label>
							<input
								class="form-control"
								id="uri-path"
								name="uriPath"
								type="text"
								value="/report/abcd-1234"
							/>
						</div>
					</div>
				</div>
				<button type="submit" class="btn btn-danger">
					Send SSE<i class="bi bi-send"></i>
				</button>
			</form>
		</div>
	</div>
</template>

<script setup lang="ts">
import { useSessionStore } from "@/store/session";

const session = useSessionStore();
</script>
