// src/api/axios.js or similar
import axios from 'axios';
import router from '@/router';

const apiClient = axios.create({
    baseURL: '/api',
    withCredentials: true
});

// Response interceptor to catch auth failures
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            // Session expired or not authenticated
            router.push('/login');
        }
        return Promise.reject(error);
    }
);
apiClient.isAuthenticated = () => {
	return true;
}
export default apiClient;
