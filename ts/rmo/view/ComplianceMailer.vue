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
import { useStoreLocation } from "@/store/location";
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
import { Contact, Address, Location, PermissionType } from "@/type/api";

interface Props {
	public_id: string;
}

const districtStore = useStoreDistrict();

const props = defineProps<Props>();
const router = useRouter();
const storeLocal = useStoreLocal();
const storePublicReport = useStorePublicReport();
async function doMounted() {
	const client_id = storeLocal.getClientID();
	const report_uri = storeLocal.getExistingComplianceReportURI();
	if (report_uri && report_uri.endsWith(props.public_id)) {
		console.log("Loading previous report", report_uri);
	} else {
		const report = await storePublicReport.createCompliance({
			client_id: client_id,
			mailer_id: props.public_id,
		});
		storeLocal.setExistingComplianceReportURI(report.uri);
		console.log("Created new compliance report", report);
	}
	router.replace(`/compliance/${props.public_id}`);
}
onMounted(() => {
	doMounted();
});
</script>
