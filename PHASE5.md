# Phase 5: Bot/Agent System

Read ALL existing Go source files first. Then implement everything below.

## Overview
Add Bot account system with Telegram-compatible Bot API, self-registration, email-based user binding, 
WebSocket real-time messaging, and an Agent integration page in the frontend.

## 5A: Bot Account Model

### Database Changes (migrations/003_bots.sql)
```sql
-- Bot accounts
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_bot BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_token VARCHAR(64) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_owner_id VARCHAR(26) REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_description TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_webhook_url TEXT DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_users_bot_token ON users(bot_token) WHERE bot_token IS NOT NULL;
```

### Model Changes (internal/model/user.go)
Add fields: IsBot, BotToken, BotOwnerID, BotDescription, BotWebhookURL

## 5B: Bot API (Telegram-compatible style)

### New file: internal/api/v1/bot_api.go

**Public endpoints (token in URL, no JWT):**

`POST /bot/{token}/sendMessage`
```json
{"channel_id": "xxx", "content": "Hello!", "root_id": "optional"}
```

`GET /bot/{token}/getUpdates?offset=0&limit=100`
Returns unread messages for this bot since offset. Store last_update_id per bot.
```json
{"ok": true, "result": [{"update_id": 1, "message": {"id": "xxx", "channel_id": "xxx", "user_id": "xxx", "content": "Hi bot", "created_at": "..."}}]}
```

`POST /bot/{token}/setWebhook`
```json
{"url": "https://example.com/webhook"}
```
When set, POST new messages to this URL instead of requiring getUpdates.

`GET /bot/{token}/getMe`
Returns bot info.

`POST /bot/{token}/sendReaction`
```json
{"channel_id": "xxx", "message_id": "xxx", "emoji_name": "thumbsup"}
```

### Bot message delivery
- When a message is sent in a channel where bot is a member:
  - If bot has webhook_url set → POST to webhook
  - Otherwise → store in bot_updates table for getUpdates polling

### New table: bot_updates
```sql
CREATE TABLE IF NOT EXISTS bot_updates (
    id BIGSERIAL PRIMARY KEY,
    bot_user_id VARCHAR(26) NOT NULL REFERENCES users(id),
    message_id VARCHAR(26) NOT NULL REFERENCES messages(id),
    channel_id VARCHAR(26) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_bot_updates_bot ON bot_updates(bot_user_id, id);
```

## 5C: Self-Registration API

### New file: internal/api/v1/agent.go

`POST /api/v1/agents/register`
```json
{"name": "MyAgent", "email": "owner@example.com"}
```
Response:
```json
{
  "bot_user_id": "xxx",
  "bot_token": "bf_xxx...",
  "owner_user_id": "xxx",
  "message": "Bot created and bound to owner. Use bot_token for Bot API."
}
```

Flow:
1. Find or create user by email (password auto-generated, can reset later)
2. Create bot user (is_bot=true, username=name, bot_owner_id=owner)
3. Generate bot_token (bf_ + 48 random hex chars)
4. Auto-add bot to owner's teams
5. Return credentials

`POST /api/v1/agents/bind`
Bind existing bot to owner by email.
```json
{"email": "owner@example.com"}
```
Header: `X-Bot-Token: bf_xxx`

## 5D: Skill.md Endpoint

`GET /skill.md` — Serve a dynamic skill.md that agents can read to self-configure.

Content should include:
- API base URL
- Registration instructions
- Bot API reference
- OpenClaw plugin install instructions
- Quick start prompt

## 5E: Frontend - Agent Integration Page

### New route: /agents
### New components:
- `web/src/components/Agent/AgentPage.tsx` - Main agent management page
- Shows: 
  1. "Create Bot" form (name + description)
  2. List of user's bots with tokens
  3. "Connect Existing Agent" section with copyable prompt
  4. Bot API reference section

### Agent page UI:
```
┌─────────────────────────────────────────────┐
│  🤖 My Agents                               │
│                                              │
│  ┌─ Create New Bot ──────────────────────┐   │
│  │ Name: [________]                      │   │
│  │ Description: [________________]       │   │
│  │ [Create Bot]                          │   │
│  └───────────────────────────────────────┘   │
│                                              │
│  ┌─ My Bots ─────────────────────────────┐   │
│  │ 🤖 MyAgent                            │   │
│  │ Token: bf_abc... [Copy] [Regenerate]  │   │
│  │ Webhook: https://... [Edit]           │   │
│  │ Status: Online                        │   │
│  └───────────────────────────────────────┘   │
│                                              │
│  ┌─ Quick Connect ───────────────────────┐   │
│  │ Copy this prompt to your AI agent:    │   │
│  │ ┌──────────────────────────────────┐  │   │
│  │ │ "Read https://xxx/skill.md       │  │   │
│  │ │  to install Nexus-MM channel...  │  │   │
│  │ │  My email is xxx@xxx.com"        │  │   │
│  │ └──────────────────────────────────┘  │   │
│  │ [Copy Prompt]                         │   │
│  └───────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## 5F: OpenClaw Channel Plugin Skeleton

Create `internal/skill/skill.go` that serves the skill.md content at GET /skill.md
The skill.md should be a Go template that injects the current server URL.

## Implementation Notes
- Follow existing code patterns exactly
- Bot token format: "bf_" + 48 random hex chars  
- Bot API routes go OUTSIDE the JWT auth middleware (token-in-URL auth)
- Register route is public (no auth)
- Add bot_api routes in router.go
- Agent page needs new React route in App.tsx
- Make sure `go build ./...` passes
- Make sure `npm run build` passes

## When Done
1. go build -o nexus-mm ./cmd/server/
2. cd web && npm run build
3. sudo systemctl restart nexus-mm
4. Test: register bot via API, send message via bot API, verify getUpdates
5. git add -A && git commit -m "feat: bot/agent system with self-registration and Telegram-compatible API" && git push origin master
6. openclaw system event --text "Done: Nexus-MM Phase 5 - Bot/Agent system with self-registration, Bot API, agent page" --mode now
