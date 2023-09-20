BEGIN;
ALTER TABLE users DROP COLUMN two_factor_type;
ALTER TABLE users DROP COLUMN totp_secret;
ALTER TABLE users DROP COLUMN totp_created;
ALTER TABLE users DROP COLUMN totp_url;
ALTER TABLE two_factor_requests DROP COLUMN send_type;
DROP TYPE two_factor_type;
DROP TYPE send_type;
COMMIT;