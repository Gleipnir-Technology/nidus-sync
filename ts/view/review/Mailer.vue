<style scoped>
.table th {
	font-weight: 600;
	border-bottom: 2px solid #dee2e6;
}

.badge {
	font-weight: 500;
	padding: 0.375rem 0.75rem;
}
</style>
<template>
	<div class="container-fluid py-4">
		<div class="row mb-4">
			<div class="col">
				<h1 class="h2">Mailers</h1>
				<p class="text-muted">Track the status of your postal mailers</p>
			</div>
		</div>

		<div class="row">
			<div class="col">
				<div class="card">
					<div class="card-body">
						<div class="table-responsive">
							<table class="table table-hover">
								<thead>
									<tr>
										<th>Created</th>
										<th>Status</th>
										<th>Site Address</th>
										<th>PDF</th>
									</tr>
								</thead>
								<tbody>
									<tr v-for="mailer in mailers" :key="mailer.id">
										<td>{{ formatDate(mailer.created) }}</td>
										<td>
											<span
												class="badge"
												:class="getStatusBadgeClass(mailer.status)"
											>
												{{ formatStatus(mailer.status) }}
											</span>
										</td>
										<td>
											<RouterLink
												:to="{
													name: 'Site Review',
													query: { site: mailer.site_id },
												}"
											>
												{{ formatAddressShort(mailer.address) }}
											</RouterLink>
										</td>
										<td>
											<a
												:href="mailer.pdfUrl()"
												target="_blank"
												class="btn btn-sm btn-outline-primary"
											>
												<i class="bi bi-file-pdf"></i> View PDF
											</a>
										</td>
									</tr>
								</tbody>
							</table>
						</div>

						<div
							v-if="mailers && mailers.length === 0"
							class="text-center py-5 text-muted"
						>
							<p>No mailers found</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { computedAsync } from "@vueuse/core";

import { formatDate, formatAddressShort } from "@/format";
import { useStoreMailer } from "@/store/mailer";
import { Mailer, MailerStatus } from "@/type/api";

const formatStatus = (status: MailerStatus): string => {
	return status.charAt(0).toUpperCase() + status.slice(1);
};

const getStatusBadgeClass = (status: MailerStatus): string => {
	const classes: Record<MailerStatus, string> = {
		created: "bg-secondary",
		printed: "bg-info",
		mailed: "bg-primary",
		completed: "bg-success",
	};
	return classes[status];
};
const storeMailer = useStoreMailer();
const mailers = computedAsync(async (): Promise<Mailer[]> => {
	return await storeMailer.list();
});
</script>
