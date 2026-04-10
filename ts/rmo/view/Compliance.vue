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
					@doEvidence="doEvidence"
					@doContact="doContact"
					@doLocator="doLocator"
					@doPermission="doPermission"
					v-model="compliance"
				/>
			</LoadingOverlay>
		</router-view>
	</template>
	<template v-else>
		<p>loading {{ slug }}...</p>
	</template>
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
import type { District, PublicReport } from "@/type/api";
import { Address, Location, PermissionAccess } from "@/type/api";
import { Locator } from "@/type/map";
import { type Contact } from "@/rmo/content/compliance/Contact.vue";
import { type Permission } from "@/rmo/content/compliance/Permission.vue";

export interface Compliance {
	comments: string;
	contact: Contact;
	id: string;
	images: Image[];
	location: Location;
	locator: Locator;
	permission: Permission;
	uri: string;
}
interface Props {
	slug: string;
}

const districtStore = useStoreDistrict();

const compliance = ref<Compliance>({
	comments: "",
	contact: {
		name: "",
		phone: "",
		can_text: true,
		email: "",
	},
	id: "",
	images: [],
	location: {
		latitude: 0,
		longitude: 0,
	},
	locator: {
		address: new Address(),
		location: new Location(),
	},
	permission: {
		access: PermissionAccess.UNSELECTED,
		access_instructions: "",
		availability_notes: "",
		gate_code: "",
		has_dog: false,
		wants_scheduled: false,
	},
	uri: "",
});
const isLoading = ref<boolean>(true);
const props = defineProps<Props>();
const report = ref<PublicReport | null>();
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
		formData.append("longitude", "0");
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
}
function doEvidence() {
	console.log("evidence", compliance.value);
}
function doContact() {
	console.log("contact", compliance.value.contact);
}
function doLocator() {
	console.log("locator done", compliance.value.locator);
	updateReport({
		locator: compliance.value.locator,
	});
}
function doPermission() {
	console.log("permission", compliance.value.permission);
}
interface ComplianceUpdate {
	locator?: Locator;
}
async function updateReport(updates: ComplianceUpdate) {
	const resp = await fetch(compliance.value.uri, {
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
onMounted(() => {
	const client_id = storeLocal.getClientID();
	storeLocation
		.get()
		.then((loc: GeolocationPosition) => {
			compliance.value.location = loc.coords;
			createReport(client_id, loc);
		})
		.catch((e) => {
			console.log("failed to get location", e);
			createReport(client_id);
		});
});
</script>
