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
		<LoadingOverlay :is-loading="isLoading" loading-text="Loading report">
			<template v-if="!isLoading">
				<component
					:is="Component"
					:district="district"
					@doAddress="doAddress"
					@doContact="doContact"
					@doEvidence="doEvidence"
					@doPermission="doPermission"
					@doSubmit="doSubmit"
					:publicID="report?.public_id ?? 'unset'"
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
import { useStorePublicReport } from "@/rmo/store/publicreport";
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

const isLoading = ref<boolean>(true);
const isUploading = ref<boolean>(false);
const props = defineProps<Props>();
const report = ref<PublicReportCompliance | undefined>(undefined);
const district = ref<District | undefined>(undefined);
const storeDistrict = useStoreDistrict();
const storeLocal = useStoreLocal();
const storeLocation = useStoreLocation();
const storePublicReport = useStorePublicReport();
async function createReport(client_id: string, district_uri: string) {
	let content = {
		client_id: client_id,
		district: district_uri,
		location: {
			accuracy: 0,
			latitude: 0,
			longitude: 0,
		},
	};
	const resp = await fetch("/api/rmo/compliance", {
		body: JSON.stringify(content),
		headers: {
			"Content-Type": "application/json",
		},
		method: "POST",
	});
	if (!resp.ok) {
		const content = await resp.text();
		console.error("Failed to create compliance report", resp.status, content);
		return;
	}
	const body = await resp.json();
	storeLocal.setExistingComplianceReportURI(body.uri);
	report.value!.public_id = body.public_id;
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
async function doMounted() {
	const r = await storePublicReport.byID(props.public_id);
	report.value = r as PublicReportCompliance;
	const d = await storeDistrict.byURI(r.district);
	district.value = d;
	isLoading.value = false;
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
async function updateReport(updates: ComplianceUpdate) {
	if (!(report.value && report.value.uri)) {
		console.log("Refusing to update report without URI");
		return;
	}
	storePublicReport.update(report.value.uri, updates);
}
async function updateLocation() {
	if (!report.value) return;
	const loc = await storeLocation.get();
	report.value.location = loc.coords;
	updateReport({
		location: report.value.location,
	});
}
async function uploadImages(images: Image[]) {
	if (images.length == 0) return;
	if (!report.value) return;

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
	const r = await storePublicReport.fetchByURI(report.value.uri);
	Object.assign(report.value, r);
}
onMounted(() => {
	doMounted();
});
</script>
