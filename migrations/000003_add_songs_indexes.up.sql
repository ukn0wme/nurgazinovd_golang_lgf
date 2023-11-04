CREATE INDEX IF NOT EXISTS songs_title_idx ON songs USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS songs_genres_idx ON songs USING GIN (genres)