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
										<td>{{ formatDate(mailer.createdAt) }}</td>
										<td>
											<span
												class="badge"
												:class="getStatusBadgeClass(mailer.status)"
											>
												{{ formatStatus(mailer.status) }}
											</span>
										</td>
										<td>
											<a :href="`/sites/${mailer.siteId}`">
												{{ mailer.siteAddress }}
											</a>
										</td>
										<td>
											<a
												:href="mailer.pdfUrl"
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
							v-if="mailers.length === 0"
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
import { onMounted, ref } from "vue";

import { useStoreMailer } from "@/store/mailer";

type MailerStatus = "created" | "printed" | "mailed" | "completed";

interface Mailer {
	id: string;
	createdAt: Date;
	status: MailerStatus;
	pdfUrl: string;
	siteId: string;
	siteAddress: string;
}

// Mock data
const mailers = ref<Mailer[]>([
	{
		id: "1",
		createdAt: new Date("2024-01-15T10:30:00"),
		status: "completed",
		pdfUrl: "/pdfs/mailer-1.pdf",
		siteId: "101",
		siteAddress: "123 Main St, Springfield, IL 62701",
	},
	{
		id: "2",
		createdAt: new Date("2024-01-16T14:20:00"),
		status: "mailed",
		pdfUrl: "/pdfs/mailer-2.pdf",
		siteId: "102",
		siteAddress: "456 Oak Ave, Chicago, IL 60601",
	},
	{
		id: "3",
		createdAt: new Date("2024-01-17T09:15:00"),
		status: "printed",
		pdfUrl: "/pdfs/mailer-3.pdf",
		siteId: "103",
		siteAddress: "789 Pine Rd, Naperville, IL 60540",
	},
	{
		id: "4",
		createdAt: new Date("2024-01-18T11:45:00"),
		status: "created",
		pdfUrl: "/pdfs/mailer-4.pdf",
		siteId: "104",
		siteAddress: "321 Elm St, Peoria, IL 61602",
	},
	{
		id: "5",
		createdAt: new Date("2024-01-18T16:00:00"),
		status: "mailed",
		pdfUrl: "/pdfs/mailer-5.pdf",
		siteId: "105",
		siteAddress: "654 Maple Dr, Rockford, IL 61101",
	},
]);

const formatDate = (date: Date): string => {
	return new Intl.DateTimeFormat("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	}).format(date);
};

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
onMounted(() => {
	storeMailer.list();
});
</script>
