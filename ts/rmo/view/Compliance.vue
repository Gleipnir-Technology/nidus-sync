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
			<component
				:is="Component"
				:district="district"
				@doComments="doComments"
				@doContact="doContact"
				@doImages="doImages"
				@doLocator="doLocator"
				@doPermission="doPermission"
				v-model="compliance"
			/>
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
import Intro from "@/rmo/content/compliance/Intro.vue";
import type { District } from "@/type/api";
import { Locator } from "@/type/map";
import { type Contact } from "@/rmo/content/compliance/Contact.vue";
import { type Permission } from "@/rmo/content/compliance/Permission.vue";

export interface Compliance {
	comments: string;
	contact?: Contact;
	locator: Locator;
	images: Image[];
	permission?: Permission;
}
interface Props {
	slug: string;
}

const districtStore = useStoreDistrict();

const compliance = ref<Compliance>({
	comments: "",
	images: [],
	locator: {
		address: {
			country: "",
			gid: "",
			locality: "",
			number: "",
			postal_code: "",
			raw: "",
			region: "",
			street: "",
			unit: "",
		},
		location: {
			latitude: 0,
			longitude: 0,
		},
	},
});
const props = defineProps<Props>();
const district = computedAsync(async (): Promise<District | undefined> => {
	const districts = await districtStore.list();
	return districts.find((district: District) => district.slug == props.slug);
});
function doComments(comments: string) {
	console.log("comments", comments);
	compliance.value.comments = comments;
}
function doContact(contact: Contact | undefined) {
	console.log("contact", contact);
	compliance.value.contact = contact;
}
function doImages(images: Image[]) {
	console.log("images", images);
	compliance.value.images = images;
}
function doLocator() {
	console.log("locator done", compliance.value.locator);
}
function doPermission(permission: Permission | undefined) {
	console.log("permission", permission);
	compliance.value.permission = permission;
}
</script>
