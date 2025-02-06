CREATE TYPE movieStatus AS ENUM('released', 'upcoming');

CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title TEXT,
    synopsis TEXT,
    status movieStatus,
    profile_path TEXT,
    background_path TEXT,
    genre_names TEXT[],
    release_date DATE,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);


CREATE INDEX idx_title_synopsis ON movies USING GIN (to_tsvector('english', title || ' ' || synopsis));
CREATE INDEX idx_release_date ON movies(release_date);
CREATE INDEX idx_genre_names ON movies USING GIN (genre_names);
