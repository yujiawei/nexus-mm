-- Users
CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(26) PRIMARY KEY,
    username    VARCHAR(32) UNIQUE NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    nickname    VARCHAR(64) NOT NULL,
    avatar_url  TEXT DEFAULT '',
    role        VARCHAR(16) NOT NULL DEFAULT 'member',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Teams
CREATE TABLE IF NOT EXISTS teams (
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(32) UNIQUE NOT NULL,
    display_name VARCHAR(128) NOT NULL,
    description  TEXT DEFAULT '',
    creator_id   VARCHAR(26) NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Team Members
CREATE TABLE IF NOT EXISTS team_members (
    team_id    VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id    VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(16) NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, user_id)
);

-- Channels
CREATE TABLE IF NOT EXISTS channels (
    id           VARCHAR(26) PRIMARY KEY,
    team_id      VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name         VARCHAR(32) NOT NULL,
    display_name VARCHAR(128) NOT NULL,
    type         VARCHAR(16) NOT NULL DEFAULT 'open',
    purpose      TEXT DEFAULT '',
    creator_id   VARCHAR(26) NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, name)
);

-- Channel Members
CREATE TABLE IF NOT EXISTS channel_members (
    channel_id VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id    VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(16) NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (channel_id, user_id)
);

-- Messages
CREATE TABLE IF NOT EXISTS messages (
    id         VARCHAR(26) PRIMARY KEY,
    channel_id VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id    VARCHAR(26) NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    type       VARCHAR(16) NOT NULL DEFAULT 'text',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_channel_created ON messages(channel_id, created_at DESC);

-- Incoming Webhooks
CREATE TABLE IF NOT EXISTS incoming_webhooks (
    id           VARCHAR(26) PRIMARY KEY,
    channel_id   VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    team_id      VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    creator_id   VARCHAR(26) NOT NULL REFERENCES users(id),
    display_name VARCHAR(128) NOT NULL,
    description  TEXT DEFAULT '',
    token        VARCHAR(64) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Outgoing Webhooks
CREATE TABLE IF NOT EXISTS outgoing_webhooks (
    id           VARCHAR(26) PRIMARY KEY,
    channel_id   VARCHAR(26) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    team_id      VARCHAR(26) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    creator_id   VARCHAR(26) NOT NULL REFERENCES users(id),
    display_name VARCHAR(128) NOT NULL,
    description  TEXT DEFAULT '',
    token        VARCHAR(64) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outgoing_webhook_triggers (
    webhook_id   VARCHAR(26) NOT NULL REFERENCES outgoing_webhooks(id) ON DELETE CASCADE,
    trigger_word VARCHAR(128) NOT NULL,
    PRIMARY KEY (webhook_id, trigger_word)
);

CREATE TABLE IF NOT EXISTS outgoing_webhook_urls (
    webhook_id   VARCHAR(26) NOT NULL REFERENCES outgoing_webhooks(id) ON DELETE CASCADE,
    callback_url TEXT NOT NULL,
    PRIMARY KEY (webhook_id, callback_url)
);
