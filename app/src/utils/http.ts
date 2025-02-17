import axios from 'axios';

const http = axios.create({
  baseURL: '/',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor
http.interceptors.request.use(
  (config) => {
    // You can add auth headers or other request modifications here
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor
http.interceptors.response.use(
  (response) => {
    // You can transform response data here if needed
    return response;
  },
  (error) => {
    // Handle response errors here
    return Promise.reject(error);
  }
);

export { http };
