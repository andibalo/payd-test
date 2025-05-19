package v1

import (
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/request"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/internal/service"
	"github.com/andibalo/payd-test/backend/pkg/apperr"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type AuthController struct {
	cfg     config.Config
	authSvc service.AuthService
}

func NewAuthController(cfg config.Config, authSvc service.AuthService) *AuthController {

	return &AuthController{
		cfg:     cfg,
		authSvc: authSvc,
	}
}

func (h *AuthController) AddRoutes(r *gin.Engine) {
	ar := r.Group("/api/v1/auth")

	ar.POST("/register", h.Register)
	ar.POST("/login", h.Login)
}

func (h *AuthController) Register(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.RegisterUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Register] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	token, err := h.authSvc.Register(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Register] Failed to register user", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, token, nil)
	return
}

func (h *AuthController) Login(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.LoginUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Login] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Login] Failed to login", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, token, nil)
	return
}
