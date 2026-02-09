# 配置管理文档

## 概述

项目已优化配置管理，支持多种配置方式，不再需要硬编码Redis地址和服务端口。

## 配置方式（优先级从高到低）

### 1. 环境变量（最高优先级）
```bash
# Windows PowerShell
$env:SERVER_PORT="9000"
$env:REDIS_ADDR="prod-redis:6379"
./server.exe

# Linux/Unix
export SERVER_PORT=9000
export REDIS_ADDR=prod-redis:6379
./server
```

### 2. YAML配置文件
在项目根目录创建 `config.yaml`：
```yaml
server:
  port: "9000"
  host: "0.0.0.0"
  mode: "release"
  read_timeout: 120

redis:
  addr: "prod-redis:6379"
  password: "your_password"
  db: 1
  pool_size: 20
  min_idle_conns: 10
```

### 3. 默认配置（最低优先级）
如果没有设置环境变量或配置文件，使用以下默认值：
- 服务端口: 8080
- Redis地址: localhost:6379
- Gin模式: debug

## 配置项说明

### 服务器配置
- `SERVER_PORT`: 服务端口，默认 8080
- `SERVER_HOST`: 绑定地址，默认空（绑定所有接口）
- `GIN_MODE`: Gin运行模式，可选 debug/release/test，默认 debug
- `SERVER_READ_TIMEOUT`: 读取超时时间（秒），默认 60

### Redis配置
- `REDIS_ADDR`: Redis服务器地址，默认 localhost:6379
- `REDIS_PASSWORD`: Redis密码，默认空
- `REDIS_DB`: 数据库编号，默认 0
- `REDIS_POOL_SIZE`: 连接池大小，默认 10
- `REDIS_MIN_IDLE_CONNS`: 最小空闲连接数，默认 5

## 部署示例

### 开发环境
```bash
# 使用默认配置
./server

# 或者设置端口
$env:SERVER_PORT="3000"
./server
```

### 生产环境
```bash
# 方式一：环境变量
$env:SERVER_PORT="80"
$env:SERVER_HOST="0.0.0.0"
$env:GIN_MODE="release"
$env:REDIS_ADDR="prod-redis.internal:6379"
$env:REDIS_PASSWORD="prod_password"
$env:REDIS_DB="1"
./server

# 方式二：配置文件
# 复制 config.prod.yaml.example 为 config.yaml 并修改
cp config.prod.yaml.example config.yaml
./server
```

### Docker部署

```dockerfile
# Dockerfile 示例
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY .. .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080

# 使用环境变量
ENV SERVER_PORT=8080
ENV REDIS_ADDR=redis:6379
ENV GIN_MODE=release

CMD ["./server"]
```

```yaml
# docker-compose.yml 示例
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - REDIS_ADDR=redis:6379
      - GIN_MODE=release
    depends_on:
      - redis

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

## 配置文件搜索路径

程序会按以下顺序搜索配置文件：
1. `./config.yaml`
2. `./config/config.yaml`
3. `./configs/config.yaml`

## 注意事项

1. 环境变量会覆盖配置文件中的设置
2. 生产环境建议使用 `GIN_MODE=release`
3. 配置文件支持YAML格式，注意缩进
4. Redis连接池配置可以根据负载调整
5. 不要将包含敏感信息的配置文件提交到版本控制系统
