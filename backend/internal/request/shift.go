package request

import (
	"github.com/andibalo/payd-test/backend/internal/constants"
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/guregu/null/v6"
	"time"
)

type GetShiftAssignmentListReq struct {
	UserID int64 `json:"user_id" form:"user_id"`
	Limit  int   `json:"limit" form:"limit"`
	Offset int   `json:"offset" form:"offset"`

	UserEmail string `json:"-"`
}
type GetShiftListReq struct {
	ShowOnlyUnassigned bool `json:"show_only_unassigned" form:"show_only_unassigned"`
	Limit              int  `json:"limit" form:"limit"`
	Offset             int  `json:"offset" form:"offset"`

	UserEmail string `json:"-"`
}

type GetShiftRequestListReq struct {
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
	Status string `json:"status" form:"status"`
	UserID int64  `json:"user_id" form:"user_id"`

	UserEmail string `json:"-"`
}

type CreateShiftReq struct {
	Date      time.Time `json:"date" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	RoleID    int       `json:"role_id" binding:"required"`
	Location  string    `json:"location"`

	UserEmail string `json:"-"`
}

func (r *CreateShiftReq) ToModel() *model.Shift {

	shift := &model.Shift{
		Date:      r.Date,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		RoleID:    r.RoleID,
		Location:  null.StringFrom(r.Location),
		IsActive:  true,
		CreatedBy: r.UserEmail,
	}

	return shift
}

type UpdateShiftReq struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	RoleID    int       `json:"role_id"`
	Location  string    `json:"location"`
	IsActive  bool      `json:"is_active"`

	UserEmail string `json:"-"`
}

func (r *UpdateShiftReq) ToModel() *model.Shift {

	shift := &model.Shift{
		Date:      r.Date,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		RoleID:    r.RoleID,
		Location:  null.StringFrom(r.Location),
		IsActive:  r.IsActive,
		UpdatedBy: null.StringFrom(r.UserEmail),
	}

	return shift
}

type CreateShiftRequestReq struct {
	UserID  int64 `json:"user_id" binding:"required"`
	ShiftID int64 `json:"shift_id" binding:"required"`

	UserEmail string `json:"-"`
}

func (r *CreateShiftRequestReq) ToModel() *model.ShiftRequest {

	shiftReq := &model.ShiftRequest{
		UserID:      r.UserID,
		ShiftID:     r.ShiftID,
		Status:      constants.SHIFT_REQUEST_STATUS_PENDING,
		RequestedBy: r.UserEmail,
		CreatedBy:   r.UserEmail,
	}

	return shiftReq
}

type ApproveShiftRequestReq struct {
	RequestedShiftID int64 `json:"requested_shift_id" binding:"required"`

	UserEmail string `json:"-"`
}

type RejectShiftRequestReq struct {
	RequestedShiftID int64  `json:"requested_shift_id"`
	Reason           string `json:"reason" binding:"required"`

	UserEmail string `json:"-"`
}
