-- Team invite links
CREATE TABLE IF NOT EXISTS team_invite_links (
    id          TEXT PRIMARY KEY,
    team_id     TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    code        TEXT NOT NULL UNIQUE,
    creator_id  TEXT NOT NULL REFERENCES users(id),
    expires_at  TIMESTAMPTZ,
    max_uses    INT NOT NULL DEFAULT 0,  -- 0 = unlimited
    use_count   INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_invite_links_code ON team_invite_links(code);
CREATE INDEX IF NOT EXISTS idx_invite_links_team ON team_invite_links(team_id);
