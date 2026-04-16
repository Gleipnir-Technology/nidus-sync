<template>
	<p>loading</p>
</template>
<script setup lang="ts">
interface Emits {
	(e: "doSelectTask", id: number): void;
}
interface Props {}
const emit = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {});
</script>
