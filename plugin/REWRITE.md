# Rewrite openclaw-channel-nexusmm as TypeScript Plugin SDK

## Reference
The DMWork plugin at ~/.openclaw/extensions/dmwork/ is the reference implementation.
Read ALL files there to understand the exact plugin pattern.

## Task
Rewrite `plugin/openclaw-channel-nexusmm/` as a proper OpenClaw TypeScript channel plugin.

### Files to create:
```
plugin/openclaw-channel-nexusmm/
  package.json
  index.ts
  tsconfig.json
  openclaw.plugin.json
  src/
    channel.ts        - Main ChannelPlugin implementation
    config-schema.ts  - Zod schema for config
    types.ts          - Type definitions
    polling.ts        - getUpdates polling logic
```

### Key differences from DMWork:
- DMWork uses WuKongIM WebSocket for real-time
- Nexus-MM uses HTTP polling (getUpdates) + HTTP API (sendMessage)
- Much simpler - no WebSocket, no WuKongIM SDK
- Bot API is Telegram-compatible

### Config Schema (channels.nexusmm):
```typescript
{
  botToken: string;    // bf_xxx bot token
  apiUrl: string;      // http://host:9876
  pollIntervalMs?: number;  // default 3000
}
```

### Channel Plugin Implementation:
1. `connect()`: Start polling getUpdates every pollIntervalMs
2. `disconnect()`: Stop polling
3. `send()`: POST /bot/{token}/sendMessage
4. Map Nexus-MM messages to OpenClaw InboundMessage format
5. Map OpenClaw OutboundMessage to Nexus-MM sendMessage format

### InboundMessage mapping:
```typescript
{
  id: update.message.id,
  channelId: `nexusmm:${update.message.channel_id}`,
  content: update.message.content,
  senderId: update.message.user_id,
  senderName: update.message.user_id, // or username if available
  timestamp: new Date(update.message.created_at).getTime(),
}
```

### package.json key fields:
- "type": "module"
- "main": "index.ts"
- peerDependencies: { "openclaw": ">=2026.2.0" }
- devDependencies: { "openclaw": "2026.2.12", "typescript": "^5.4.0" }
- dependencies: { "axios": "^1.7.0", "zod": "^4.0.0" }
- openclaw.extensions: ["./index.ts"]
- openclaw.channel: { id: "nexusmm", label: "Nexus-MM", ... }

### openclaw.plugin.json:
```json
{
  "id": "nexusmm",
  "channels": ["nexusmm"],
  "configSchema": { "type": "object", "additionalProperties": false, "properties": {} }
}
```

## After rewriting:
1. npm install in plugin dir (to get types)
2. npx tsc --noEmit to type-check (may have errors, that's ok for now - OpenClaw loads .ts directly)
3. npm publish --access public (bump to 0.5.0)
4. Test with Docker: docker run openclaw with the plugin
5. git add -A && git commit && git push

## Docker E2E Test:
After publishing, run this Docker test:
```bash
docker rm -f nexus-agent-alpha 2>/dev/null
docker run -d --name nexus-agent-alpha \
  ghcr.io/openclaw/openclaw:latest \
  sh -c '
    openclaw config set channels.nexusmm.botToken "bf_0a047a1bb3fbf41e51d245d94b078a9e98e94b839a3fc341"
    openclaw config set channels.nexusmm.apiUrl "http://35.221.229.58:9876"
    openclaw config set gateway.mode local
    openclaw plugins install openclaw-channel-nexusmm
    openclaw plugins enable nexusmm
    exec openclaw gateway run
  '
```

When done: openclaw system event --text "Done: Nexus-MM plugin rewritten as TypeScript SDK, published, Docker E2E tested" --mode now
