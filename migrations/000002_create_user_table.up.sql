CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  email text UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL
)
