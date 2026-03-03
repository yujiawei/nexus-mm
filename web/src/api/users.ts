import client from './client';
import type { LoginRequest, RegisterRequest, LoginResponse, User } from './types';

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await client.post<LoginResponse>('/users/login', data);
  return res.data;
}

export async function register(data: RegisterRequest): Promise<LoginResponse> {
  const res = await client.post<LoginResponse>('/users/register', data);
  return res.data;
}

export async function getMe(): Promise<User> {
  const res = await client.get<User>('/users/me');
  return res.data;
}
