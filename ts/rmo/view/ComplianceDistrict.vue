<template>
	<div class="container">
		<div class="row min-vh-100 align-items-center justify-content-center">
			<div class="col-auto text-center">
				<div class="spinner-border text-primary" role="status">
					<span class="visually-hidden">Loading...</span>
				</div>
				<p class="mt-3 text-muted">Loading report details...</p>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { computedAsync } from "@vueuse/core";

import type { Image } from "@/components/ImageUpload.vue";
import { useStoreDistrict } from "@/rmo/store/district";
import { useStoreLocal } from "@/store/local";
import { useStorePublicReport } from "@/store/publicreport";
import Intro from "@/rmo/content/compliance/Intro.vue";
import LoadingOverlay from "@/components/LoadingOverlay.vue";
import {
	type ComplianceUpdate,
	type District,
	PublicReport,
	PublicReportCompliance,
	PublicReportComplianceOptions,
} from "@/type/api";
import { Contact, Address, PermissionType } from "@/type/api";

interface Props {
	slug: string;
}

const districtStore = useStoreDistrict();

const props = defineProps<Props>();
const router = useRouter();
const storeLocal = useStoreLocal();
const storePublicReport = useStorePublicReport();
async function doMounted() {
	const client_id = storeLocal.getClientID();
	const report_uri = storeLocal.getExistingComplianceReportURI();
	if (report_uri) {
		const report = await storePublicReport.byURI(report_uri);
		if (report && report.public_id) {
			router.replace(`/compliance/${report.public_id}`);
			return;
		}
	}
	const districts = await districtStore.list();
	const district = districts.find(
		(district: District) => district.slug == props.slug,
	);
	if (!district) {
		console.error("failed to find matching district", props.slug, districts);
		return;
	}
	const report = await storePublicReport.createCompliance({
		client_id: client_id,
		district: district.uri,
	});
	storeLocal.setExistingComplianceReportURI(report.uri);
	router.replace(`/compliance/${report.public_id}`);
	console.log("Created new compliance report", report);
}
onMounted(() => {
	doMounted();
});
</script>
