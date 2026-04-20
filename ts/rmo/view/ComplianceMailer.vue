<style scoped>
body {
	background-color: #f8f9fa;
	min-height: 100vh;
	display: flex;
	flex-direction: column;
}
body > .container-fluid {
	flex: 1;
}
.progress-bar {
	background-color: #0d6efd;
	transition: width 0.3s ease;
}
.reference-number {
	text-align: center;
	color: #6c757d;
	font-size: 0.9rem;
	margin-top: 24px;
}
</style>
<template>
	<router-view v-slot="{ Component }">
		<LoadingOverlay
			:is-loading="isLoading"
			loading-text="Loading previous data"
		>
			<template v-if="!isLoading">
				<component
					:is="Component"
					:district="district"
					@doAddress="doAddress"
					@doContact="doContact"
					@doEvidence="doEvidence"
					@doPermission="doPermission"
					@doSubmit="doSubmit"
					v-model="report"
				/>
			</template>
		</LoadingOverlay>
	</router-view>
	<!-- Reference Number -->
	<div class="reference-number" v-if="report && report.public_id">
		<small>
			Reference number: <strong>{{ report.public_id }}</strong>
		</small>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { computedAsync } from "@vueuse/core";

import type { Image } from "@/components/ImageUpload.vue";
import { useStoreDistrict } from "@/rmo/store/district";
import { useStoreLocal } from "@/store/local";
import { useStoreLocation } from "@/store/location";
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

const isLoading = ref<boolean>(true);
const isUploading = ref<boolean>(false);
const props = defineProps<Props>();
const report = ref<PublicReportCompliance>(new PublicReportCompliance());
const district = ref<District | undefined>(undefined);
const storeLocal = useStoreLocal();
const storeLocation = useStoreLocation();
async function beginReport(client_id: string) {
	const report_uri = "/api/publicreport/compliance/" + props.public_id;
	const [districts, r] = await Promise.all([
		districtStore.list(),
		fetchExistingReport(report_uri),
	]);
	Object.assign(report.value, r);
	const d = districts.find((district: District) => district.uri == r.district);
	if (!d) {
		console.error("Failed to find district with uri", districts, r.district);
		return;
	}
	district.value = d;
	isLoading.value = false;
	await updateLocation();
}
function doAddress() {
	if (!report.value) {
		console.log("can't do address, null report");
		return;
	}
	console.log("address done", report.value.address);
	updateReport({
		address: report.value.address,
	});
}
function doEvidence(images: Image[]) {
	if (!report.value) {
		console.log("can't do evidence, null report");
		return;
	}
	uploadImages(images);
	if (report.value.comments) {
		updateReport({
			comments: report.value.comments,
		});
	}
}
function doContact() {
	if (!report.value) {
		console.log("can't do contact, null report");
		return;
	}
	console.log(
		"contact",
		JSON.stringify(report.value.reporter),
		report.value.reporter,
	);
	updateReport({
		reporter: report.value.reporter,
	});
}
function doPermission() {
	if (!report.value) {
		console.log("can't do permission, null report");
		return;
	}
	console.log("report.value.has_dog", report.value.has_dog);
	updateReport({
		access_instructions: report.value.access_instructions,
		availability_notes: report.value.availability_notes,
		gate_code: report.value.gate_code,
		has_dog: report.value.has_dog,
		permission_type: report.value.permission_type,
		wants_scheduled: report.value.wants_scheduled,
	});
}
function doSubmit() {
	console.log("submit", report.value);
	storeLocal.delExistingComplianceReportURI();
}
async function fetchExistingReport(
	report_uri: string,
): Promise<PublicReportCompliance> {
	const resp = await fetch(report_uri);
	if (!resp.ok) {
		const content = await resp.text();
		throw new Error(
			`Failed to fetch existing report ${report_uri}: ${resp.status} ${content}`,
		);
	}
	const body = (await resp.json()) as PublicReportComplianceOptions;
	console.log("fetched existing report", report.value);
	return new PublicReportCompliance(body);
}
async function updateReport(updates: ComplianceUpdate) {
	if (!report.value.uri) {
		console.log("Refusing to update report without URI");
		return;
	}
	const resp = await fetch(report.value.uri, {
		method: "PUT",
		body: JSON.stringify(updates),
		headers: {
			"Content-Type": "application/json",
		},
	});
	if (!resp.ok) {
		const content = await resp.text();
		console.error("Failed to update compliance", resp.status, content);
		return;
	}
}
async function updateLocation() {
	const loc = await storeLocation.get();
	report.value.location = loc.coords;
	updateReport({
		location: report.value.location,
	});
}
async function uploadImages(images: Image[]) {
	if (images.length == 0) return;
	isUploading.value = true;
	const formData = new FormData();
	images.map(async (image, index) => {
		formData.append(`image[${index}]`, image.file, image.name);
	});
	const url = `${report.value.uri}/image`;
	const response = await fetch(url, {
		body: formData,
		method: "POST",
	});
	if (!response.ok) {
		const content = await response.text();
		console.error(
			"Failed to POST images",
			url,
			response.status,
			response.statusText,
			content,
		);
		isUploading.value = false;
		return;
	}
	isUploading.value = false;
	// after everything is done update the report so that we see the correct number of images
	// on the report summary
	await fetchExistingReport(report.value.uri);
}
onMounted(() => {
	const client_id = storeLocal.getClientID();
	beginReport(client_id);
});
</script>
