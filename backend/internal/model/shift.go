package model

import (
	"github.com/guregu/null/v6"
	"time"
)

type Shift struct {
	ID        int         `json:"id"`
	Date      time.Time   `json:"date"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	RoleID    int         `json:"role_id"`
	Location  null.String `json:"location"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	DeletedAt null.Time   `json:"deleted_at"`
	DeletedBy null.String `json:"deleted_by"`
}

type ShiftRequest struct {
	ID              int64       `json:"id"`
	UserID          int64       `json:"user_id"`
	ShiftID         int64       `json:"shift_id"`
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

type WorkerShift struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	ShiftID    int64       `json:"shift_id"`
	AssignedAt time.Time   `json:"assigned_at"`
	AssignedBy string      `json:"assigned_by_by"`
	CreatedAt  time.Time   `json:"created_at"`
	CreatedBy  string      `json:"created_by"`
	UpdatedAt  null.Time   `json:"updated_at"`
	UpdatedBy  null.String `json:"updated_by"`
	DeletedAt  null.Time   `json:"deleted_at"`
	DeletedBy  null.String `json:"deleted_by"`
}
