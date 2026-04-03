<style scoped>
.photo-upload-area {
	border: 2px dashed #ccc;
	border-radius: 8px;
	padding: 20px;
	text-align: center;
	margin-bottom: 20px;
	background-color: #f9f9f9;
}

.photo-preview {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
	margin-top: 15px;
}

.photo-preview img {
	width: 100px;
	height: 100px;
	object-fit: cover;
	border-radius: 4px;
}

.btn-danger {
	font-size: 10px;
	padding: 0 5px;
}
</style>
<template>
	<div class="photo-upload-area">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			fill="currentColor"
			class="bi bi-camera mb-2"
			viewBox="0 0 16 16"
		>
			<path
				d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1v6zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4H2z"
			/>
			<path
				d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5zm0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7zM3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0z"
			/>
		</svg>
		<div
			class="file-upload-container"
			:style="{ backgroundColor: dropAreaBgColor }"
			@dragover.prevent="handleDragOver"
			@dragleave="handleDragLeave"
			@drop.prevent="handleDrop"
		>
			<input
				ref="fileInputRef"
				type="file"
				class="d-none"
				accept="image/jpeg,image/jpg,image/png,image/gif,image/webp,image/bmp"
				multiple
				@change="handleFileSelect"
			/>
			<button
				type="button"
				class="btn btn-outline-primary mb-2"
				@click="openFileDialog"
			>
				Add Photos
			</button>
		</div>
		<small class="d-block text-muted">
			Take pictures of the mosquito problem area
		</small>

		<!-- Photo Preview Area -->
		<div class="photo-preview mt-3 d-flex flex-wrap">
			<div
				v-for="image in modelValue"
				:key="image.id"
				class="position-relative m-1"
			>
				<img :src="image.preview" class="img-thumbnail" :alt="image.name" />
				<button
					type="button"
					class="btn btn-sm btn-danger position-absolute top-0 end-0"
					@click="removeImage(image.id)"
				>
					&times;
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";

export interface Image {
	id: number;
	file: File;
	name: string;
	preview: string;
}

interface Emits {
	(e: "update:modelValue", value: Image[]): void;
	(e: "fileAdded", image: Image): void;
	(e: "fileRemoved", imageId: number): void;
	(e: "filesDropped", files: File[]): void;
	(e: "error", error: string): void;
}

interface Props {
	modelValue: Image[];
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const fileInputRef = ref<HTMLInputElement | null>(null);
const dropAreaBgColor = ref("#f8f9fa");
let fileCounter = 0;

const openFileDialog = () => {
	fileInputRef.value?.click();
};

const handleDragOver = (e: DragEvent) => {
	dropAreaBgColor.value = "#e9ecef";
};

const handleDragLeave = () => {
	dropAreaBgColor.value = "#f8f9fa";
};

const handleDrop = (e: DragEvent) => {
	dropAreaBgColor.value = "#f8f9fa";

	if (e.dataTransfer?.files.length) {
		const files = Array.from(e.dataTransfer.files);
		emit("filesDropped", files);
		processFiles(files);
	}
};

const handleFileSelect = (e: Event) => {
	const target = e.target as HTMLInputElement;
	if (target.files && target.files.length > 0) {
		processFiles(Array.from(target.files));
	}
};

const processFiles = (files: File[]) => {
	const newImages: Image[] = [];

	files.forEach((file) => {
		if (!file.type.match("image.*")) {
			emit("error", `File ${file.name} is not an image`);
			return;
		}

		const reader = new FileReader();
		reader.onload = (e) => {
			const image: Image = {
				id: fileCounter++,
				file,
				name: file.name,
				preview: e.target?.result as string,
			};

			newImages.push(image);
			emit("fileAdded", image);

			// Update model after all files in this batch are processed
			if (
				newImages.length === files.filter((f) => f.type.match("image.*")).length
			) {
				emit("update:modelValue", [...props.modelValue, ...newImages]);
			}
		};
		reader.readAsDataURL(file);
	});
};

const removeImage = (imageId: number) => {
	const updatedImages = props.modelValue.filter((img) => img.id !== imageId);
	emit("update:modelValue", updatedImages);
	emit("fileRemoved", imageId);
};
</script>
