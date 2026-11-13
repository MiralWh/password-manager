export interface APIKey {
  id: number;
  user_id: number;
  name: string;
  key_prefix: string;
  created_at: string;
  last_used_at?: string | null;
}

export interface APIKeyCreated extends APIKey {
  key: string;
}
