
import { PUBLIC_RMS_API_BASE_URL } from '$env/static/public';
import type { CreateShiftReq, GetShiftListReq, GetShiftRequestListReq } from '$lib/request/shift';
import type { GetShiftListResponse, GetShiftRequestListResponse } from '$lib/response/shift';

export async function fetchShifts(filter: GetShiftListReq): Promise<GetShiftListResponse> {
  const token = localStorage.getItem('token');
  const params = new URLSearchParams();
  if (filter.showOnlyUnassigned) params.append('show_only_unassigned', 'true');

  if (filter.limit <= 0) {
    filter.limit = 10
  }

  params.append('limit', filter.limit.toString());
  params.append('offset', filter.offset.toString());

  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift?${params.toString()}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });


  return res.json() as Promise<GetShiftListResponse>;
}

export async function createShift(shift: CreateShiftReq) {
  const token = localStorage.getItem('token');
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(shift),
  });
  return res.json();
}

export async function deleteShift(id: number) {
  const token = localStorage.getItem('token');
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift/${id}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return res.json();
}

export async function fetchShiftRequests(filter: GetShiftRequestListReq): Promise<GetShiftRequestListResponse> {
  const token = localStorage.getItem('token');
  const params = new URLSearchParams();
  if (filter.status) params.append('status', filter.status);
  params.append('limit', filter.limit?.toString() ?? '10');
  params.append('offset', filter.offset?.toString() ?? '0');

  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift/request?${params.toString()}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return res.json() as Promise<GetShiftRequestListResponse>;
}

export async function approveShiftRequest(id: number) {
  const token = localStorage.getItem('token');
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift/request/${id}/approve`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return res.json();
}

export async function rejectShiftRequest(id: number, reason: string) {
  const token = localStorage.getItem('token');
  const res = await fetch(`${PUBLIC_RMS_API_BASE_URL}/api/v1/shift/request/${id}/reject`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ reason }),
  });
  return res.json();
}