BEGIN;
ALTER TABLE users DROP COLUMN two_factor_method;
ALTER TABLE users DROP COLUMN totp_secret;
ALTER TABLE users DROP COLUMN totp_created;
ALTER TABLE users DROP COLUMN totp_url;
ALTER TABLE two_factor_requests DROP COLUMN send_method;
DROP TYPE two_factor_method;
DROP TYPE send_method;
COMMIT;