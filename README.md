# Nexus-MM

Enterprise collaboration platform built on [WuKongIM](https://github.com/WuKongIM/WuKongIM), inspired by Mattermost.

## Architecture

```
                        +-----------------+
                        |   Web Frontend  |
                        |  (React + Vite) |
                        +--------+--------+
                                 |
                                 v
+-------------+         +-------+--------+         +-----------+
|  PostgreSQL |<------->|  Nexus-MM API  |<------->| WuKongIM  |
|             |         |   (Go + Gin)   |         | (IM Core) |
+-------------+         +-------+--------+         +-----------+
                                 |
                    +------------+------------+
                    |            |             |
               +----+----+ +----+----+ +------+------+
               |  Redis  | |  Meili  | |  Bot API    |
               |  Cache  | | Search  | | (Telegram-  |
               |         | |         | |  style)     |
               +---------+ +---------+ +-------------+
```

## Features

- Team & channel management
- Real-time messaging via WuKongIM WebSocket
- Thread replies
- Emoji reactions & message pinning
- Full-text search (MeiliSearch)
- Incoming/Outgoing Webhooks
- Slash commands
- Telegram-compatible Bot API with getUpdates polling & webhook delivery
- Agent self-registration system
- Channel categories
- Message retention policies
- Audit logging
- CORS, rate limiting, request logging
- OpenClaw channel plugin

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker & Docker Compose

### 1. Start infrastructure

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, WuKongIM, and MeiliSearch.

### 2. Run database migrations

```bash
psql -h localhost -U nexus -d nexus_mm -f migrations/001_init.sql
```

### 3. Start the API server

```bash
cp configs/nexus.yaml.example configs/nexus.yaml
# Edit configs/nexus.yaml with your settings
go run ./cmd/server/
```

### 4. Start the frontend (dev mode)

```bash
cd web
npm install
npm run dev
```

Open http://localhost:3000 in your browser.

### 5. Build for production

```bash
# Backend
go build -o nexus-mm ./cmd/server/

# Frontend
cd web && npm run build
```

## API Reference

### User API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/users/register` | Register (rate limited) |
| POST | `/api/v1/users/login` | Login, returns JWT |
| GET | `/api/v1/users/me` | Current user info |

### Team API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/teams` | Create a team |
| GET | `/api/v1/teams` | List user's teams |
| GET | `/api/v1/teams/all` | List all teams |
| POST | `/api/v1/teams/:id/join` | Join a team |
| POST | `/api/v1/teams/:id/members` | Add team member |
| GET | `/api/v1/teams/:id/members` | List members |

### Channel API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/teams/:id/channels` | Create channel |
| GET | `/api/v1/teams/:id/channels` | List team channels |
| POST | `/api/v1/channels/:id/join` | Join channel |
| POST | `/api/v1/channels/:id/messages` | Send message |
| GET | `/api/v1/channels/:id/messages` | List messages (cursor pagination) |
| GET | `/api/v1/channels/:id/messages/:msg_id/thread` | Get thread |

### Bot API (Telegram-compatible)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/bot/{token}/getMe` | Get bot info |
| POST | `/bot/{token}/sendMessage` | Send message as bot |
| GET | `/bot/{token}/getUpdates` | Poll for new messages |
| POST | `/bot/{token}/setWebhook` | Set webhook URL |
| POST | `/bot/{token}/sendReaction` | Add emoji reaction |

**sendMessage Body:**
```json
{
  "channel_id": "channel-ulid",
  "content": "Hello from bot!",
  "root_id": "optional-parent-msg-id"
}
```

**getUpdates Response:**
```json
{
  "ok": true,
  "result": [
    {
      "update_id": 1,
      "channel_id": "...",
      "message_id": "...",
      "user_id": "...",
      "content": "Hello",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

## Configuration

Via `configs/nexus.yaml` or environment variables (prefix `NEXUS_`):

| Key | Env Var | Default | Description |
|-----|---------|---------|-------------|
| `server.host` | `NEXUS_SERVER_HOST` | `0.0.0.0` | Listen host |
| `server.port` | `NEXUS_SERVER_PORT` | `8065` | Listen port |
| `database.host` | `NEXUS_DATABASE_HOST` | `localhost` | PostgreSQL host |
| `database.port` | `NEXUS_DATABASE_PORT` | `5432` | PostgreSQL port |
| `database.user` | `NEXUS_DATABASE_USER` | `nexus` | DB user |
| `database.password` | `NEXUS_DATABASE_PASSWORD` | `nexus` | DB password |
| `database.dbname` | `NEXUS_DATABASE_DBNAME` | `nexus_mm` | DB name |
| `jwt.secret` | `NEXUS_JWT_SECRET` | - | JWT signing secret |
| `jwt.expire_hour` | `NEXUS_JWT_EXPIRE_HOUR` | `72` | Token TTL hours |
| `wukong.api_url` | `NEXUS_WUKONG_API_URL` | `http://localhost:5001` | WuKongIM API |
| `wukong.manager_token` | `NEXUS_WUKONG_MANAGER_TOKEN` | - | WuKongIM admin token |
| `wukong.webhook_addr` | `NEXUS_WUKONG_WEBHOOK_ADDR` | `0.0.0.0:6979` | Webhook listen addr |
| `meilisearch.url` | `NEXUS_MEILISEARCH_URL` | `http://localhost:7700` | MeiliSearch URL |
| `redis.addr` | `NEXUS_REDIS_ADDR` | `localhost:6379` | Redis address |

## OpenClaw Plugin

An OpenClaw channel plugin is included at `plugin/openclaw-channel-nexusmm/`. It connects OpenClaw to Nexus-MM's Bot API via polling.

See [plugin/openclaw-channel-nexusmm/README.md](plugin/openclaw-channel-nexusmm/README.md).

## Tech Stack

| Component | Technology |
|-----------|-----------|
| API Server | Go 1.22+ (Gin) |
| IM Engine | WuKongIM |
| Database | PostgreSQL (sqlx) |
| Search | MeiliSearch |
| Cache | Redis (go-redis) |
| Auth | JWT (golang-jwt) |
| Config | Viper |
| Logging | zerolog |
| Frontend | React 18 + TypeScript + Vite |
| State | Zustand |
| Styling | Tailwind CSS |

## License

MIT
