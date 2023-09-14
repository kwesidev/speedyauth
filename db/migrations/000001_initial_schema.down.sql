BEGIN;
DROP TABLE IF EXISTS user_refresh_tokens;
DROP TABLE IF EXISTS reset_password_requests;
DROP TABLE IF EXISTS two_factor_requests;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
COMMIT;