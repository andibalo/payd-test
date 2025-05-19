package v1

import (
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/middleware"
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

type UserController struct {
	cfg     config.Config
	userSvc service.UserService
}

func NewUserController(cfg config.Config, userSvc service.UserService) *UserController {

	return &UserController{
		cfg:     cfg,
		userSvc: userSvc,
	}
}

func (h *UserController) AddRoutes(r *gin.Engine) {
	ur := r.Group("/api/v1/user")

	ur.GET("/worker", middleware.JwtMiddleware(h.cfg), h.GetWorkerList)
}

func (h *UserController) GetWorkerList(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.GetWorkersList", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetWorkerListReq

	data.UserEmail = claims.Email
	user, err := h.userSvc.GetWorkerList(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetWorkerList] Failed to get worker list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, user, nil)
	return
}
