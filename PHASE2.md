# Phase 2 + 3 Implementation Plan

Read all existing Go code first to understand the patterns, then implement these features.

## Phase 2: Mattermost Core Features

### 2.1 Thread/Reply System
- Add `root_id` and `reply_count` fields to messages table (migration 002)
- Messages with root_id are replies; root messages track reply_count
- API: GET /api/v1/channels/:id/messages/:msg_id/thread - get all replies
- API: POST /api/v1/channels/:id/messages with `root_id` field for replies
- Update message_store and message_service

### 2.2 Full-Text Search (MeiliSearch)
- Add MeiliSearch client in internal/search/meilisearch.go
- Index messages on creation (async via goroutine)
- API: GET /api/v1/search?q=keyword&channel_id=xxx&team_id=xxx
- Search results include message content, channel name, sender info
- Add search handler in api/v1/search.go

### 2.3 Reactions
- New table: reactions (message_id, user_id, emoji_name, created_at)
- Migration 002
- API: POST /api/v1/channels/:id/messages/:msg_id/reactions {emoji_name}
- API: DELETE /api/v1/channels/:id/messages/:msg_id/reactions/:emoji
- API: GET included in message response
- New model, store, handler

### 2.4 Pin Messages
- New table: pinned_messages (channel_id, message_id, user_id, created_at)
- API: POST /api/v1/channels/:id/messages/:msg_id/pin
- API: DELETE /api/v1/channels/:id/messages/:msg_id/pin
- API: GET /api/v1/channels/:id/pinned

### 2.5 Webhook Enhancement
- Incoming webhook: POST /api/v1/hooks/incoming/:token - accept payload, post as bot
- Outgoing webhook: trigger on message patterns, POST to callback URL
- Integrate with message_service to trigger outgoing webhooks on new messages

### 2.6 Slash Commands
- New table: slash_commands (id, team_id, trigger, url, method, creator_id, ...)
- API: POST /api/v1/teams/:id/commands - register command
- When message starts with /trigger, POST to command URL with context
- Response posted back to channel

## Phase 3: Enterprise Features

### 3.1 Channel Categories / Sidebar
- New table: channel_categories (id, user_id, team_id, display_name, sort_order)
- New table: channel_category_entries (category_id, channel_id, sort_order)
- API: CRUD for categories per user per team

### 3.2 Audit Log
- New table: audit_logs (id, user_id, action, entity_type, entity_id, ip_addr, details JSONB, created_at)
- Middleware to log key actions (login, channel create, message delete, etc.)
- API: GET /api/v1/admin/audit?action=xxx&user_id=xxx (admin only)

### 3.3 Message Retention Policy
- Add retention_days field to teams and channels tables
- Background goroutine (ticker) to delete expired messages
- API: PUT /api/v1/teams/:id/retention, PUT /api/v1/channels/:id/retention

## Technical Notes
- Follow existing code patterns (Gin handlers, sqlx stores, service layer)
- Use existing JWT middleware
- Put migration in migrations/002_phase2.sql (single file for all new tables/columns)
- Register all new routes in router.go
- Keep error handling consistent with existing code
- Run `go build ./...` to verify compilation
- Commit in logical groups

When completely finished:
1. Run `go build ./...` to verify
2. Push to origin master  
3. Run: openclaw system event --text "Done: Nexus-MM Phase 2+3 complete - threads, search, reactions, pins, slash commands, audit, retention" --mode now
