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
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type userService struct {
	cfg      config.Config
	userRepo repository.UserRepository
}

func NewUserService(cfg config.Config, userRepo repository.UserRepository) UserService {

	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *userService) GetWorkerList(ctx context.Context, req request.GetWorkerListReq) ([]model.User, error) {
	//ctx, endFunc := trace.Start(ctx, "UserService.GetWorkerList", "service")
	//defer endFunc()

	workers, err := s.userRepo.GetList(repository.GetUserListFilter{Role: constants.WORKER_ROLE})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetWorkerList] Failed to get worker list", zap.Error(err))
		return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user devices")
	}

	return workers, nil
}
