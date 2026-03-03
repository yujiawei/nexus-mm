import axios from 'axios';
import { useToastStore } from '../components/common/Toast';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    } else if (error.response?.status === 429) {
      useToastStore.getState().addToast('Too many requests. Please slow down.', 'error');
    } else if (error.response?.status >= 500) {
      useToastStore.getState().addToast('Server error. Please try again later.', 'error');
    }
    return Promise.reject(error);
  }
);

export default client;
