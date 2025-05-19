package service

import (
	"context"
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/andibalo/payd-test/backend/internal/request"
	"github.com/andibalo/payd-test/backend/internal/response"
)

type UserService interface {
	GetWorkerList(ctx context.Context, req request.GetWorkerListReq) ([]model.User, error)
}

type AuthService interface {
	Register(ctx context.Context, req request.RegisterUserReq) (token string, err error)
	Login(ctx context.Context, req request.LoginUserReq) (token string, err error)
}

type ShiftService interface {
	CreateShift(ctx context.Context, req request.CreateShiftReq) error
	GetShiftByID(ctx context.Context, id int64) (*model.Shift, error)
	GetShiftList(ctx context.Context, req request.GetShiftListReq) (resp response.GetShiftListResponse, err error)
	UpdateShiftByID(ctx context.Context, id int64, req request.UpdateShiftReq) error
	DeleteShiftByID(ctx context.Context, id int64, deletedBy string) error
	CreateShiftRequest(ctx context.Context, req request.CreateShiftRequestReq) error
	ApproveShiftRequest(ctx context.Context, req request.ApproveShiftRequestReq) error
	RejectShiftRequest(ctx context.Context, req request.RejectShiftRequestReq) error
	GetShiftRequestList(ctx context.Context, req request.GetShiftRequestListReq) (resp response.GetShiftRequestListResponse, err error)
	GetShiftAssignmentList(ctx context.Context, req request.GetShiftAssignmentListReq) (resp response.GetShiftAssignmentListResponse, err error)
}
