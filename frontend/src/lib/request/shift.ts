
export interface CreateShiftReq {
    date: string;
    start_time: string;
    end_time: string;
    role_id: number;
    location?: string;
}

export interface GetShiftListReq {
  showOnlyUnassigned: boolean; 
  limit: number;              
  offset: number;           
}

export interface GetShiftRequestListReq {
  limit: number;
  offset: number;
  status?: string;
}