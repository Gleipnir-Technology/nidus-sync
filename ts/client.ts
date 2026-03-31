// src/api/axios.ts
import axios, { AxiosInstance } from "axios";
import router from "@/router";

// Extend the AxiosInstance interface
declare module "axios" {
	interface AxiosInstance {
		isAuthenticated(): boolean;
	}
}

const apiClient = axios.create({
	baseURL: "/api",
	withCredentials: true,
});

apiClient.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response && error.response.status === 401) {
			router.push("/login");
		}
		return Promise.reject(error);
	},
);

apiClient.isAuthenticated = () => {
	return true;
};

export default apiClient;
