import { apiClient } from '@/shared/api/client';
import { APIKey, APIKeyCreated } from './model';

export const apiKeyApi = {
  async list(): Promise<APIKey[]> {
    const { data } = await apiClient.get<APIKey[]>('/auth/api-keys');
    return data;
  },
  async create(name: string): Promise<APIKeyCreated> {
    const { data } = await apiClient.post<APIKeyCreated>('/auth/api-keys', { name });
    return data;
  },
  async remove(id: number): Promise<void> {
    await apiClient.delete(`/auth/api-keys/${id}`);
  },
};
