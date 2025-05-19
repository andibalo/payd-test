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
	"github.com/andibalo/payd-test/backend/pkg"
	"github.com/andibalo/payd-test/backend/pkg/apperr"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type authService struct {
	cfg      config.Config
	userRepo repository.UserRepository
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository) AuthService {

	return &authService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req request.RegisterUserReq) (token string, err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.Register", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to get user by email", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Invalid Email/Password")
	}

	if existingUser != nil && existingUser.ID != 0 {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] User already exists")
		return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User already exists")
	}

	user, err := s.mapCreateUserReqToUserModel(ctx, req)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to map payload to user model", zap.Error(err))
		return "", oops.Wrapf(err, "[Register] Failed to map payload to user model")
	}

	err = s.userRepo.Save(user)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to insert user to database", zap.Error(err))

		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	token, err = pkg.GenerateToken(existingUser)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to generate JWT Token for user", zap.String("email", req.Email))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return token, nil
}

func (s *authService) mapCreateUserReqToUserModel(ctx context.Context, data request.RegisterUserReq) (*model.User, error) {

	hasedPassword, err := pkg.HashPassword(data.Password)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[mapCreateUserReqToUserModel] Failed to hash password", zap.Error(err))

		return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return &model.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Role:      constants.WORKER_ROLE,
		Password:  hasedPassword,
		CreatedBy: data.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req request.LoginUserReq) (token string, err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.Login", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[Login] Invalid email/password", zap.Error(err))
			return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Invalid Email/Password")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Failed to get user by email", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	isMatch := pkg.CheckPasswordHash(req.Password, existingUser.Password)
	if !isMatch {
		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Invalid password for user", zap.String("email", req.Email))
		return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Invalid Email/Password")
	}

	token, err = pkg.GenerateToken(existingUser)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Failed to generate JWT Token for user", zap.String("email", req.Email))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return token, nil
}
