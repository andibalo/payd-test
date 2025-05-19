package rms

import (
	"context"
	"database/sql"
	"github.com/andibalo/payd-test/backend/internal/api"
	v1 "github.com/andibalo/payd-test/backend/internal/api/v1"
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/middleware"
	"github.com/andibalo/payd-test/backend/internal/repository"
	"github.com/andibalo/payd-test/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	gin *gin.Engine
	srv *http.Server
}

func NewServer(cfg config.Config, db *sql.DB) *Server {

	router := gin.New()

	router.Use(middleware.LogPreReq(cfg.Logger()))

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")

	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(db)
	shiftRepo := repository.NewShiftRepository(db)

	userSvc := service.NewUserService(cfg, userRepo)
	authSvc := service.NewAuthService(cfg, userRepo)
	shiftSvc := service.NewShiftService(cfg, shiftRepo)

	uc := v1.NewUserController(cfg, userSvc)
	ac := v1.NewAuthController(cfg, authSvc)
	sc := v1.NewShiftController(cfg, shiftSvc)

	registerHandlers(router, &api.HealthCheck{}, uc, ac, sc)

	return &Server{
		gin: router,
	}
}

func (s *Server) Start(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.gin,
	}

	s.srv = srv

	return srv.ListenAndServe()
}

func (s *Server) GetGin() *gin.Engine {

	return s.gin
}

func (s *Server) Shutdown(ctx context.Context) error {

	return s.srv.Shutdown(ctx)
}

func registerHandlers(g *gin.Engine, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(g)
	}
}
