-- +goose Up
CREATE TABLE users(
    id uuid PRIMARY KEY,
    name TEXT,
    email TEXT,
    password TEXT,
    created_at timestamp
);

CREATE TABLE sessions(
  user_id TEXT PRIMARY KEY,
  refresh_token TEXT,
  expires_at timestamp
);

--+goose Down
DROP TABLE sessions;
DROP TABLE users;