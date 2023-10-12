CREATE OR REPLACE function remove_expired_tokens() RETURNS TRIGGER AS $$
BEGIN
  DELETE FROM tokens WHERE expires_at < now() - interval '1 minute';
  return new;
END;
$$
LANGUAGE PLPGSQL;

CREATE OR REPLACE TRIGGER remove_expired_tokens_trigger AFTER INSERT ON tokens EXECUTE PROCEDURE remove_expired_tokens();
