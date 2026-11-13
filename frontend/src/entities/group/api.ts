import { apiClient } from '@/shared/api/client';
import { Group, GroupInput } from './model';

export const groupApi = {
  async list(search?: string): Promise<Group[]> {
    const { data } = await apiClient.get<Group[]>('/groups', {
      params: search ? { q: search } : undefined,
    });
    return data;
  },
  async create(input: GroupInput): Promise<Group> {
    const { data } = await apiClient.post<Group>('/groups', input);
    return data;
  },
  async rename(id: number, input: GroupInput): Promise<Group> {
    const { data } = await apiClient.put<Group>(`/groups/${id}`, input);
    return data;
  },
  async remove(id: number): Promise<void> {
    await apiClient.delete(`/groups/${id}`);
  },
};
