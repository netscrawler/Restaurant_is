export interface User {
  id: string;
  email: string;
  phone: string;
  name: string;
  roles: string[];
  created_at: string;
  updated_at: string;
}

export interface Staff {
  id: string;
  work_email: string;
  position: string;
  roles: string[];
  created_at: string;
  updated_at: string;
}

export interface Dish {
  id: string;
  name: string;
  description: string;
  price: number;
  category_id: number;
  category_name: string;
  image_url: string;
  available: boolean;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  dish_id: string;
  dish_name: string;
  quantity: number;
  price: number;
}

export interface Order {
  id: string;
  user_id: string;
  items: OrderItem[];
  total_amount: number;
  status: OrderStatus;
  delivery_address: string;
  comment: string;
  created_at: string;
  updated_at: string;
}

export type OrderStatus = 
  | 'ORDER_STATUS_CREATED'
  | 'ORDER_STATUS_CONFIRMED'
  | 'ORDER_STATUS_COOKING'
  | 'ORDER_STATUS_READY'
  | 'ORDER_STATUS_DELIVERED'
  | 'ORDER_STATUS_CANCELLED';

export interface AuthResponse {
  access_token: string;
  expires_in: number;
  refresh_token: string;
  refresh_token_expires_in: number;
  user: User;
}

export interface LoginInitRequest {
  phone: string;
}

export interface LoginConfirmRequest {
  phone: string;
  code: string;
}

export interface StaffLoginRequest {
  staff: Staff;
  password: string;
}

export interface CreateOrderRequest {
  items: Array<{
    dish_id: string;
    quantity: number;
  }>;
  delivery_address?: string;
  comment?: string;
}

export interface PaginatedResponse<T> {
  [key: string]: T[] | number;
  total: number;
  page: number;
  page_size: number;
}

export interface ApiError {
  code: number;
  message: string;
  details: string[];
} 