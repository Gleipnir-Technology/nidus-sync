<style scoped>
.modal.show {
	background-color: rgba(0, 0, 0, 0.5);
}
img {
	min-width: 512px;
	min-height: 512px;
}
</style>
<template>
	<div
		class="modal fade"
		:class="{ 'show d-block': show }"
		tabindex="-1"
		v-show="show"
		@click.self="emit('close')"
	>
		<div class="modal-dialog modal-lg modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Image Viewer</h5>
					<button
						type="button"
						class="btn-close"
						@click="emit('close')"
					></button>
				</div>
				<div class="modal-body text-center">
					<div v-if="image && show">
						<img
							:src="image.src"
							class="img-fluid rounded"
							style="max-height: 60vh"
						/>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div
		class="modal-backdrop fade show"
		v-show="show"
		@click="emit('close')"
	></div>
</template>

<script setup lang="ts">
interface Emits {
	(e: "close"): void;
}
export interface Image {
	src: string;
}
interface Props {
	image?: Image;
	show: boolean;
}
const emit = defineEmits<Emits>();
const props = defineProps<Props>();
</script>
