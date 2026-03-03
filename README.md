# Nexus-MM

> Mattermost 企业级功能 × WuKongIM 高性能消息引擎

Nexus-MM 是一个基于 WuKongIM 的企业协作平台，借鉴 Mattermost 的产品设计，提供线程讨论、全文搜索、Webhook 集成、插件系统等企业级功能。

## 架构

```
┌─────────────────────────────────────────────┐
│                 Nexus-MM                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────────┐ │
│  │ Threads  │ │ Search   │ │ Integrations │ │
│  │ 线程讨论  │ │ 全文搜索  │ │ Webhook/Bot  │ │
│  └────┬─────┘ └────┬─────┘ └──────┬───────┘ │
│       │            │              │          │
│  ┌────┴────────────┴──────────────┴───────┐  │
│  │          Nexus API Server (Go)         │  │
│  │     RESTful API + WebSocket Gateway    │  │
│  └────────────────┬───────────────────────┘  │
│                   │                          │
│  ┌────────────────┴───────────────────────┐  │
│  │         WuKongIM (消息引擎)             │  │
│  │   高性能投递 · 百万连接 · 自研协议       │  │
│  └────────────────────────────────────────┘  │
│                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────────┐ │
│  │PostgreSQL│ │  Redis   │ │    MinIO     │ │
│  │ 业务数据  │ │  缓存     │ │  文件存储    │ │
│  └──────────┘ └──────────┘ └──────────────┘ │
└─────────────────────────────────────────────┘
```

## 核心功能（按优先级）

### Phase 1 - 基础框架
- [ ] Go API Server 骨架（Gin）
- [ ] WuKongIM 集成层（gRPC webhook + WS）
- [ ] 用户认证（JWT + OAuth2）
- [ ] 频道/团队 CRUD
- [ ] 消息收发（通过 WuKongIM）

### Phase 2 - Mattermost 核心功能
- [ ] **线程讨论** - 消息回复树
- [ ] **全文搜索** - MeiliSearch 集成
- [ ] **Incoming/Outgoing Webhook**
- [ ] **Slash Commands**
- [ ] **Reactions / Pin**

### Phase 3 - 企业级
- [ ] Plugin System (Go 插件)
- [ ] SSO/LDAP
- [ ] 审计日志
- [ ] 消息保留策略

### Phase 4 - 前端
- [ ] React Web 客户端
- [ ] 移动端适配

## 技术栈

| 组件 | 技术 |
|------|------|
| API Server | Go (Gin) |
| 消息引擎 | WuKongIM |
| 数据库 | PostgreSQL |
| 搜索 | MeiliSearch |
| 缓存 | Redis |
| 文件存储 | MinIO |
| 前端 | React + TypeScript |
| 部署 | Docker Compose |

## License

MIT
