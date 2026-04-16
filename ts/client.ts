// src/api/axios.ts
import axios, { AxiosInstance, AxiosRequestConfig } from "axios";
import router from "@/router";

class ApiClient {
	private client: AxiosInstance;
	private _isAuthenticated: boolean = false;

	constructor() {
		this.client = axios.create({
			timeout: 10000,
			withCredentials: true,
		});

		// Request interceptor for auth headers, content-type, etc.
		this.client.interceptors.request.use((config) => {
			// Content-type negotiation
			config.headers["Accept"] = "application/json";
			config.headers["X-Requested-With"] = "nidus-web 0.1";

			// Add auth token if logged in
			const token = localStorage.getItem("authToken");
			if (token) {
				config.headers["Authorization"] = `Bearer ${token}`;
			}

			return config;
		});

		// Response interceptor for handling auth errors
		this.client.interceptors.response.use(
			(response) => response,
			(error) => {
				if (error.response?.status === 401) {
					this._isAuthenticated = false;
					// Could emit event or redirect here
				}
				return Promise.reject(error);
			},
		);
	}

	get isAuthenticated(): boolean {
		return this._isAuthenticated;
	}

	setLoggedIn(value: boolean): void {
		this._isAuthenticated = value;
	}

	async JSONGet<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
		const response = await this.client.get<T>(url, {
			...config,
			headers: {
				Accept: "application/json",
				...config?.headers,
			},
		});
		return response.data;
	}

	async JSONPost<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig,
	): Promise<T> {
		const response = await this.client.post<T>(url, data, {
			...config,
			headers: {
				"Content-Type": "application/json",
				Accept: "application/json",
				...config?.headers,
			},
		});
		return response.data;
	}

	async JSONPut<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig,
	): Promise<T> {
		const response = await this.client.put<T>(url, data, {
			...config,
			headers: {
				"Content-Type": "application/json",
				Accept: "application/json",
				...config?.headers,
			},
		});
		return response.data;
	}

	async JSONDelete<T = any>(
		url: string,
		config?: AxiosRequestConfig,
	): Promise<T> {
		const response = await this.client.delete<T>(url, config);
		return response.data;
	}
}

// Single instance export - this IS the singleton
export const apiClient = new ApiClient();
