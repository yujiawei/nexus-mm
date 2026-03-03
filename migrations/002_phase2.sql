-- Phase 2+3: threads, reactions, pins, slash commands, categories, audit, retention

-- Thread/Reply support
ALTER TABLE messages ADD COLUMN root_id VARCHAR(26) NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN reply_count INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_messages_root_id ON messages(root_id) WHERE root_id != '';

-- Reactions
CREATE TABLE IF NOT EXISTS reactions (
    id         VARCHAR(26) PRIMARY KEY,
    message_id VARCHAR(26) NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id    VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji_name VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (message_id, user_id, emoji_name)
);
CREATE INDEX IF NOT EXISTS idx_reactions_message ON reactions(message_id);

-- Pinned Messages
CREATE TABLE IF NOT EXISTS pinned_messages (
    id         VARCHAR(26) PRIMARY KEY,
    channel_id VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    message_id VARCHAR(26) NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id    VARCHAR(26) NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (channel_id, message_id)
);

-- Slash Commands
CREATE TABLE IF NOT EXISTS slash_commands (
    id         VARCHAR(26) PRIMARY KEY,
    team_id    VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    trigger    VARCHAR(128) NOT NULL,
    url        TEXT NOT NULL,
    method     VARCHAR(8) NOT NULL DEFAULT 'POST',
    creator_id VARCHAR(26) NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, trigger)
);

-- Channel Categories
CREATE TABLE IF NOT EXISTS channel_categories (
    id           VARCHAR(26) PRIMARY KEY,
    user_id      VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id      VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    display_name VARCHAR(128) NOT NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS channel_category_entries (
    category_id VARCHAR(26) NOT NULL REFERENCES channel_categories(id) ON DELETE CASCADE,
    channel_id  VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (category_id, channel_id)
);

-- Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id          VARCHAR(26) PRIMARY KEY,
    user_id     VARCHAR(26) NOT NULL REFERENCES users(id),
    action      VARCHAR(64) NOT NULL,
    entity_type VARCHAR(32) NOT NULL,
    entity_id   VARCHAR(26) NOT NULL,
    ip_addr     VARCHAR(45) DEFAULT '',
    details     JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- Message Retention
ALTER TABLE teams ADD COLUMN retention_days INTEGER NOT NULL DEFAULT 0;
ALTER TABLE channels ADD COLUMN retention_days INTEGER NOT NULL DEFAULT 0;
