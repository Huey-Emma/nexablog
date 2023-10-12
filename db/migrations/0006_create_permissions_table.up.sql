CREATE TABLE IF NOT EXISTS permissions (
  permission_id INT generated always as identity,
  code VARCHAR,
  PRIMARY KEY(permission_id)
);

INSERT INTO permissions (code) VALUES ('posts:read'), ('posts:write');

CREATE TABLE IF NOT EXISTS users_permissions (
  permission_id INT,
  user_id INT,
  PRIMARY KEY(permission_id, user_id),
  CONSTRAINT users_permissions_users_fk 
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  CONSTRAINT users_permissions_permissions_fk 
    FOREIGN KEY(permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
