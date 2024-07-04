package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"time-tracker/internal/config"
	"time-tracker/internal/connectors"
	taskcontroller "time-tracker/internal/controllers/v1/task"
	usercontroller "time-tracker/internal/controllers/v1/user"
	"time-tracker/internal/db/postgres"
	_ "time-tracker/internal/docs/v1"
	"time-tracker/internal/domain"
	"time-tracker/internal/storages"
	"time-tracker/internal/utils/slogutils"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title           Time Tracker
// @version         1.0
// @description     API to track performance

// @BasePath  /api/v1
func Run(cfg config.Config, logger *slog.Logger) {
	slog.SetDefault(logger)

	slog.Info("Setting up server dependencies")

	db, err := postgres.NewDatabase(cfg.DB)
	if err != nil {
		slog.Error("initialize db", slogutils.ErrorAttr(err))
		return
	}

	userStorage := storages.NewUserStorage(db)
	tasksStorage := storages.NewTaskStorage(db)

	userService := domain.NewUserService(
		userStorage,
		connectors.NewUserInfoConnector())
	tasksService := domain.NewTaskService(tasksStorage)

	userController := usercontroller.NewUserController(userService)
	tasksController := taskcontroller.NewTasksController(tasksService)

	switch cfg.Env {
	case config.EnvLocal:
		gin.SetMode(gin.DebugMode)
	case config.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(sloggin.New(logger))
	engine.Use(gin.Recovery())
	engine.GET("api/v1/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	userController.RegisterRoutes(engine)
	tasksController.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:    cfg.HTTPServer.IpAddress + ":" + cfg.HTTPServer.Port,
		Handler: engine.Handler(),
	}

	slog.Info("Starting server ...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen: %s\n", slogutils.ErrorAttr(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", slogutils.ErrorAttr(err))
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")
}
