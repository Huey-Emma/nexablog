CREATE TABLE IF NOT EXISTS tokens (
  hash BYTEA NOT NULL, 
  user_id INT NOT NULL, 
  scope VARCHAR NOT NULL CHECK(scope = 'authentication'),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL, 
  CONSTRAINT tokens_users_fk FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS tokens_hash_idx ON tokens(hash);
