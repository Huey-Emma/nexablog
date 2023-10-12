CREATE TABLE IF NOT EXISTS users(
  user_id INT generated always as identity,
  username VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  password BYTEA NOT NULL, 
  version INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  PRIMARY KEY(user_id)
);
