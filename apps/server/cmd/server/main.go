package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fullstack-engineering-lab/server/internal/config"
	"fullstack-engineering-lab/server/internal/database"
	"fullstack-engineering-lab/server/internal/logger"
	"fullstack-engineering-lab/server/internal/router"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.Load("./configs")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Init logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	zap.L().Info("starting server",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port),
	)

	// Init database
	db, err := database.Init(&cfg.MySQL)
	if err != nil {
		zap.L().Fatal("failed to init database", zap.Error(err))
	}

	// Init Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		zap.L().Warn("redis not available, continuing without redis", zap.Error(err))
		rdb = nil
	}

	// Setup router
	r := router.Setup(cfg, db, rdb)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		zap.L().Info("server listening", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("server exited")
}
