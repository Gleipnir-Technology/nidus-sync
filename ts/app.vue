<script setup lang="ts">
import { ref } from 'vue';
import MapView from './components/map-view.vue';

const count = ref(0);
const message = ref('Hello from Vue 3!');
</script>

<template>
  <div>
    <p>{{ message }}</p>
    <button @click="count++">Count: {{ count }}</button>
    
    <h2>Map Example:</h2>
    <MapView />
  </div>
</template>

<style scoped>
button {
  padding: 0.5rem 1rem;
  font-size: 1rem;
}
</style>
