# lingbao-market-backend

## 项目简介

基于Gin开发的灵宝市集分享码共享平台后端。

## 快速开始

### 1. 编译运行

```bash
# 编译
go build ./cmd/server

# 运行（使用默认配置）
./server

# 服务将启动在 http://localhost:8080
```

### 2. 配置管理

项目支持多种配置方式，优先级从高到低：

1. **环境变量**（最高优先级）
2. **YAML配置文件**
3. **默认配置**（最低优先级）

#### 环境变量配置
```bash
# Windows PowerShell
$env:SERVER_PORT="9000"
$env:REDIS_ADDR="prod-redis:6379"
./server

# Linux/Unix
export SERVER_PORT=9000
export REDIS_ADDR=prod-redis:6379
./server
```

#### 配置文件
创建 `config.yaml` 文件：
```yaml
server:
  port: "9000"
  host: "0.0.0.0"
  mode: "release"

redis:
  addr: "prod-redis:6379"
  password: "your_password"
  db: 1
  pool_size: 20
```

详细配置说明请参考 [CONFIG.md](docs/CONFIG.md)

### 生产环境
1. 设置环境变量或创建配置文件
2. 确保Redis服务可用
3. 设置 `GIN_MODE=release`
4. 配置反向代理（如Nginx）

## 技术栈

- **Web框架**: Gin
- **数据库**: Redis
- **语言**: Go 1.25
- **配置**: YAML + 环境变量

## 项目结构

```
├── cmd/server/          # 程序入口
├── internal/
│   ├── api/            # API处理器
│   ├── config/         # 配置管理
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   └── service/        # 业务逻辑层
├── pkg/response/       # 通用响应
├── config.yaml.example # 配置文件示例
├── .env.example        # 环境变量示例
└── CONFIG.md          # 详细配置文档
```
