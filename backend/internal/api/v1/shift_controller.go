package v1

import (
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/middleware"
	"github.com/andibalo/payd-test/backend/internal/request"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/internal/service"
	"github.com/andibalo/payd-test/backend/pkg"
	"github.com/andibalo/payd-test/backend/pkg/apperr"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type ShiftController struct {
	cfg      config.Config
	shiftSvc service.ShiftService
}

func NewShiftController(cfg config.Config, shiftSvc service.ShiftService) *ShiftController {

	return &ShiftController{
		cfg:      cfg,
		shiftSvc: shiftSvc,
	}
}

func (h *ShiftController) AddRoutes(r *gin.Engine) {
	ar := r.Group("/api/v1/shift")

	ar.POST("", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.CreateShift)
	ar.GET("/:id", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.GetShiftByID)
	ar.GET("", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.GetShiftList)
	ar.PUT("/:id", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.UpdateShiftByID)
	ar.DELETE("/:id", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.DeleteShiftByID)
	ar.GET("/assignment", middleware.JwtMiddleware(h.cfg), h.GetShiftAssignmentsList)
	ar.GET("/request", middleware.JwtMiddleware(h.cfg), h.GetShiftRequestList)
	ar.POST("/request", middleware.JwtMiddleware(h.cfg), h.CreateShiftRequest)
	ar.PUT("/request/:id/approve", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.ApproveShiftRequest)
	ar.PUT("/request/:id/reject", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.RejectShiftRequest)
}

func (h *ShiftController) CreateShift(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ShiftController.CreateShift", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CreateShiftReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateShift] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserEmail = claims.Email

	err := h.shiftSvc.CreateShift(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateShift] Failed to create shift", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) GetShiftByID(c *gin.Context) {

	id, err := pkg.GetIntParam(c, "id")
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftByID] Error getting param", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	shift, err := h.shiftSvc.GetShiftByID(c.Request.Context(), id)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftByID] Failed to get shift", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, shift, nil)
	return
}

func (h *ShiftController) GetShiftList(c *gin.Context) {
	var req request.GetShiftListReq

	if err := c.ShouldBindQuery(&req); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftList] Failed to bind query parameters", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	shifts, err := h.shiftSvc.GetShiftList(c.Request.Context(), req)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftList] Failed to get shift list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, shifts, nil)
	return
}

func (h *ShiftController) UpdateShiftByID(c *gin.Context) {
	id, err := pkg.GetIntParam(c, "id")
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateShiftByID] Error getting param", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	var data request.UpdateShiftReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateShiftByID] Failed to bind JSON", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err = h.shiftSvc.UpdateShiftByID(c.Request.Context(), id, data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateShiftByID] Failed to update shift", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) DeleteShiftByID(c *gin.Context) {
	id, err := pkg.GetIntParam(c, "id")
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteShiftByID] Error getting param", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	err = h.shiftSvc.DeleteShiftByID(c.Request.Context(), id, claims.Email)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteShiftByID] Failed to delete shift", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) CreateShiftRequest(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ShiftController.CreateShiftRequest", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CreateShiftRequestReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateShiftRequestReq] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserEmail = claims.Email

	err := h.shiftSvc.CreateShiftRequest(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateShiftRequestReq] Failed to create shift request", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) ApproveShiftRequest(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ShiftController.ApproveShiftRequest", "controller")
	//defer endFunc()
	id, err := pkg.GetIntParam(c, "id")
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ApproveShiftRequest] Error getting param", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.ApproveShiftRequestReq

	data.UserEmail = claims.Email
	data.RequestedShiftID = id

	err = h.shiftSvc.ApproveShiftRequest(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ApproveShiftRequest] Failed to approve shift request", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) RejectShiftRequest(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ShiftController.RejectShiftRequest", "controller")
	//defer endFunc()
	id, err := pkg.GetIntParam(c, "id")
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RejectShiftRequest] Error getting param", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.RejectShiftRequestReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RejectShiftRequest] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserEmail = claims.Email
	data.RequestedShiftID = id

	err = h.shiftSvc.RejectShiftRequest(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RejectShiftRequest] Failed to reject shift request", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ShiftController) GetShiftRequestList(c *gin.Context) {
	var req request.GetShiftRequestListReq

	if err := c.ShouldBindQuery(&req); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftRequestList] Failed to bind query parameters", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	req.UserEmail = claims.Email

	shifts, err := h.shiftSvc.GetShiftRequestList(c.Request.Context(), req)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftRequestList] Failed to get shift request list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, shifts, nil)
	return
}

func (h *ShiftController) GetShiftAssignmentsList(c *gin.Context) {
	var req request.GetShiftAssignmentListReq

	if err := c.ShouldBindQuery(&req); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftAssignmentsList] Failed to bind query parameters", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	shifts, err := h.shiftSvc.GetShiftAssignmentList(c.Request.Context(), req)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetShiftAssignmentsList] Failed to get shift assignment list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, shifts, nil)
	return
}
