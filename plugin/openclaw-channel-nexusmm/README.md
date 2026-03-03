# openclaw-channel-nexusmm

OpenClaw channel plugin for [Nexus-MM](https://github.com/yujiawei/nexus-mm).

## Setup

1. Install the plugin:
   ```bash
   cd plugin/openclaw-channel-nexusmm
   npm install
   ```

2. Add to your `openclaw.json`:
   ```json
   {
     "channels": {
       "nexusmm": {
         "botToken": "bf_your_bot_token_here",
         "apiUrl": "http://your-nexusmm-server:9876"
       }
     }
   }
   ```

3. The plugin will automatically start polling for messages.

## How It Works

- **Inbound**: Polls `GET /bot/{token}/getUpdates` every 3 seconds for new messages, emits them to OpenClaw.
- **Outbound**: Sends replies via `POST /bot/{token}/sendMessage`.
- **Reconnection**: Retries up to 5 times on error, then backs off for 5 seconds before retrying.

## Configuration

| Key        | Required | Description                        |
|------------|----------|------------------------------------|
| `botToken` | Yes      | Bot token from Nexus-MM bot API    |
| `apiUrl`   | Yes      | Nexus-MM server URL (e.g. `http://host:9876`) |

## Plugin Interface

```js
module.exports = function(config) {
  return {
    start(emitter),  // Start receiving messages
    stop(),          // Clean up
    send(message),   // Send outbound message
  };
};
```
