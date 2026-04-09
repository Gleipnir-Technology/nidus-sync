<template>
	<div class="row">
		<div class="col-md-4">
			<strong>{{ label }}</strong>
		</div>
		<div class="col-md-6">
			<span>{{ value }}</span>
		</div>
	</div>
</template>
<script setup lang="ts">
interface Props {
	label: string;
	value: string;
}
const props = defineProps<Props>();
</script>
