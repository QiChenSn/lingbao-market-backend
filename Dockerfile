# 多阶段构建 Dockerfile
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 优化点 1: 替换 Alpine 为国内阿里云源，并安装 git
# 将 sed 和 apk add 合并为一层，减少镜像层数
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git

# 优化点 2: 设置 Go 代理（使用 goproxy.cn 通常比 aliyun 更快且稳定）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# 复制 go mod 文件并下载依赖（利用 Docker 缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（-ldflags="-s -w" 可以大幅缩小二进制体积）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w" -o server ./cmd/server

# --- 最终镜像 ---
FROM alpine:latest

# 优化点 3: 同样替换运行环境的 Alpine 源，加速 ca-certificates 安装
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates && \
    adduser -D -s /bin/sh appuser

WORKDIR /home/appuser

# 从 builder 阶段复制二进制文件
COPY --from=builder /app/server .
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