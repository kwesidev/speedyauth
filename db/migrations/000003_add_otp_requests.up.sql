BEGIN;
CREATE TABLE otp_requests (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    ip_address VARCHAR(40) ,
    user_agent VARCHAR(200) ,
    request_id VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(6) NOT NULL,
    expiry_time TIMESTAMP NOT NULL,
    send_method enum_send_method NOT NULL,
    created_at TIMESTAMP NOT NULL
);
COMMIT;
