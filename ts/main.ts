import { createApp } from 'vue';

// Simple example without SFCs
const App = {
  data() {
    return {
      count: 0,
      message: 'Hello from Vue 3!'
    }
  },
  template: `
    <div>
      <p>{{ message }}</p>
      <button @click="count++">Count: {{ count }}</button>
    </div>
  `
};

createApp(App).mount('#app');
