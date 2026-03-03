# Final Delivery Sprint

Read ALL existing source files. Then implement everything below to produce a production-ready product.

## Task 1: OpenClaw Channel Plugin (openclaw-channel-nexusmm)

Create a complete OpenClaw channel plugin at `plugin/openclaw-channel-nexusmm/`.

### Structure:
```
plugin/openclaw-channel-nexusmm/
  package.json
  index.js        - Main plugin entry
  README.md
```

### package.json:
```json
{
  "name": "openclaw-channel-nexusmm",
  "version": "0.1.0",
  "description": "OpenClaw channel plugin for Nexus-MM",
  "main": "index.js",
  "openclaw": {
    "extensions": {
      "nexusmm": {
        "type": "channel"
      }
    }
  },
  "dependencies": {
    "axios": "^1.7.0"
  }
}
```

### index.js - Plugin Logic:
The plugin should:
1. On init: read config (botToken, apiUrl) from openclaw config
2. Start polling getUpdates every 3 seconds
3. When message received from getUpdates → emit to OpenClaw as inbound message
4. When OpenClaw wants to send reply → call POST /bot/{token}/sendMessage
5. Handle reconnection and errors gracefully

Reference the DMWork plugin at ~/.openclaw/extensions/dmwork/ for the plugin API pattern.

Key OpenClaw channel plugin interface:
- `module.exports = function(config) { return { start, stop, send } }`
- `start(emitter)` - Start receiving messages, call emitter.emit('message', {...})
- `stop()` - Cleanup
- `send(message)` - Send outbound message
- Message format: { id, channelId, content, senderId, senderName, timestamp }

### Config (openclaw.json):
```json
{
  "channels": {
    "nexusmm": {
      "botToken": "bf_xxx",
      "apiUrl": "http://35.221.229.58:9876"
    }
  }
}
```

## Task 2: Production Hardening

### 2.1 CORS (internal/server/server.go)
Add CORS middleware to allow frontend on different origins:
```go
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Bot-Token")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})
```

### 2.2 Request Logging
Add request logging middleware with zerolog (method, path, status, latency).

### 2.3 Graceful Error Messages
Review all handlers - ensure consistent error JSON format: `{"error": "message"}` for 4xx, `{"error": "internal error"}` for 5xx (don't leak internals).

### 2.4 Rate Limiting (basic)
Add simple in-memory rate limiter for public endpoints (register, login): max 10 req/min per IP.

### 2.5 Input Validation
- Username: 3-30 chars, alphanumeric + underscore
- Email: valid format
- Password: min 8 chars
- Message content: max 10000 chars
- Channel name: 3-50 chars, lowercase alphanumeric + hyphen

## Task 3: Frontend Polish

### 3.1 Navigation
- Add proper nav bar with links: Chat | Agents | Settings
- Show current user info in sidebar header
- Logout button

### 3.2 Error Handling
- Show toast/alert on API errors
- Loading spinners on async operations
- Empty state messages ("No messages yet", "No channels")

### 3.3 Auto-scroll
- Chat auto-scrolls to bottom on new messages
- Thread panel scrolls independently

## Task 4: Documentation

### 4.1 Update README.md
Complete production README with:
- Project description
- Architecture diagram (text)
- Quick start (docker compose up + go run)
- API reference summary
- Bot API reference
- Frontend screenshots description
- Configuration reference
- License

### 4.2 Deploy Guide (DEPLOY.md)
- Prerequisites
- Docker Compose setup
- Nginx reverse proxy config
- GCP firewall rules
- Systemd service setup
- SSL/TLS (certbot instructions)

## Task 5: Build & Verify

1. `go build -o nexus-mm ./cmd/server/`
2. `cd web && npm run build`
3. `sudo systemctl restart nexus-mm`
4. Run comprehensive E2E test
5. `cd plugin/openclaw-channel-nexusmm && npm install`
6. git add -A && git commit && git push

When finished:
openclaw system event --text "Done: Nexus-MM final delivery - plugin, hardening, docs, polish" --mode now
