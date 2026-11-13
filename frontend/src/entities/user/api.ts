import { apiClient } from '@/shared/api/client';
import { LoginResponse } from './model';

export const authApi = {
  async login(username: string, password: string): Promise<LoginResponse> {
    const { data } = await apiClient.post<LoginResponse>('/auth/login', { username, password });
    return data;
  },
};
