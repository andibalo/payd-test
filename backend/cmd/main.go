package main

import (
	"context"
	"errors"
	"fmt"
	rms "github.com/andibalo/payd-test/backend"
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/pkg/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.InitConfig()

	database := db.InitDB(cfg)

	server := rms.NewServer(cfg, database)

	cfg.Logger().Info(fmt.Sprintf("Server starting at port %s", cfg.AppAddress()))

	go func() {
		if err := server.Start(cfg.AppAddress()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cfg.Logger().Fatal("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cfg.Logger().Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger().Fatal("Server force to shutdown")
	}

	_ = database.Close()

	cfg.Logger().Info("Server exiting")
}
