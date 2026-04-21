<template>
	<PublicReportLoading />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

import { useStoreDistrict } from "@/rmo/store/district";
import { useStoreLocal } from "@/store/local";
import { useStorePublicReport } from "@/store/publicreport";
import PublicReportLoading from "@/rmo/components/PublicReportLoading.vue";

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
