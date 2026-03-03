# Nexus-MM Development Guide

## 项目定位
基于 WuKongIM 的企业协作平台，借鉴 Mattermost 的产品设计。

## Phase 1 任务：搭建完整 Go API Server 骨架

### 目录结构
```
cmd/server/main.go          - 入口
internal/
  config/config.go          - 配置管理 (Viper)
  server/server.go          - HTTP Server 启动
  api/
    router.go               - 路由注册
    middleware/auth.go       - JWT 认证中间件
    v1/
      user.go               - 用户注册/登录/信息
      team.go               - 团队 CRUD
      channel.go            - 频道 CRUD
      message.go            - 消息收发 (通过 WuKongIM)
      webhook.go            - Incoming/Outgoing Webhook
  model/                    - 数据模型
    user.go
    team.go
    channel.go
    message.go
    webhook.go
  store/                    - 数据访问层
    postgres/
      user_store.go
      team_store.go
      channel_store.go
      message_store.go
  wkim/                     - WuKongIM 集成层
    client.go               - WuKongIM HTTP API 客户端
    webhook_handler.go      - gRPC webhook 回调处理
    ws_client.go            - WebSocket 消息监听
  service/                  - 业务逻辑层
    user_service.go
    team_service.go
    channel_service.go
    message_service.go
docker-compose.yaml         - PostgreSQL + Redis + WuKongIM + MeiliSearch
configs/
  nexus.yaml.example        - 配置模板
migrations/
  001_init.sql              - 初始数据库迁移
Makefile                    - 常用命令
```

### WuKongIM 集成关键点
- WuKongIM 通过 gRPC webhook 回调 getSubscribers 获取群成员
- 必须通过 API 添加群成员，直接 INSERT 数据库不会同步到 WuKongIM
- WuKongIM v2 需要 header 带 `token: {managerToken}`
- 消息通过 WuKongIM WebSocket 投递，Nexus-MM 存储业务元数据

### API 设计 (RESTful, Mattermost 风格)
```
POST   /api/v1/users/login
POST   /api/v1/users/register
GET    /api/v1/users/me

POST   /api/v1/teams
GET    /api/v1/teams
GET    /api/v1/teams/:id

POST   /api/v1/teams/:team_id/channels
GET    /api/v1/teams/:team_id/channels
GET    /api/v1/channels/:id

POST   /api/v1/channels/:id/messages
GET    /api/v1/channels/:id/messages

POST   /api/v1/hooks/incoming
POST   /api/v1/hooks/outgoing
```

### 技术要求
- Go 1.22+, Gin framework
- PostgreSQL (sqlx), Redis (go-redis)
- JWT 认证 (golang-jwt)
- 配置: Viper
- 日志: zerolog
- Docker Compose 一键启动
- 所有代码要有合理的错误处理
- 写好 go.mod

### 不要做
- 不要写前端代码
- 不要过度设计，保持简洁
- 不要 mock WuKongIM，写真实的 HTTP 客户端调用
