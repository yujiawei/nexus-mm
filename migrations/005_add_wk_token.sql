-- Add WuKongIM token to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS wk_token TEXT DEFAULT '';
