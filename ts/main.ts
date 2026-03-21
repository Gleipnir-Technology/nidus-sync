import Alpine from './vendor/alpinejs-3.15.8.js';

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

Alpine.start();

