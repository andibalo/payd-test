package repository

import (
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"time"
)

type UserRepository interface {
	Save(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetList(filter GetUserListFilter) ([]model.User, error)
}

type ShiftRepository interface {
	Save(shift *model.Shift) error
	GetByID(id int64) (*model.Shift, error)
	GetList(filter GetShiftListFilter) ([]response.GetShiftListData, *httpresp.Pagination, error)
	UpdateByID(id int64, shift *model.Shift) error
	DeleteByID(id int64, deletedBy string) error
	GetShiftRequestByID(id int64) (*model.ShiftRequest, error)
	GetShiftRequestList(filter GetShiftRequestListFilter) ([]response.GetShiftRequestListData, *httpresp.Pagination, error)
	SaveShiftRequest(shiftRequest *model.ShiftRequest) error
	SaveWorkerShift(workerShift *model.WorkerShift) error
	UpdateShiftRequestByID(id int64, sr *model.ShiftRequest) error
	CheckUserAssignedShiftExistsByDate(userID int64, shiftDate time.Time) (bool, error)
	CheckIfShiftIsAlreadyAssigned(shiftID int64) (bool, error)
	GetUserWeeklyAssignedShiftCountByDate(userID int64, shiftDate time.Time) (int, error)
	CheckIfShiftRequestTimeOverlaps(userID int64, shiftDate, requestedStartTime, requestedEndTime time.Time) (bool, error)
	GetShiftAssignmentList(filter GetShiftAssignmentListFilter) ([]response.GetShiftAssignmentListData, *httpresp.Pagination, error)
}
