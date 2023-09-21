BEGIN;
-- CREATE TYPES
CREATE TYPE enum_two_factor_type AS ENUM ('SMS','EMAIL','TOTP','NONE');
CREATE TYPE enum_send_type AS ENUM ('SMS','EMAIL');
-- CREATE new columns
ALTER TABLE users ADD COLUMN two_factor_type enum_two_factor_type;
-- Update to email 
UPDATE users SET two_factor_type = 'EMAIL' WHERE two_factor_enabled = True;
UPDATE users SET two_factor_type = 'NONE' WHERE two_factor_enabled = False;
-- Add hotop columns
ALTER TABLE users ADD COLUMN totp_secret VARCHAR ;
UPDATE users SET totp_secret = '';
ALTER TABLE users ADD COLUMN totp_url VARCHAR  ;
UPDATE users SET totp_url = '';
ALTER TABLE users ADD COLUMN totp_created TIMESTAMP ;
-- Update previous records send type to email since there was no other option except email
ALTER TABLE two_factor_requests ADD COLUMN send_type enum_send_type ;
UPDATE two_factor_requests SET send_type = 'EMAIL';
ALTER TABLE two_factor_requests ALTER COLUMN send_type SET  NOT NULL;
COMMIT;