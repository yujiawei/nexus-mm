import { create } from 'zustand';
import type { User } from '../api/types';
import * as usersApi from '../api/users';

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string, nickname: string) => Promise<void>;
  logout: () => void;
  loadUser: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem('token'),
  loading: false,
  error: null,

  login: async (username, password) => {
    set({ loading: true, error: null });
    try {
      const res = await usersApi.login({ username, password });
      localStorage.setItem('token', res.token);
      set({ user: res.user, token: res.token, loading: false });
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Login failed';
      set({ error: message, loading: false });
      throw err;
    }
  },

  register: async (username, email, password, nickname) => {
    set({ loading: true, error: null });
    try {
      const res = await usersApi.register({ username, email, password, nickname });
      localStorage.setItem('token', res.token);
      set({ user: res.user, token: res.token, loading: false });
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Registration failed';
      set({ error: message, loading: false });
      throw err;
    }
  },

  logout: () => {
    localStorage.removeItem('token');
    set({ user: null, token: null });
  },

  loadUser: async () => {
    const token = localStorage.getItem('token');
    if (!token) return;
    set({ loading: true });
    try {
      const user = await usersApi.getMe();
      set({ user, token, loading: false });
    } catch {
      localStorage.removeItem('token');
      set({ user: null, token: null, loading: false });
    }
  },

  clearError: () => set({ error: null }),
}));
