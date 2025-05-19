package response

import (
	"github.com/guregu/null/v6"
	"time"
)

type GetShiftRequestListData struct {
	ID              int64       `json:"id"`
	UserID          int64       `json:"user_id"`
	ShiftID         int64       `json:"shift_id"`
	ShiftDate       time.Time   `json:"shift_date"`
	ShiftStartTime  time.Time   `json:"shift_start_time"`
	ShiftEndTime    time.Time   `json:"shift_end_time"`
	ShiftRoleID     int64       `json:"shift_role_id"`
	ShiftRoleName   string      `json:"shift_role_name"`
	Status          string      `json:"status"`
	RequestedBy     string      `json:"requested_by"`
	AdminActor      null.String `json:"admin_actor"`
	RejectionReason null.String `json:"rejection_reason"`
	CreatedAt       time.Time   `json:"created_at"`
	CreatedBy       string      `json:"created_by"`
	UpdatedAt       null.Time   `json:"updated_at"`
	UpdatedBy       null.String `json:"updated_by"`
	DeletedAt       null.Time   `json:"deleted_at"`
	DeletedBy       null.String `json:"deleted_by"`
}

type GetShiftListData struct {
	ID        int         `json:"id"`
	Date      time.Time   `json:"date"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	RoleID    int         `json:"role_id"`
	RoleName  string      `json:"role_name"`
	Location  null.String `json:"location"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	DeletedAt null.Time   `json:"deleted_at"`
	DeletedBy null.String `json:"deleted_by"`
}

type GetShiftRequestListResponse struct {
	Data []GetShiftRequestListData `json:"request_shifts"`
	Meta PaginationMeta            `json:"meta"`
}

type GetShiftListResponse struct {
	Data []GetShiftListData `json:"shifts"`
	Meta PaginationMeta     `json:"meta"`
}

type GetShiftAssignmentListData struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	Email          string      `json:"email"`
	ShiftID        int64       `json:"shift_id"`
	ShiftDate      time.Time   `json:"shift_date"`
	ShiftStartTime time.Time   `json:"shift_start_time"`
	ShiftEndTime   time.Time   `json:"shift_end_time"`
	ShiftRoleName  string      `json:"shift_role_name"`
	AssignedAt     time.Time   `json:"assigned_at"`
	AssignedBy     string      `json:"assigned_by_by"`
	CreatedAt      time.Time   `json:"created_at"`
	CreatedBy      string      `json:"created_by"`
	UpdatedAt      null.Time   `json:"updated_at"`
	UpdatedBy      null.String `json:"updated_by"`
	DeletedAt      null.Time   `json:"deleted_at"`
	DeletedBy      null.String `json:"deleted_by"`
}

type GetShiftAssignmentListResponse struct {
	Data []GetShiftAssignmentListData `json:"shifts"`
	Meta PaginationMeta               `json:"meta"`
}
