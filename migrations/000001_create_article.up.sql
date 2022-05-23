CREATE TABLE IF NOT EXISTS articles (
  id serial PRIMARY KEY,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  title text NOT NULL,
  body text NOT NULL,
  categories text[] NOT NULL
);
