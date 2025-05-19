
import type { PaginationMeta } from './common';
import type { Shift } from '../model/shift';

export interface ShiftListData {
    shifts: Shift[];
    meta: PaginationMeta;
}

export interface GetShiftListResponse {
    data: ShiftListData;
    success: string;
}

export interface ShiftRequest {
  id: number;
  user_id: number;
  shift_id: number;
  shift_date: string;
  shift_start_time: string;
  shift_end_time: string;
  shift_role_id: number;
  shift_role_name: string;
  status: string;
  requested_by: string;
  admin_actor: string;
  rejection_reason: string | null;
  created_at: string;
  created_by: string;
  updated_at: string | null;
  updated_by: string | null;
  deleted_at: string | null;
  deleted_by: string | null;
}

export interface ShiftRequestListData {
  request_shifts: ShiftRequest[];
  meta: PaginationMeta;
}

export interface GetShiftRequestListResponse {
  data: ShiftRequestListData;
  success: string;
}