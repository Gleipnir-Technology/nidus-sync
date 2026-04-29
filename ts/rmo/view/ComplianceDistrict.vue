<template>
	<PublicReportLoading />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

import { useStoreDistrict } from "@/rmo/store/district";
import { useStoreLocal } from "@/store/local";
import { useStorePublicReport } from "@/rmo/store/publicreport";
import PublicReportLoading from "@/rmo/components/PublicReportLoading.vue";
import { type District } from "@/type/api";

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
