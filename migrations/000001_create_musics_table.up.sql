CREATE TABLE IF NOT EXISTS musics (
                                      id bigserial PRIMARY KEY,
                                      added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                      title text NOT NULL,
                                      year integer NOT NULL,
                                      duration integer NOT NULL,
                                      genres text[] NOT NULL,
                                      version integer NOT NULL DEFAULT 1
);