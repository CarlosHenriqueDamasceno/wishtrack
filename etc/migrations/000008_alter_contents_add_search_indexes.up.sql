CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_contents_name ON contents USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_contents_category ON contents USING gin (category gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_contents_summary ON contents USING gin (summary gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_genres ON contents USING gin (genres);