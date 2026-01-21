-- Videos: core content entity
CREATE TABLE IF NOT EXISTS videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    language TEXT DEFAULT 'oji',
    duration_seconds INTEGER DEFAULT 0,
    file_path TEXT NOT NULL,
    thumbnail_path TEXT DEFAULT '',
    file_size_bytes INTEGER DEFAULT 0,
    mime_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    rejection_reason TEXT DEFAULT '',
    device_id TEXT,
    created_at DATETIME NOT NULL,
    published_at DATETIME,
    moderated_at DATETIME,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

-- Devices: trust tracking for kiosks
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    api_key_hash TEXT NOT NULL,
    trust_level TEXT NOT NULL DEFAULT 'trusted',
    last_seen_at DATETIME,
    created_at DATETIME NOT NULL
);

-- Tags: categorization
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS video_tags (
    video_id TEXT NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (video_id, tag_id),
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_created ON videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_device ON videos(device_id);
