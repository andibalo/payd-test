export interface Shift {
  id: number;
  date: string;
  start_time: string;
  end_time: string;
  role_id: number;
  role_name: string;
  location: string;
  is_active: boolean;
  created_at: string;
  created_by: string;
  updated_at: string | null;
  updated_by: string | null;
  deleted_at: string | null;
  deleted_by: string | null;
}