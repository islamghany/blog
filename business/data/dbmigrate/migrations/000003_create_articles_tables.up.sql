CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    tags TEXT [] NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS articles_author_id_idx ON articles (author_id);