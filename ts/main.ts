import Alpine from './vendor/alpinejs-3.15.8.js';
import bootstrap from '../static/vendor/bootstrap-5.3.8/bootstrap.bundle.min.js';
import { SSEManager } from './sse-manager';

// Make Alpine available on window for inline Alpine
window.Alpine = Alpine;

// Make bootstrap available on window for various scripts
window.bootstrap = bootstrap;

// Make SSEManager available to all the JavaScript
window.SSEManager = SSEManager;

// Wait for DOM to be ready, then initialize Alpine
document.addEventListener("DOMContentLoaded", () => {
	Alpine.start();
	SSEManager.connect("/api/events");
});
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
