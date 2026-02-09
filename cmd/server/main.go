package main

import (
	"context"
	"log"

	"lingbao-market-backend/internal/api"
	"lingbao-market-backend/internal/config"
	"lingbao-market-backend/internal/repository"
	"lingbao-market-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	// 依赖注入
	repo := repository.NewShareCodeRepo(rdb)
	svc := service.NewShareCodeService(repo)
	handler := api.NewHandler(svc)

	// 启动 Gin
	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	// 路由定义
	r.POST("/sharecode", handler.CreateShareCode)
	r.GET("/sharecode", handler.ListShareCodes)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("服务启动在 %s (模式: %s)", serverAddr, cfg.Server.Mode)
	err := r.Run(serverAddr)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
		return
	}
}
