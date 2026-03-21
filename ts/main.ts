import Alpine from './vendor/alpinejs-3.15.8.js';
import { createApp } from 'vue';
import App from './app.vue';
import { SSEManager } from './sse-manager';
import { SetupSidebar } from "./sidebar";
import 'maplibre-gl/dist/maplibre-gl.css';

// Import Bootstrap Icons CSS
import 'bootstrap-icons/font/bootstrap-icons.css';
// Import Bootstrap SCSS
import './style/style.scss';

// Import Bootstrap JavaScript and make it available globally
import * as bootstrap from 'bootstrap';
window.bootstrap = bootstrap;

import { Planning } from './app/planning';

// Make Alpine available on window for inline Alpine
window.Alpine = Alpine;

// Make SSEManager available to all the JavaScript
window.SSEManager = SSEManager;

function createAppPlanning() {
	const app = createApp({
		data() {
			return {
				count: 0
			}
		}
	});
}
window.createAppPlanning = createAppPlanning;

// Wait for DOM to be ready, then initialize Alpine
document.addEventListener("DOMContentLoaded", () => {
	Alpine.start();
	SSEManager.connect("/api/events");
	SetupSidebar();
});
	document.addEventListener("alpine:init", () => {
		 const user = {
			"display_name":"",
			"initials":"",
			"notifications":[],
			"notification_counts":{
				"communication":0,
				"home":0,
				"review":0
			},
			"organization":{
				"name":""
			},
			"role":"",
			"username":""
		};
		 Alpine.store("user",user);
	})
interface GreetingComponent {
    message: string;
    name: string;
    updateMessage(): void;
}

Alpine.data('greeting', (): GreetingComponent => ({
    message: 'Welcome to Alpine + TypeScript!',
    name: 'World',
    
    updateMessage() {
        this.message = 'Message updated at ' + new Date().toLocaleTimeString();
    }
}));

createApp(App).mount('#app');
