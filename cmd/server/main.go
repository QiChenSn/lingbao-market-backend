package main

import (
	"context"
	"log"

	"lingbao-market-backend/internal/api"
	"lingbao-market-backend/internal/repository"
	"lingbao-market-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	// 依赖注入
	repo := repository.NewShareCodeRepo(rdb)
	svc := service.NewShareCodeService(repo)
	handler := api.NewHandler(svc)

	// 启动 Gin
	r := gin.Default()

	// 路由定义
	r.POST("/sharecode", handler.CreateShareCode)
	r.GET("/sharecode", handler.ListShareCodes)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
		return
	}
	log.Println("服务启动在 :8080")
}
