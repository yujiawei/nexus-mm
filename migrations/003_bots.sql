-- Bot accounts
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_bot BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_token VARCHAR(64) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_owner_id VARCHAR(26) REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_description TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS bot_webhook_url TEXT DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_users_bot_token ON users(bot_token) WHERE bot_token IS NOT NULL;

-- Bot updates for getUpdates polling
CREATE TABLE IF NOT EXISTS bot_updates (
    id BIGSERIAL PRIMARY KEY,
    bot_user_id VARCHAR(26) NOT NULL REFERENCES users(id),
    message_id VARCHAR(26) NOT NULL REFERENCES messages(id),
    channel_id VARCHAR(26) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_bot_updates_bot ON bot_updates(bot_user_id, id);
