CREATE TABLE search_caches (
    id bigserial PRIMARY KEY,
    query VARCHAR(255) NOT NULL,
    embedding vector(1536)
);