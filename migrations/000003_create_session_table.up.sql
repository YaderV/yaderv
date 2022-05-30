CREATE TABLE IF NOT EXISTS sessions (
  token text PRIMARY KEY,
  data bytea NOT NULL,
  expiry timestamp NOT NULL
)
