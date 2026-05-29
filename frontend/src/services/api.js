import axios from 'axios';
// import { API_BASE_URL, ADMIN_API_KEY } from '../config';
import { API_BASE_URL, VITE_DEBUG_MODE } from '../config';
import { getToken } from '../utils/storage';

const api = axios.create({
  baseURL: API_BASE_URL || 'http://localhost:8080/api',
  withCredentials: true,
});

// VULNERABILITY: Logging sensitive data
api.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    {/*
      // VULNERABILITY: Logging requests with sensitive data
      if (console && console.log) {
        console.log('API Request:', config.method, config.url, config.data);
      }
       */}
    
    return config;
  },
  (error) => {
    if (VITE_DEBUG_MODE) {
      console.error('Request Error:', error);
    }
    return Promise.reject(error);
  }
);

// VULNERABILITY: Logging sensitive response data
api.interceptors.response.use(
  (response) => {
    // console.log('API Response:', response.data);
    return response;
  },
  (error) => {
    // console.error('Response Error:', error.response?.data);
    return Promise.reject(error);
  }
);

// Auth APIs
export const register = (email, password, name) => {
  return api.post('/auth/register', { email, password, name });
};

export const login = (email, password) => {
  return api.post('/auth/login', { email, password });
};

// Task APIs
export const getTasks = () => {
  return api.get('/tasks');
};

export const createTask = (taskData) => {
  return api.post('/tasks', taskData);
};

export const updateTask = (id, taskData) => {
  return api.put(`/tasks/${id}`, taskData);
};

export const deleteTask = (id) => {
  return api.delete(`/tasks/${id}`);
};

// VULNERABILITY #1: No input sanitization before sending to backend
export const searchTasks = (searchTerm) => {
  // This will be vulnerable to SQL injection on the backend
  // return api.get(`/tasks/search?q=${searchTerm}`);

  //fix
  return api.get(`/tasks/search`, { params: { q: searchTerm } });
};

// User APIs
export const getCurrentUser = () => {
  return api.get('/users/me');
};

export const updateProfile = (userId, profileData) => {
  return api.put(`/users/${userId}/profile`, profileData);
};

// FIXED: Admin endpoint sekarang pakai Bearer token dari interceptor
// Backend yang validasi apakah user ini role admin atau bukan via JWT.
// Tidak perlu static API key — itu shared secret yang kalau bocor, siapapun jadi admin.
export const getAllUsers = () => {
  // return api.get('/admin/users', {
  //   headers: {
  //     // 'X-Admin-Key': ADMIN_API_KEY  // Hardcoded admin key!
  //     'X-Admin-Key': import.meta.env.VITE_ADMIN_API_KEY
  //   }
  // });
  return api.get('/admin/users');
};

export default api;
