BEGIN;
-- CREATE methodS
CREATE TYPE enum_two_factor_method AS ENUM ('SMS','EMAIL','TOTP','NONE');
CREATE TYPE enum_send_method AS ENUM ('SMS','EMAIL');
-- CREATE new columns
ALTER TABLE users ADD COLUMN two_factor_method enum_two_factor_method;
-- Update to email 
UPDATE users SET two_factor_method = 'EMAIL' WHERE two_factor_enabled = True;
UPDATE users SET two_factor_method = 'NONE' WHERE two_factor_enabled = False;
-- Add hotop columns
ALTER TABLE users ADD COLUMN totp_secret VARCHAR ;
UPDATE users SET totp_secret = '';
ALTER TABLE users ADD COLUMN totp_url VARCHAR  ;
UPDATE users SET totp_url = '';
ALTER TABLE users ADD COLUMN totp_created TIMESTAMP ;
-- Update previous records send method to email since there was no other option except email
ALTER TABLE two_factor_requests ADD COLUMN send_method enum_send_method ;
UPDATE two_factor_requests SET send_method = 'EMAIL';
ALTER TABLE two_factor_requests ALTER COLUMN send_method SET  NOT NULL;
COMMIT;