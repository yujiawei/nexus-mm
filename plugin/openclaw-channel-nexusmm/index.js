/**
 * openclaw-channel-nexusmm
 *
 * OpenClaw channel plugin for Nexus-MM.
 * Polls the Nexus-MM Bot API (getUpdates) for inbound messages
 * and sends replies via POST /bot/{token}/sendMessage.
 */

const axios = require('axios');

const POLL_INTERVAL_MS = 3000;
const MAX_RETRIES = 5;
const RETRY_BACKOFF_MS = 5000;

module.exports = function nexusmm(config) {
  const botToken = config.botToken;
  const apiUrl = (config.apiUrl || '').replace(/\/$/, '');

  if (!botToken || !apiUrl) {
    throw new Error('nexusmm plugin requires botToken and apiUrl in config');
  }

  let emitter = null;
  let polling = false;
  let pollTimer = null;
  let offset = 0;
  let retryCount = 0;

  const client = axios.create({
    baseURL: `${apiUrl}/bot/${botToken}`,
    timeout: 30000,
    headers: { 'Content-Type': 'application/json' },
  });

  async function pollUpdates() {
    if (!polling) return;

    try {
      const res = await client.get('/getUpdates', {
        params: { offset, limit: 100 },
      });

      const data = res.data;
      if (data.ok && Array.isArray(data.result)) {
        for (const update of data.result) {
          // Advance offset past this update.
          if (update.update_id >= offset) {
            offset = update.update_id + 1;
          }

          // Emit inbound message to OpenClaw.
          if (emitter) {
            emitter.emit('message', {
              id: update.message_id || String(update.update_id),
              channelId: update.channel_id || '',
              content: update.content || '',
              senderId: update.user_id || '',
              senderName: update.user_id || '',
              timestamp: update.created_at
                ? new Date(update.created_at).getTime()
                : Date.now(),
            });
          }
        }
      }

      // Reset retry count on success.
      retryCount = 0;
    } catch (err) {
      retryCount++;
      const msg = err.response
        ? `HTTP ${err.response.status}: ${JSON.stringify(err.response.data)}`
        : err.message;
      console.error(`[nexusmm] poll error (attempt ${retryCount}): ${msg}`);

      if (retryCount >= MAX_RETRIES) {
        console.error('[nexusmm] max retries reached, backing off...');
        retryCount = 0;
        await new Promise((r) => setTimeout(r, RETRY_BACKOFF_MS));
      }
    }

    // Schedule next poll.
    if (polling) {
      pollTimer = setTimeout(pollUpdates, POLL_INTERVAL_MS);
    }
  }

  return {
    /**
     * Start receiving messages from Nexus-MM.
     * @param {EventEmitter} em - OpenClaw emitter for inbound messages.
     */
    start(em) {
      emitter = em;
      polling = true;
      offset = 0;
      retryCount = 0;
      console.log(`[nexusmm] starting poll loop -> ${apiUrl}/bot/***`);
      pollUpdates();
    },

    /**
     * Stop the polling loop and clean up.
     */
    stop() {
      polling = false;
      if (pollTimer) {
        clearTimeout(pollTimer);
        pollTimer = null;
      }
      emitter = null;
      console.log('[nexusmm] stopped');
    },

    /**
     * Send an outbound message via the Nexus-MM Bot API.
     * @param {Object} message - { channelId, content }
     */
    async send(message) {
      if (!message || !message.content) return;

      const channelId = message.channelId;
      if (!channelId) {
        console.error('[nexusmm] send: missing channelId');
        return;
      }

      try {
        const res = await client.post('/sendMessage', {
          channel_id: channelId,
          content: message.content,
        });

        if (res.data && !res.data.ok) {
          console.error(`[nexusmm] sendMessage failed: ${res.data.error}`);
        }
      } catch (err) {
        const msg = err.response
          ? `HTTP ${err.response.status}`
          : err.message;
        console.error(`[nexusmm] sendMessage error: ${msg}`);
      }
    },
  };
};
