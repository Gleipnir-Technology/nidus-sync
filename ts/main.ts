import Alpine from './vendor/alpinejs-3.15.8.js';

// Make Alpine available on window for inline Alpine
window.Alpine = Alpine;

// Wait for DOM to be ready, then initialize Alpine
document.addEventListener("DOMContentLoaded", () => {
	Alpine.start();
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
