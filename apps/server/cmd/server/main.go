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
	mqttPkg "fullstack-engineering-lab/server/pkg/mqtt"
	tcpPkg "fullstack-engineering-lab/server/pkg/tcp"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("./configs")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 初始化日志
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	zap.L().Info("starting server",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port),
	)

	// 初始化数据库
	db, err := database.Init(&cfg.MySQL)
	if err != nil {
		zap.L().Fatal("failed to init database", zap.Error(err))
	}

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试 Redis 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		zap.L().Warn("redis not available, continuing without redis", zap.Error(err))
		rdb = nil
	}

	// 启动 TCP Server
	tcpCfg := &tcpPkg.ServerConfig{
		Host:            cfg.TCP.Host,
		Port:            cfg.TCP.Port,
		ReadTimeoutSec:  cfg.TCP.ReadTimeoutSec,
		WriteTimeoutSec: cfg.TCP.WriteTimeoutSec,
	}
	tcpServer := tcpPkg.NewServer(tcpCfg)
	if err := tcpServer.Start(context.Background()); err != nil {
		zap.L().Fatal("TCP 服务器启动失败", zap.Error(err))
	}
	defer tcpServer.Stop()
	zap.L().Info("TCP 服务器已启动", zap.String("addr", tcpServer.Addr()))

	// 初始化 MQTT Client（失败只 Warn，不 Fatal）
	var mqttClient *mqttPkg.MQTTClient
	if cfg.MQTT.Enabled {
		mqttClient = mqttPkg.NewMQTTClient(&mqttPkg.ClientConfig{
			Broker:   cfg.MQTT.Broker,
			ClientID: cfg.MQTT.ClientID,
			Username: cfg.MQTT.Username,
			Password: cfg.MQTT.Password,
			Topics:   cfg.MQTT.SubscribeTopics,
		})

		mqttCtx, mqttCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer mqttCancel()
		if err := mqttClient.Connect(mqttCtx); err != nil {
			zap.L().Warn("MQTT 连接失败，MQTT 功能不可用", zap.Error(err))
			mqttClient = nil
		} else {
			zap.L().Info("MQTT 客户端已连接", zap.String("broker", cfg.MQTT.Broker))
			defer mqttClient.Disconnect()
		}
	}

	// 初始化路由
	r, hub := router.Setup(cfg, db, rdb, mqttClient, tcpServer)

	// 启动 WebSocket Hub
	go hub.Run()

	// 启动服务
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

	// 优雅关闭
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
