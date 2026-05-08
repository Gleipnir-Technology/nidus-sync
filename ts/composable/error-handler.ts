import { ref } from "vue";

interface ErrorState {
	hasError: boolean;
	message: string;
	timestamp: Date;
}

const globalError = ref<ErrorState | null>(null);

export function useErrorHandler() {
	const setError = (error: Error) => {
		globalError.value = {
			hasError: true,
			message: error.message,
			timestamp: new Date(),
		};
	};

	const errorClear = () => {
		globalError.value = null;
	};

	return {
		error: globalError,
		setError,
		errorClear,
	};
}
