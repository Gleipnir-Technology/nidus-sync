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
	<template v-if="district">
		<router-view v-slot="{ Component }">
			<LoadingOverlay
				:is-loading="isLoading"
				loading-text="Loading previous data"
			>
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
			</LoadingOverlay>
		</router-view>
	</template>
	<template v-else>
		<p>loading {{ slug }}...</p>
	</template>
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
	slug: string;
}

const districtStore = useStoreDistrict();

const isLoading = ref<boolean>(true);
const isUploading = ref<boolean>(false);
const props = defineProps<Props>();
const report = ref<PublicReportCompliance>(new PublicReportCompliance());
const district = computedAsync(async (): Promise<District | undefined> => {
	const districts = await districtStore.list();
	return districts.find((district: District) => district.slug == props.slug);
});
const storeLocal = useStoreLocal();
const storeLocation = useStoreLocation();
async function createReport(client_id: string, loc?: GeolocationPosition) {
	const formData = new FormData();
	formData.append("client_id", client_id);
	if (loc) {
		formData.append("location.accuracy", loc.coords.accuracy.toString());
		formData.append("location.latitude", loc.coords.latitude.toString());
		formData.append("location.longitude", loc.coords.longitude.toString());
	} else {
		formData.append("location.accuracy", "0");
		formData.append("location.latitude", "0");
		formData.append("location.longitude", "0");
	}
	const resp = await fetch("/api/rmo/compliance", {
		method: "POST",
		body: formData,
		// Don't set Content-Type, the borwser should do it
	});
	if (!resp.ok) {
		const content = await resp.text();
		console.error("Failed to create compliance report", resp.status, content);
		return;
	}
	const body = await resp.json();
	storeLocal.setExistingComplianceReportURI(body.uri);
	report.value!.uri = body.uri;
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
async function fetchExistingReport(report_uri: string) {
	isLoading.value = true;
	const resp = await fetch(report_uri);
	if (!resp.ok) {
		isLoading.value = false;
		const content = await resp.text();
		console.error(
			"Failed to fetch existing report",
			report_uri,
			resp.status,
			content,
		);
		return;
	}
	const body = (await resp.json()) as PublicReportComplianceOptions;
	Object.assign(report.value, new PublicReportCompliance(body));
	console.log("fetched existing report", report.value);
	isLoading.value = false;
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
}
onMounted(() => {
	const client_id = storeLocal.getClientID();
	const report_uri = storeLocal.getExistingComplianceReportURI();
	if (report_uri) {
		fetchExistingReport(report_uri);
	} else {
		isLoading.value = false;
		createReport(client_id);
	}
	storeLocation
		.get()
		.then((loc: GeolocationPosition) => {
			report.value.location = loc.coords;
			updateReport({
				location: report.value.location,
			});
		})
		.catch((e) => {
			console.log("failed to get location", e);
		});
});
</script>
