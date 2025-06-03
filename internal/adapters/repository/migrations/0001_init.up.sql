CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    content_type VARCHAR(50) NOT NULL,
    storage_path TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_media_owner ON media(owner_id);

CREATE INDEX IF NOT EXISTS idx_media_created_at ON media(created_at);