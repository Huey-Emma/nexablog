CREATE TABLE IF NOT EXISTS posts (
  post_id INT generated always as identity, 
  title VARCHAR NOT NULL,
  body TEXT NOT NULL,
  version INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  PRIMARY KEY(post_id)
);
