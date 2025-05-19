package model

import (
	"github.com/guregu/null/v6"
	"time"
)

type User struct {
	ID        int64       `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Role      string      `json:"role"`
	CreatedBy string      `json:"created_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedBy null.String `json:"updated_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedBy null.String `json:"-"`
	DeletedAt null.Time   `json:"-"`
}
