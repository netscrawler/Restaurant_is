import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { 
  User, 
  Dish, 
  Order, 
  AuthResponse, 
  LoginInitRequest, 
  LoginConfirmRequest,
  CreateOrderRequest,
  PaginatedResponse
} from '../types/api';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: 'http://localhost:8080/api/v1',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Добавляем интерцептор для автоматического добавления токена
    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('access_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Добавляем интерцептор для обработки ошибок
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async loginInit(data: LoginInitRequest): Promise<{ status: string; error?: string }> {
    const response: AxiosResponse = await this.api.post('/auth/client/login/init', data);
    return response.data;
  }

  async loginConfirm(data: LoginConfirmRequest): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/client/login/confirm', data);
    return response.data;
  }

  async validateToken(token: string): Promise<{ valid: boolean; user: User }> {
    const response: AxiosResponse = await this.api.post('/auth/validate', { token });
    return response.data;
  }

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
  }

  // Menu endpoints
  async getDishes(params?: {
    category_id?: number;
    only_available?: boolean;
    page?: number;
    page_size?: number;
  }): Promise<PaginatedResponse<Dish>> {
    const response: AxiosResponse<PaginatedResponse<Dish>> = await this.api.get('/menu/dishes', { params });
    return response.data;
  }

  async getDish(id: string): Promise<Dish> {
    const response: AxiosResponse<Dish> = await this.api.get(`/menu/dishes/${id}`);
    return response.data;
  }

  // User endpoints
  async getProfile(): Promise<User> {
    const response: AxiosResponse<User> = await this.api.get('/users/profile');
    return response.data;
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    const response: AxiosResponse<User> = await this.api.put('/users/profile', data);
    return response.data;
  }

  // Order endpoints
  async getOrders(params?: {
    status?: string;
    page?: number;
    page_size?: number;
  }): Promise<PaginatedResponse<Order>> {
    const response: AxiosResponse<PaginatedResponse<Order>> = await this.api.get('/orders', { params });
    return response.data;
  }

  async getOrder(id: string): Promise<Order> {
    const response: AxiosResponse<Order> = await this.api.get(`/orders/${id}`);
    return response.data;
  }

  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const response: AxiosResponse<Order> = await this.api.post('/orders', data);
    return response.data;
  }
}

export const apiService = new ApiService(); 