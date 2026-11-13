export interface Group {
  id: number;
  user_id: number;
  name: string;
  password_count: number;
  created_at: string;
}

export interface GroupInput {
  name: string;
}
