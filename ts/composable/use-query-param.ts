import { computed, ComputedRef } from "vue";
import { useRoute, useRouter, LocationQueryValue } from "vue-router";

export function useQueryParam(paramName: string) {
	const route = useRoute();
	const router = useRouter();

	// Returns string | null for easier handling
	const value = computed<string | null>(() => {
		const param = route.query[paramName];

		// Handle arrays by taking first value, or return null
		if (Array.isArray(param)) {
			return param[0] ?? null;
		}

		return param ?? null;
	});

	const setValue = (newValue: string | number) => {
		router.replace({
			name: route.name,
			query: {
				...route.query,
				[paramName]: String(newValue),
			},
		});
	};

	const removeValue = () => {
		const { [paramName]: _, ...rest } = route.query;
		router.replace({
			name: route.name,
			query: rest,
		});
	};

	return {
		value,
		setValue,
		removeValue,
	};
}
