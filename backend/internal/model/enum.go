package model

import (
	"github.com/guregu/null/v6"
	"time"
)

type ShiftRoleEnum struct {
	ID        int64       `json:"id"`
	RoleName  string      `json:"role_name"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	DeletedAt null.Time   `json:"deleted_at"`
	DeletedBy null.String `json:"deleted_by"`
}
