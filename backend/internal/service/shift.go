package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/constants"
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/andibalo/payd-test/backend/internal/repository"
	"github.com/andibalo/payd-test/backend/internal/request"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type shiftService struct {
	cfg       config.Config
	shiftRepo repository.ShiftRepository
}

func NewShiftService(cfg config.Config, shiftRepo repository.ShiftRepository) ShiftService {

	return &shiftService{
		cfg:       cfg,
		shiftRepo: shiftRepo,
	}
}

func (s *shiftService) CreateShift(ctx context.Context, req request.CreateShiftReq) error {

	err := s.shiftRepo.Save(req.ToModel())
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShift] Failed to create shift", zap.Error(err))
		return err
	}

	return nil
}

func (s *shiftService) GetShiftByID(ctx context.Context, id int64) (*model.Shift, error) {
	shift, err := s.shiftRepo.GetByID(id)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetShiftByID] Failed to get shift by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return shift, nil
}

func (s *shiftService) GetShiftList(ctx context.Context, req request.GetShiftListReq) (resp response.GetShiftListResponse, err error) {

	filter := repository.GetShiftListFilter{
		ShowOnlyUnassigned: req.ShowOnlyUnassigned,
		Limit:              req.Limit,
		Offset:             req.Offset,
	}

	shifts, pagination, err := s.shiftRepo.GetList(filter)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetShiftList] Failed to get shifts list", zap.Any("filter", filter), zap.Error(err))
		return resp, err
	}

	resp = response.GetShiftListResponse{
		Data: shifts,
		Meta: response.PaginationMeta{
			CurrentPage: pagination.CurrentPage,
			TotalPages:  pagination.TotalPages,
			TotalItems:  pagination.TotalElements,
		},
	}

	return resp, nil
}

func (s *shiftService) UpdateShiftByID(ctx context.Context, id int64, req request.UpdateShiftReq) error {

	_, err := s.shiftRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateShiftByID] Shift not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Shift not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateShiftByID] Failed to get shift by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get shift by id")
	}

	err = s.shiftRepo.UpdateByID(id, req.ToModel())
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateShiftByID] Failed to update shift", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *shiftService) DeleteShiftByID(ctx context.Context, id int64, deletedBy string) error {
	_, err := s.shiftRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[DeleteShiftByID] Shift not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Shift not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteShiftByID] Failed to get shift by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get shift by id")
	}

	err = s.shiftRepo.DeleteByID(id, deletedBy)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteShiftByID] Failed to delete shift", zap.Int64("id", id), zap.String("deletedBy", deletedBy), zap.Error(err))
		return err
	}
	return nil
}

func (s *shiftService) GetShiftRequestList(ctx context.Context, req request.GetShiftRequestListReq) (resp response.GetShiftRequestListResponse, err error) {

	filter := repository.GetShiftRequestListFilter{
		Limit:  req.Limit,
		Status: req.Status,
		UserID: req.UserID,
		Offset: req.Offset,
	}

	shiftRequests, pagination, err := s.shiftRepo.GetShiftRequestList(filter)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetShiftRequestList] Failed to get shift request list", zap.Any("filter", filter), zap.Error(err))
		return resp, err
	}

	resp = response.GetShiftRequestListResponse{
		Data: shiftRequests,
		Meta: response.PaginationMeta{
			CurrentPage: pagination.CurrentPage,
			TotalPages:  pagination.TotalPages,
			TotalItems:  pagination.TotalElements,
		},
	}

	return resp, nil
}

func (s *shiftService) CreateShiftRequest(ctx context.Context, req request.CreateShiftRequestReq) error {

	shiftDetail, err := s.shiftRepo.GetByID(req.ShiftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Shift not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Shift not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to get shift by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get shift by id")
	}

	isShiftAlreadyAssigned, err := s.shiftRepo.CheckIfShiftIsAlreadyAssigned(req.ShiftID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to check if shift already assigned", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to check if shift already assigned")
	}

	if isShiftAlreadyAssigned {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Shift is already assigned", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Shift is already assigned")
	}

	isShiftRequestTimeOverlaps, err := s.shiftRepo.CheckIfShiftRequestTimeOverlaps(req.UserID, shiftDetail.Date, shiftDetail.StartTime, shiftDetail.EndTime)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to check shift request overlaps", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to check shift request overlaps")
	}

	if isShiftRequestTimeOverlaps {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Shift request time overlaps", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Shift request time overlaps")
	}

	isUserHasOtherShiftAssignment, err := s.shiftRepo.CheckUserAssignedShiftExistsByDate(req.UserID, shiftDetail.Date)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to check assigned shift exists", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to check assigned shift exists")
	}

	if isUserHasOtherShiftAssignment {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Already have an assigned shift on this day", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Already have an assigned shift on this day")
	}

	userWeeklyShiftCount, err := s.shiftRepo.GetUserWeeklyAssignedShiftCountByDate(req.UserID, shiftDetail.Date)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to check user weekly shift count", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to check user weekly shift count")
	}

	if userWeeklyShiftCount == constants.MAX_ASSIGNED_SHIFT_PER_WEEK {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] User already reached shift assignment limit this week", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User already reached shift assignment limit this week")
	}

	err = s.shiftRepo.SaveShiftRequest(req.ToModel())
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateShiftRequest] Failed to create shift request", zap.String("requestedBy", req.UserEmail), zap.Error(err))
		return err
	}
	return nil
}

func (s *shiftService) ApproveShiftRequest(ctx context.Context, req request.ApproveShiftRequestReq) error {

	shiftRequest, err := s.shiftRepo.GetShiftRequestByID(req.RequestedShiftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[ApproveShiftRequest] Shift request not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Shift request not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[ApproveShiftRequest] Failed to fetch shift request by ID", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get shift request by id")
	}

	if shiftRequest.Status != constants.SHIFT_REQUEST_STATUS_PENDING {
		s.cfg.Logger().ErrorWithContext(ctx, "[ApproveShiftRequest] Shift request status is not pending", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Shift request status is not pending")
	}

	shiftRequest.Status = constants.SHIFT_REQUEST_STATUS_APPROVED
	shiftRequest.AdminActor = null.StringFrom(req.UserEmail)
	shiftRequest.UpdatedBy = null.StringFrom(req.UserEmail)

	err = s.shiftRepo.UpdateShiftRequestByID(req.RequestedShiftID, shiftRequest)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ApproveShiftRequest] Failed to update shift request", zap.Error(err))
		return err
	}

	workerShift := &model.WorkerShift{
		UserID:     shiftRequest.UserID,
		ShiftID:    shiftRequest.ShiftID,
		AssignedBy: req.UserEmail,
		CreatedBy:  req.UserEmail,
	}

	err = s.shiftRepo.SaveWorkerShift(workerShift)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ApproveShiftRequest] Failed to create worker shift", zap.Error(err))
		return err
	}

	return nil
}

func (s *shiftService) RejectShiftRequest(ctx context.Context, req request.RejectShiftRequestReq) error {

	shiftRequest, err := s.shiftRepo.GetShiftRequestByID(req.RequestedShiftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[RejectShiftRequest] Shift request not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Shift request not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[RejectShiftRequest] Failed to fetch shift request by ID", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get shift request by id")
	}

	if shiftRequest.Status != constants.SHIFT_REQUEST_STATUS_PENDING {
		s.cfg.Logger().ErrorWithContext(ctx, "[RejectShiftRequest] Shift request status is not pending", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Shift request status is not pending")
	}

	shiftRequest.Status = constants.SHIFT_REQUEST_STATUS_REJECTED
	shiftRequest.AdminActor = null.StringFrom(req.UserEmail)
	shiftRequest.RejectionReason = null.StringFrom(req.Reason)
	shiftRequest.UpdatedBy = null.StringFrom(req.UserEmail)

	err = s.shiftRepo.UpdateShiftRequestByID(req.RequestedShiftID, shiftRequest)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[RejectShiftRequest] Failed to update shift request", zap.Error(err))
		return err
	}

	return nil
}

func (s *shiftService) GetShiftAssignmentList(ctx context.Context, req request.GetShiftAssignmentListReq) (resp response.GetShiftAssignmentListResponse, err error) {

	filter := repository.GetShiftAssignmentListFilter{
		UserID: req.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	shiftAssignments, pagination, err := s.shiftRepo.GetShiftAssignmentList(filter)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetShiftAssignmentList] Failed to get shifts assignment list", zap.Any("filter", filter), zap.Error(err))
		return resp, err
	}

	resp = response.GetShiftAssignmentListResponse{
		Data: shiftAssignments,
		Meta: response.PaginationMeta{
			CurrentPage: pagination.CurrentPage,
			TotalPages:  pagination.TotalPages,
			TotalItems:  pagination.TotalElements,
		},
	}

	return resp, nil
}
