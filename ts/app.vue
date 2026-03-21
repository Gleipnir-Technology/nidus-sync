<script setup lang="ts">
import { ref } from 'vue';

const count = ref(0);
const message = ref('Hello from Vue 3!');
</script>

<template>
  <div>
    <p>{{ message }}</p>
    <button @click="count++">Count: {{ count }}</button>
  </div>
</template>

<style scoped>
button {
  padding: 0.5rem 1rem;
  font-size: 1rem;
}
</style>
