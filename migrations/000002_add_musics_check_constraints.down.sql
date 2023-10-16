ALTER TABLE musics ADD CONSTRAINT movies_runtime_check CHECK (duration >= 0);
ALTER TABLE musics ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE musics ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
