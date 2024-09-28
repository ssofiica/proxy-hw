CREATE TABLE request (
    id SERIAL PRIMARY KEY,
    data JSONB
);

CREATE TABLE response (
    id SERIAL PRIMARY KEY,
    request_id INTEGER REFERENCES request (id) ON DELETE CASCADE,
    data JSONB
);