# 多阶段构建 Dockerfile
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go代理
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
ENV GOSUMDB=off

# 安装git（如果需要）
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 最终镜像
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates

# 创建非root用户
RUN adduser -D -s /bin/sh appuser

WORKDIR /home/appuser

# 从builder阶段复制二进制文件
COPY --from=builder /app/server .

# 复制配置文件示例（可选）
COPY --from=builder /app/config.yaml.example .
COPY --from=builder /app/.env.example .

# 改变文件所有者
RUN chown -R appuser:appuser /home/appuser

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置默认环境变量
ENV SERVER_PORT=8080
ENV GIN_MODE=release

# 启动应用
CMD ["./server"]
