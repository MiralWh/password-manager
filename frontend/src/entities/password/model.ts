export interface PasswordEntry {
  id: number;
  user_id: number;
  group_id?: number | null;
  title: string;
  username: string;
  password: string;
  url: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface PasswordInput {
  group_id?: number | null;
  title: string;
  username: string;
  password: string;
  url: string;
  notes: string;
}
