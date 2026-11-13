import { apiClient } from '@/shared/api/client';
import { PasswordEntry, PasswordInput } from './model';

export const passwordApi = {
  async list(search?: string): Promise<PasswordEntry[]> {
    const { data } = await apiClient.get<PasswordEntry[]>('/passwords', {
      params: search ? { q: search } : undefined,
    });
    return data;
  },
  async create(input: PasswordInput): Promise<PasswordEntry> {
    const { data } = await apiClient.post<PasswordEntry>('/passwords', input);
    return data;
  },
  async update(id: number, input: PasswordInput): Promise<PasswordEntry> {
    const { data } = await apiClient.put<PasswordEntry>(`/passwords/${id}`, input);
    return data;
  },
  async remove(id: number): Promise<void> {
    await apiClient.delete(`/passwords/${id}`);
  },
};
