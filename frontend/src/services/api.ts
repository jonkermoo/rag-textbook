import axios from 'axios';
import type { LoginRequest, LoginResponse, Textbook, 
              TextbookStatus, QueryRequest, QueryResponse } from '../types';

// Base URL for Go backend
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

// Axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Auto add JWT token to request if it exists
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth API calls
export const authAPI = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', credentials);
    return response.data;
  },

  register: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/register', credentials);
    return response.data;
  },
};

// Textbook API calls
export const textbookAPI = {
  list: async (): Promise<Textbook[]> => {
    const response = await api.get<Textbook[]>('/textbooks');
    return response.data;
  },

  get: async (id: number): Promise<Textbook> => {
    const response = await api.get<Textbook>(`/textbooks/${id}`);
    return response.data;
  },

  getStatus: async (id: number): Promise<TextbookStatus> => {
    const response = await api.get<TextbookStatus>(`/textbooks/${id}/status`);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/textbooks/${id}`);
  },

  upload: async (file: File, title: string): Promise<Textbook> => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('title', title);

    const response = await api.post<Textbook>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

// Query API calls
export const queryAPI = {
  ask: async (request: QueryRequest): Promise<QueryResponse> => {
    const response = await api.post<QueryResponse>('/query', request);
    return response.data;
  },
};