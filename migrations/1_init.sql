-- +goose Up

CREATE TABLE sessions(
  user_id BIGINT,
    user_name TEXT,
  refresh_token TEXT,
  expires_at timestamp,
  PRIMARY KEY(user_id, refresh_token)
);

--+goose Down
DROP TABLE sessions;
