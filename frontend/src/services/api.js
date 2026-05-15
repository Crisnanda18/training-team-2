import axios from 'axios';
import { API_BASE_URL, ADMIN_API_KEY } from '../config';
import { getToken } from '../utils/storage';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// VULNERABILITY: Logging sensitive data
api.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // VULNERABILITY: Logging requests with sensitive data
    if (console && console.log) {
      console.log('API Request:', config.method, config.url, config.data);
    }
    
    return config;
  },
  (error) => {
    console.error('Request Error:', error);
    return Promise.reject(error);
  }
);

// VULNERABILITY: Logging sensitive response data
api.interceptors.response.use(
  (response) => {
    console.log('API Response:', response.data);
    return response;
  },
  (error) => {
    console.error('Response Error:', error.response?.data);
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
  return api.get(`/tasks/search?q=${searchTerm}`);
};

// User APIs
export const getCurrentUser = () => {
  return api.get('/users/me');
};

export const updateProfile = (userId, profileData) => {
  return api.put(`/users/${userId}/profile`, profileData);
};

// VULNERABILITY #2: Admin endpoint accessible without proper authorization check
// VULNERABILITY #4: Hardcoded API key sent in request
export const getAllUsers = () => {
  return api.get('/admin/users', {
    headers: {
      'X-Admin-Key': ADMIN_API_KEY  // Hardcoded admin key!
    }
  });
};

export default api;
